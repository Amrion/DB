package repository

import (
	"forum/pkg/forum/repository"
	repository2 "forum/pkg/post/repository"
	repository3 "forum/pkg/service/repository"
	repository4 "forum/pkg/thread/repository"
	"forum/pkg/user/repostitory"
	"github.com/jackc/pgx"
)

type Postgres struct {
	DB *pgx.ConnPool
}

func NewPostgres() (*Postgres, error) {
	conf := pgx.ConnConfig{
		User:                 "docker",
		Database:             "docker",
		Password:             "docker",
		PreferSimpleProtocol: false,
	}

	poolConf := pgx.ConnPoolConfig{
		ConnConfig:     conf,
		MaxConnections: 100,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	}
	db, err := pgx.NewConnPool(poolConf)
	if err != nil {
		return nil, err
	}
	return &Postgres{
		DB: db,
	}, nil
}

func (p *Postgres) GetPostgres() *pgx.ConnPool {
	return p.DB
}

func (p *Postgres) Close() error {
	p.DB.Close()
	return nil
}

func (p *Postgres) ProcedureRequests() error {
	//forum
	if _, err := p.DB.Prepare("SelectUsersByForumDesc", repository.SelectUsersByForumDesc); err != nil {
		return err
	}
	if _, err := p.DB.Prepare("SelectUsersByForumSince", repository.SelectUsersByForumSince); err != nil {
		return err
	}
	if _, err := p.DB.Prepare("SelectUsersByForumSinceDesc", repository.SelectUsersByForumSinceDesc); err != nil {
		return err
	}
	if _, err := p.DB.Prepare("SelectForum", repository.SelectForum); err != nil {
		return err
	}
	if _, err := p.DB.Prepare("SelectUsersByForum", repository.SelectUsersByForum); err != nil {
		return err
	}
	if _, err := p.DB.Prepare("InsertForum", repository.InsertForum); err != nil {
		return err
	}

	//post
	if _, err := p.DB.Prepare("GetPostsTreeDesc", repository2.GetPostsTreeDesc); err != nil {
		return err
	}
	if _, err := p.DB.Prepare("GetPostsTree", repository2.GetPostsTree); err != nil {
		return err
	}
	if _, err := p.DB.Prepare("GetPostsFlat", repository2.GetPostsFlat); err != nil {
		return err
	}
	if _, err := p.DB.Prepare("GetPostParent", repository2.GetPostParent); err != nil {
		return err
	}
	if _, err := p.DB.Prepare("GetPostsFlatDesc", repository2.GetPostsFlatDesc); err != nil {
		return err
	}
	if _, err := p.DB.Prepare("GetPostsFlatSince", repository2.GetPostsFlatSince); err != nil {
		return err
	}
	if _, err := p.DB.Prepare("GetPostsFlatSinceDesc", repository2.GetPostsFlatSinceDesc); err != nil {
		return err
	}
	if _, err := p.DB.Prepare("GetPostsParent", repository2.GetPostsParent); err != nil {
		return err
	}
	if _, err := p.DB.Prepare("GetPostsTreeSinceDesc", repository2.GetPostsTreeSinceDesc); err != nil {
		return err
	}
	if _, err := p.DB.Prepare("GetPostsTreeSince", repository2.GetPostsTreeSince); err != nil {
		return err
	}
	if _, err := p.DB.Prepare("GetPostsParentDesc", repository2.GetPostsParentDesc); err != nil {
		return err
	}
	if _, err := p.DB.Prepare("GetPostsParentSinceDesc", repository2.GetPostsParentSinceDesc); err != nil {
		return err
	}
	if _, err := p.DB.Prepare("GetPostsParentSince", repository2.GetPostsParentSince); err != nil {
		return err
	}
	if _, err := p.DB.Prepare("SelectPostInfo", repository2.SelectPostInfo); err != nil {
		return err
	}
	if _, err := p.DB.Prepare("SelectPostInfoUser", repository2.SelectPostInfoUser); err != nil {
		return err
	}
	if _, err := p.DB.Prepare("SelectPostInfoThread", repository2.SelectPostInfoThread); err != nil {
		return err
	}
	if _, err := p.DB.Prepare("SelectPostInfoForum", repository2.SelectPostInfoForum); err != nil {
		return err
	}
	if _, err := p.DB.Prepare("UpdatePost", repository2.UpdatePost); err != nil {
		return err
	}

	//service
	if _, err := p.DB.Prepare("CleanDB", repository3.CleanDB); err != nil {
		return err
	}
	if _, err := p.DB.Prepare("StatusPost", repository3.StatusPost); err != nil {
		return err
	}
	if _, err := p.DB.Prepare("StatusUser", repository3.StatusUser); err != nil {
		return err
	}
	if _, err := p.DB.Prepare("StatusForum", repository3.StatusForum); err != nil {
		return err
	}
	if _, err := p.DB.Prepare("StatusThread", repository3.StatusThread); err != nil {
		return err
	}

	//thread
	if _, err := p.DB.Prepare("SelectThreadIdBySlug", repository4.SelectThreadIdBySlug); err != nil {
		return err
	}
	if _, err := p.DB.Prepare("InsertVote", repository4.InsertVote); err != nil {
		return err
	}
	if _, err := p.DB.Prepare("UpdateVote", repository4.UpdateVote); err != nil {
		return err
	}
	if _, err := p.DB.Prepare("UpdateThreadId", repository4.UpdateThreadId); err != nil {
		return err
	}
	if _, err := p.DB.Prepare("UpdateThreadSlug", repository4.UpdateThreadSlug); err != nil {
		return err
	}
	if _, err := p.DB.Prepare("SelectThreadInfoBySlug", repository4.SelectThreadInfoBySlug); err != nil {
		return err
	}
	if _, err := p.DB.Prepare("SelectThreadInfoById", repository4.SelectThreadInfoById); err != nil {
		return err
	}
	if _, err := p.DB.Prepare("SelectThreadDesc", repository4.SelectThreadDesc); err != nil {
		return err
	}
	if _, err := p.DB.Prepare("SelectThread", repository4.SelectThread); err != nil {
		return err
	}
	if _, err := p.DB.Prepare("SelectThreadSinceDesc", repository4.SelectThreadSinceDesc); err != nil {
		return err
	}
	if _, err := p.DB.Prepare("SelectThreadSince", repository4.SelectThreadSince); err != nil {
		return err
	}
	if _, err := p.DB.Prepare("InsertThread", repository4.InsertThread); err != nil {
		return err
	}

	//user
	if _, err := p.DB.Prepare("InsertUser", repostitory.InsertUser); err != nil {
		return err
	}
	if _, err := p.DB.Prepare("SelectUser", repostitory.SelectUser); err != nil {
		return err
	}
	if _, err := p.DB.Prepare("UpdateUser", repostitory.UpdateUser); err != nil {
		return err
	}
	if _, err := p.DB.Prepare("SelectUserByNick", repostitory.SelectUserByNick); err != nil {
		return err
	}

	return nil
}
