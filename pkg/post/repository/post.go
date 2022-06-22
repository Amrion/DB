package repository

import (
	"fmt"
	"forum/internal/utils/utils"
	"forum/pkg/models"
	"github.com/jackc/pgx"
	"github.com/pkg/errors"
	"log"
	"strings"
)

type PostRepositoryInterface interface {
	AddPosts(posts models.Posts, threadId int, forumName string) (models.Posts, error)
	ChangePost(updateMessage models.PostUpdate, id int) (models.Post, bool)
	GetAllInfo(params models.FullPostParams, id int) (models.FullPost, bool)
	GetAllPostByThread(id int, limit int, since int, desc bool) ([]models.Post, bool)
	GetPostsTree(id int, limit int, since int, desc bool) ([]models.Post, bool)
	GetPostsParentTree(id int, limit int, since int, desc bool) ([]models.Post, bool)
}

type PostRepository struct {
	DB *pgx.ConnPool
}

func (p PostRepository) ParseRowsToPost(rows *pgx.Rows) ([]models.Post, bool) {
	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created)
		if err != nil {
			log.Println(err)
			rows.Close()
			return posts, false
		}
		posts = append(posts, post)
	}

	rows.Close()
	return posts, true
}

func (p PostRepository) GetPostsTree(id int, limit int, since int, desc bool) ([]models.Post, bool) {
	var rows *pgx.Rows
	var err error

	if since == 0 {
		if desc {
			rows, err = p.DB.Query("GetPostsTreeDesc", id, limit)
		} else {
			rows, err = p.DB.Query("GetPostsTree", id, limit)
		}
	} else {
		if desc {
			rows, err = p.DB.Query("GetPostsTreeSinceDesc", id, since, limit)
		} else {
			rows, err = p.DB.Query("GetPostsTreeSince", id, since, limit)
		}
	}
	if err != nil {
		return []models.Post{}, false
	}

	return p.ParseRowsToPost(rows)
}

const (
	GetPostsTreeDesc      = `SELECT id, parent, author, message, isedited, forum, thread, created FROM dbproject."Post" WHERE thread = $1 ORDER BY path DESC, id DESC LIMIT $2`
	GetPostsTree          = `SELECT id, parent, author, message, isedited, forum, thread, created FROM dbproject."Post" WHERE thread = $1 ORDER BY path, id LIMIT $2`
	GetPostsTreeSinceDesc = `SELECT id, parent, author, message, isedited, forum, thread, created FROM dbproject."Post" WHERE thread = $1 AND path < (SELECT path FROM dbproject."Post" where id = $2) ORDER BY path DESC , id DESC LIMIT $3`
	GetPostsTreeSince     = `SELECT id, parent, author, message, isedited, forum, thread, created FROM dbproject."Post" WHERE thread = $1 AND path > (SELECT path FROM dbproject."Post" where id = $2) ORDER BY path, id LIMIT $3`

	GetPostsParentDesc      = `SELECT id, parent, author, message, isedited, forum, thread, created FROM dbproject."Post" WHERE path[1] IN (SELECT id FROM dbproject."Post" WHERE thread = $1 AND parent = 0 ORDER BY id DESC LIMIT $2) ORDER BY path[1] DESC, path, id`
	GetPostsParent          = `SELECT id, parent, author, message, isedited, forum, thread, created FROM dbproject."Post" WHERE path[1] IN (SELECT id FROM dbproject."Post" WHERE thread = $1 AND parent = 0 ORDER BY id LIMIT $2) ORDER BY path, id`
	GetPostsParentSinceDesc = `SELECT id, parent, author, message, isedited, forum, thread, created FROM dbproject."Post" WHERE path[1] IN (SELECT id FROM dbproject."Post" WHERE thread = $1 AND parent = 0 AND path[1] < (SELECT path[1] FROM dbproject."Post" WHERE id = $2) ORDER BY id DESC LIMIT $3) ORDER BY path[1] DESC, path, id`
	GetPostsParentSince     = `SELECT id, parent, author, message, isedited, forum, thread, created FROM dbproject."Post" WHERE path[1] IN (SELECT id FROM dbproject."Post" WHERE thread = $1 AND parent = 0 AND path[1] > (SELECT path[1] FROM dbproject."Post" WHERE id = $2) ORDER BY id LIMIT $3) ORDER BY path, id`

	GetPostsFlatDesc      = `SELECT id, parent, author, message, isedited, forum, thread, created FROM dbproject."Post" WHERE thread = $1 ORDER BY id DESC LIMIT $2`
	GetPostsFlat          = `SELECT id, parent, author, message, isedited, forum, thread, created FROM dbproject."Post" WHERE thread = $1 ORDER BY id LIMIT $2`
	GetPostsFlatSinceDesc = `SELECT id, parent, author, message, isedited, forum, thread, created FROM dbproject."Post" WHERE thread = $1 AND id < $2 ORDER BY id DESC LIMIT $3`
	GetPostsFlatSince     = `SELECT id, parent, author, message, isedited, forum, thread, created FROM dbproject."Post" WHERE thread = $1 AND id > $2 ORDER BY id LIMIT $3`

	SelectPostInfo       = `SELECT id, parent, author, message, isedited, forum, thread, created FROM dbproject."Post" WHERE id = $1`
	SelectPostInfoUser   = `SELECT nickname, fullname, about, email FROM dbproject."User" WHERE nickname = $1`
	SelectPostInfoThread = `SELECT id, title, author, forum, message, votes, slug, created FROM dbproject."Thread" WHERE id = $1`
	SelectPostInfoForum  = `SELECT title, "user", slug, posts, threads FROM dbproject."Forum" WHERE slug = $1`

	UpdatePost    = `UPDATE dbproject."Post" SET message = COALESCE(NULLIF($1, ''), message), isedited = CASE WHEN $1 = '' OR message=$1 THEN isedited else true end WHERE id = $2 RETURNING id, parent, author, message, isedited, forum, thread, created`
	GetPostParent = `SELECT id FROM dbproject."Post" WHERE thread = $1 AND id = $2`
)

func (p PostRepository) GetPostsParentTree(id int, limit int, since int, desc bool) ([]models.Post, bool) {
	var rows *pgx.Rows
	var err error

	if since == 0 {
		if desc {
			rows, err = p.DB.Query("GetPostsParentDesc", id, limit)
		} else {
			rows, err = p.DB.Query("GetPostsParent", id, limit)
		}
	} else {
		if desc {
			rows, err = p.DB.Query("GetPostsParentSinceDesc", id, since, limit)
		} else {
			rows, err = p.DB.Query("GetPostsParentSince", id, since, limit)
		}
	}

	if err != nil {
		return []models.Post{}, false
	}

	return p.ParseRowsToPost(rows)
}

func (p PostRepository) GetAllPostByThread(id int, limit int, since int, desc bool) ([]models.Post, bool) {
	var rows *pgx.Rows
	var err error

	if since == 0 {
		if desc {
			rows, err = p.DB.Query("GetPostsFlatDesc", id, limit)
		} else {
			rows, err = p.DB.Query("GetPostsFlat", id, limit)
		}
	} else {
		if desc {
			rows, err = p.DB.Query("GetPostsFlatSinceDesc", id, since, limit)
		} else {
			rows, err = p.DB.Query("GetPostsFlatSince", id, since, limit)
		}
	}

	if err != nil {
		return []models.Post{}, false
	}

	return p.ParseRowsToPost(rows)
}

func (p PostRepository) GetAllInfo(params models.FullPostParams, id int) (models.FullPost, bool) {
	var info models.FullPost

	post := models.Post{}
	user := models.User{}
	forum := models.Forum{}
	thread := models.Thread{}

	err := p.DB.QueryRow("SelectPostInfo", id).
		Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created)
	if err != nil {
		log.Println(err)
		return models.FullPost{}, false
	}
	info.Post = &post

	if params.User {
		err = p.DB.QueryRow("SelectPostInfoUser", post.Author).
			Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
		if err != nil {
			return models.FullPost{}, false
		}
		info.Author = &user
	}

	if params.Thread {
		err = p.DB.QueryRow("SelectPostInfoThread", post.Thread).
			Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
		if err != nil {
			return models.FullPost{}, false
		}
		if utils.IsValidUUID(thread.Slug) {
			thread.Slug = ""
		}

		info.Thread = &thread
	}

	if params.Forum {
		err = p.DB.QueryRow("SelectPostInfoForum", post.Forum).
			Scan(&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads)
		if err != nil {
			return models.FullPost{}, false
		}
		info.Forum = &forum
	}

	return info, true
}

func (p PostRepository) ChangePost(updateMessage models.PostUpdate, id int) (models.Post, bool) {
	var post models.Post
	err := p.DB.QueryRow("UpdatePost", updateMessage.Message, id).
		Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created)

	if err != nil {
		log.Println(err)
		return models.Post{}, false
	}

	return post, true
}

func (p PostRepository) AddPosts(posts models.Posts, threadId int, forumName string) (models.Posts, error) {
	var insertedPosts models.Posts

	var sqlValues []interface{}
	sqlQuery := `INSERT INTO dbproject."Post" (PARENT, AUTHOR, MESSAGE, FORUM, THREAD) VALUES `

	if len(posts) == 0 {
		return models.Posts{}, nil
	}

	for i, post := range posts {
		if post.Parent != 0 {
			id := -1
			err := p.DB.QueryRow("GetPostParent", threadId, post.Parent).Scan(&id)
			if err == pgx.ErrNoRows {
				return nil, errors.New("Cant Find Parent")
			}
		}

		sqlValuesString := fmt.Sprintf("($%d, $%d, $%d, $%d, $%d),", i*5+1, i*5+2, i*5+3, i*5+4, i*5+5)

		sqlQuery += sqlValuesString
		sqlValues = append(sqlValues, post.Parent, post.Author, post.Message, forumName, threadId)
	}

	sqlQuery = strings.TrimSuffix(sqlQuery, ",")
	sqlQuery += ` RETURNING id, parent, author, message, isedited, forum, thread, created;`

	rows, err := p.DB.Query(sqlQuery, sqlValues...)

	if err != nil {
		return nil, errors.New("Cant Find Parent")
	}

	defer rows.Close()

	for rows.Next() {
		post := models.Post{}

		err := rows.Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created)
		if err != nil || post.Author == "" {
			return nil, err
		}

		insertedPosts = append(insertedPosts, post)
	}

	if len(insertedPosts) == 0 {
		return nil, errors.New(models.ErrUserUnknown)
	}

	return insertedPosts, nil
}
