package repository

import (
	"forum/pkg/models"
	"github.com/jackc/pgx"
	"log"
)

const (
	SelectUsersByForumDesc      = `SELECT U.nickname, U.fullname, U.about, U.email FROM dbproject."Users_by_Forum" users INNER JOIN dbproject."User" U on U.nickname = users."user" AND users.forum = $1 ORDER BY users."user" DESC LIMIT $2`
	SelectUsersByForum          = `SELECT U.nickname, U.fullname, U.about, U.email FROM dbproject."Users_by_Forum" users INNER JOIN dbproject."User" U on U.nickname = users."user" AND users.forum = $1 ORDER BY users."user" LIMIT $2`
	SelectUsersByForumSinceDesc = `SELECT U.nickname, U.fullname, U.about, U.email FROM dbproject."Users_by_Forum" users INNER JOIN dbproject."User" U on U.nickname = users."user" AND users.forum = $1 AND U.nickname < $2 ORDER BY users."user" DESC LIMIT $3`
	SelectUsersByForumSince     = `SELECT U.nickname, U.fullname, U.about, U.email FROM dbproject."Users_by_Forum" users INNER JOIN dbproject."User" U on U.nickname = users."user" AND users.forum = $1 AND U.nickname > $2 ORDER BY users."user" LIMIT $3`
	InsertForum                 = `INSERT INTO dbproject."Forum" (title, "user", slug, posts, threads) VALUES ($1, (SELECT nickname FROM dbproject."User" WHERE nickname = $2),$3,0,0) RETURNING "user"`
	SelectForum                 = `SELECT f.slug, f.title, f."user", f.posts, f.threads from dbproject."Forum" f WHERE slug = $1`
)

type ForumRepositoryInterface interface {
	CreateForum(forum models.Forum) (models.Forum, error)
	GetForumInfo(slug string) (models.Forum, bool)
	FindUsers(slug string, params models.ParamsForSearch) ([]models.User, bool)
}

type ForumRepository struct {
	DB *pgx.ConnPool
}

func (r ForumRepository) FindUsers(slug string, params models.ParamsForSearch) ([]models.User, bool) {
	var rows *pgx.Rows
	var err error

	if params.Since == "" {
		if params.Desc {
			rows, err = r.DB.Query("SelectUsersByForumDesc", slug, params.Limit)
		} else {
			rows, err = r.DB.Query("SelectUsersByForum", slug, params.Limit)
		}
	} else {
		if params.Desc {
			rows, err = r.DB.Query("SelectUsersByForumSinceDesc", slug, params.Since, params.Limit)
		} else {
			rows, err = r.DB.Query("SelectUsersByForumSince", slug, params.Since, params.Limit)
		}
	}

	var users []models.User
	if err != nil {
		log.Println(err)
		rows.Close()
		return users, false
	}

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
		if err != nil {
			log.Println(err)
			rows.Close()
			return []models.User{}, false
		}
		users = append(users, user)
	}
	rows.Close()
	return users, true
}

func (r ForumRepository) CreateForum(forum models.Forum) (models.Forum, error) {
	err := r.DB.QueryRow("InsertForum", forum.Title, forum.User, forum.Slug).Scan(&forum.User)
	return forum, err
}

func (r ForumRepository) GetForumInfo(slug string) (models.Forum, bool) {
	var forum models.Forum = models.Forum{Slug: slug}
	err := r.DB.QueryRow("SelectForum", forum.Slug).Scan(&forum.Slug, &forum.Title, &forum.User, &forum.Posts, &forum.Threads)

	if err != nil {
		log.Println(err)
		return models.Forum{}, false
	}

	return forum, true
}
