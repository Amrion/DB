package repository

import (
	"forum/internal/utils/utils"
	"forum/pkg/models"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx"
	"log"
	"strconv"
)

const (
	SelectThreadIdBySlug   = `SELECT id FROM dbproject."Thread" WHERE slug = $1`
	InsertVote             = `INSERT INTO dbproject."Vote" (threadid, "user", value) VALUES ($1, $2, $3)`
	UpdateVote             = `UPDATE dbproject."Vote" SET value = $1 WHERE threadid = $2 AND "user" = $3`
	UpdateThreadId         = `UPDATE dbproject."Thread" SET title = COALESCE(NULLIF($1, ''), title), message = COALESCE(NULLIF($2, ''), message) WHERE id = $3 RETURNING *`
	UpdateThreadSlug       = `UPDATE dbproject."Thread" SET title = COALESCE(NULLIF($1, ''), title), message = COALESCE(NULLIF($2, ''), message) WHERE slug = $3 RETURNING *`
	SelectThreadInfoBySlug = `SELECT id, title, author, forum, message, votes, slug, created FROM dbproject."Thread" WHERE slug = $1`
	SelectThreadInfoById   = `SELECT id, title, author, forum, message, votes, slug, created FROM dbproject."Thread" WHERE id = $1`
	SelectThreadDesc       = `SELECT t.id, t.title, t.author, t.forum, t.message, t.votes, t.slug, t.created FROM dbproject."Thread" t WHERE t.forum = $1 ORDER BY t.created DESC LIMIT $2`
	SelectThread           = `SELECT t.id, t.title, t.author, t.forum, t.message, t.votes, t.slug, t.created FROM dbproject."Thread" t WHERE t.forum = $1 ORDER BY t.created LIMIT $2`
	SelectThreadSinceDesc  = `SELECT t.id, t.title, t.author, t.forum, t.message, t.votes, t.slug, t.created FROM dbproject."Thread" t WHERE t.forum = $1 AND t.created <= $2 ORDER BY t.created DESC LIMIT $3`
	SelectThreadSince      = `SELECT t.id, t.title, t.author, t.forum, t.message, t.votes, t.slug, t.created FROM dbproject."Thread" t WHERE t.forum = $1 AND t.created >= $2 ORDER BY t.created  LIMIT $3`
	InsertThread           = `INSERT INTO dbproject."Thread" (title, author, forum, message, votes, slug, created) VALUES ($1,(SELECT nickname from dbproject."User" where nickname = $2),(SELECT slug from dbproject."Forum"  where slug = $3),$4,0,$5,$6) RETURNING id, forum, author, slug`
)

type ThreadRepositoryInterface interface {
	CreateThread(thread models.Thread) (models.Thread, error)
	FindThreads(slug string, params models.ParamsForSearch) ([]models.Thread, bool)
	GetThreadInfoBySlug(slug string) (models.Thread, bool)
	GetThreadInfoById(id int) (models.Thread, bool)
	UpdateThread(update models.ThreadUpdate, slugOrId string) (models.Thread, bool)
	SetVote(vote models.Vote, id int) bool
	GetThreadIdBySlug(slug string) (int, bool)
}

type ThreadRepository struct {
	DB *pgx.ConnPool
}

func (r ThreadRepository) GetThreadIdBySlug(slug string) (int, bool) {
	id := -1
	err := r.DB.QueryRow("SelectThreadIdBySlug", slug).Scan(&id)
	if err != nil {
		log.Println(err)
		return -1, false
	}

	log.Println(id)

	return id, true
}

func (r ThreadRepository) SetVote(vote models.Vote, id int) bool {
	_, err := r.DB.Exec("InsertVote", id, vote.Nickname, int32(vote.Voice))
	if err == nil {
		log.Println("Add vote to thread")
		return true
	}

	log.Println(err)
	code := utils.PgxErrorCode(err)

	// duplicate key value violates unique constraint "onlyonevote" (SQLSTATE 23505)
	if code == "23505" {
		_, err = r.DB.Exec("UpdateVote", int32(vote.Voice), id, vote.Nickname)
		if err != nil {
			return false
		}
	}

	if code == "23503" {
		return false
	}

	return true
}

func (r ThreadRepository) UpdateThread(update models.ThreadUpdate, slugOrId string) (models.Thread, bool) {
	var thread models.Thread
	id, err := strconv.Atoi(slugOrId)
	if err != nil {
		err = r.DB.QueryRow("UpdateThreadSlug", update.Title, update.Message, slugOrId).
			Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
	} else {
		err = r.DB.QueryRow("UpdateThreadId", update.Title, update.Message, id).
			Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
	}

	if err != nil {
		return models.Thread{}, false
	}
	return thread, true
}

func (r ThreadRepository) GetThreadInfoBySlug(slug string) (models.Thread, bool) {
	var thread models.Thread
	err := r.DB.QueryRow("SelectThreadInfoBySlug", slug).
		Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
	if err != nil {
		return models.Thread{}, false
	}

	if utils.IsValidUUID(thread.Slug) {
		thread.Slug = ""
	}

	return thread, true
}

func (r ThreadRepository) GetThreadInfoById(id int) (models.Thread, bool) {
	var thread models.Thread
	err := r.DB.QueryRow("SelectThreadInfoById", id).
		Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
	if err != nil {
		return models.Thread{}, false
	}

	if utils.IsValidUUID(thread.Slug) {
		thread.Slug = ""
	}

	return thread, true
}

func (r ThreadRepository) FindThreads(slug string, params models.ParamsForSearch) ([]models.Thread, bool) {
	var rows *pgx.Rows
	var err error

	if params.Since == "" {
		if params.Desc {
			rows, err = r.DB.Query("SelectThreadDesc", slug, params.Limit)
		} else {
			rows, err = r.DB.Query("SelectThread", slug, params.Limit)
		}
	} else {
		if params.Desc {
			rows, err = r.DB.Query("SelectThreadSinceDesc", slug, params.Since, params.Limit)
		} else {
			rows, err = r.DB.Query("SelectThreadSince", slug, params.Since, params.Limit)
		}
	}

	if err != nil {
		log.Println(err)
		return nil, false
	}

	var threads []models.Thread

	for rows.Next() {
		var thread models.Thread
		err := rows.Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
		if err != nil {
			log.Println(err)
			rows.Close()
			return nil, false
		}
		if utils.IsValidUUID(thread.Slug) {
			thread.Slug = ""
		}
		threads = append(threads, thread)
	}

	rows.Close()
	return threads, true
}

func (r *ThreadRepository) CreateThread(thread models.Thread) (models.Thread, error) {
	var err error

	if thread.Slug == "" {
		gen, _ := uuid.NewV4()
		thread.Slug = gen.String()
	}

	err = r.DB.QueryRow("InsertThread", thread.Title, thread.Author, thread.Forum, thread.Message, thread.Slug, thread.Created).
		Scan(&thread.Id, &thread.Forum, &thread.Author, &thread.Slug)

	if utils.IsValidUUID(thread.Slug) {
		thread.Slug = ""
	}
	return thread, err
}
