package postgres

import (
	"context"
	"tp-db-project/internal/models"
	"tp-db-project/pkg/errs"
	"tp-db-project/pkg/queries"
	"tp-db-project/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ForumRepository struct {
	conn *pgxpool.Pool
}

func InitForumRepository(db *pgxpool.Pool) repository.ForumRepo {
	return &ForumRepository{conn: db}
}

func (s *ForumRepository) Create(cont context.Context, forum *models.CreateForumForm) (*models.Forum, error) {
	var user models.User
	err := s.conn.QueryRow(cont, queries.GetUserByNicknameCommand, forum.User).Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
	if err != nil {
		return nil, errs.ErrorUserDoesNotExist
	}

	_, err = s.conn.Exec(cont, queries.CreateForumCommand, forum.Title, user.Nickname, forum.Slug)
	if err != nil {
		forumAlreadyExist, _ := s.Get(cont, forum.Slug)
		return forumAlreadyExist, errs.ErrorForumAlreadyExist
	}

	forumToReturn := &models.Forum{
		Title:   forum.Title,
		User:    user.Nickname,
		Slug:    forum.Slug,
		Posts:   0,
		Threads: 0,
	}

	return forumToReturn, nil
}

func (s *ForumRepository) Get(cont context.Context, slug string) (*models.Forum, error) {
	var forum models.Forum
	err := s.conn.QueryRow(cont, queries.GetForumCommand, slug).Scan(&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads)
	if err != nil {
		return nil, errs.ErrorForumDoesNotExist
	}
	return &forum, nil
}

func (s *ForumRepository) GetUsers(cont context.Context, getSettings *models.GetUsersForm) (*[]models.User, error) {
	var err error
	var rows pgx.Rows
	_, err = s.Get(cont, getSettings.Slug)
	if err != nil {
		return nil, errs.ErrorForumDoesNotExist
	}
	if getSettings.Desc {
		if getSettings.Since == "" {
			rows, err = s.conn.Query(cont, queries.GetUsersOnForumWithoutSinceDescCommand, getSettings.Slug, getSettings.Limit)
		} else {
			rows, err = s.conn.Query(cont, queries.GetUsersOnForumDescCommand, getSettings.Slug, getSettings.Since, getSettings.Limit)
		}
	} else {
		if getSettings.Since == "" {
			rows, err = s.conn.Query(cont, queries.GetUsersOnForumWithoutSinceCommand, getSettings.Slug, getSettings.Limit)
		} else {
			rows, err = s.conn.Query(cont, queries.GetUsersOnForumCommand, getSettings.Slug, getSettings.Since, getSettings.Limit)
		}
	}
	if err != nil {
		return nil, errs.ErrorForumDoesNotExist
	}
	defer rows.Close()
	users := make([]models.User, 0, rows.CommandTag().RowsAffected())
	for rows.Next() {
		user := models.User{}
		err = rows.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
		if err != nil {
			return nil, errs.ErrorForumDoesNotExist
		}
		users = append(users, user)
	}
	return &users, nil
}

func (s *ForumRepository) GetThreads(cont context.Context, slug string, getSettings *models.GetThreadsForm) (*[]models.Thread, error) {
	var rows pgx.Rows
	_, err := s.Get(cont, slug)
	if err != nil {
		return nil, errs.ErrorForumDoesNotExist
	}
	if getSettings.Desc {
		if getSettings.Since == "" {
			rows, err = s.conn.Query(cont, queries.GetThreadsOnForumWithoutSinceDescCommand, slug, getSettings.Limit)
		} else {
			rows, err = s.conn.Query(cont, queries.GetThreadsOnForumDescCommand, slug, getSettings.Since, getSettings.Limit)
		}
	} else {
		if getSettings.Since == "" {
			rows, err = s.conn.Query(cont, queries.GetThreadsOnForumWithoutSinceCommand, slug, getSettings.Limit)
		} else {
			rows, err = s.conn.Query(cont, queries.GetThreadsOnForumCommand, slug, getSettings.Since, getSettings.Limit)
		}
	}
	if err != nil {
		return nil, errs.ErrorForumDoesNotExist
	}
	defer rows.Close()
	threads := make([]models.Thread, 0, rows.CommandTag().RowsAffected())
	for rows.Next() {
		thread := models.Thread{}
		err = rows.Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
		if err != nil {
			return nil, errs.ErrorForumDoesNotExist
		}

		threads = append(threads, thread)
	}
	return &threads, nil
}
