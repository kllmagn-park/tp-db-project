package postgres

import (
	"context"
	"strconv"
	"tp-db-project/internal/models"
	"tp-db-project/pkg/errs"
	"tp-db-project/pkg/queries"
	"tp-db-project/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ThreadRepository struct {
	conn *pgxpool.Pool
}

func InitThreadRepository(db *pgxpool.Pool) repository.ThreadRepo {
	return &ThreadRepository{conn: db}
}

func (s *ThreadRepository) Create(cont context.Context, forumSlug string, thread *models.CreateThreadForm) (*models.Thread, error) {
	var forum models.Forum
	err := s.conn.QueryRow(cont, queries.GetForumCommand, forumSlug).Scan(&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads)
	if err != nil {
		return nil, errs.ErrorNoAuthorOrForum
	}

	var user models.User
	err = s.conn.QueryRow(cont, queries.GetUserByNicknameCommand, thread.Author).Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
	if err != nil {
		return nil, errs.ErrorNoAuthorOrForum
	}
	thread.Forum = forum.Slug

	if thread.Slug != "" {
		var threadAlreadyExist models.Thread
		err = s.conn.QueryRow(cont, queries.GetThreadBySlugCommand, thread.Slug).Scan(&threadAlreadyExist.Id, &threadAlreadyExist.Title, &threadAlreadyExist.Author, &threadAlreadyExist.Forum, &threadAlreadyExist.Message, &threadAlreadyExist.Votes, &threadAlreadyExist.Slug, &threadAlreadyExist.Created)
		if err == nil {
			return &threadAlreadyExist, errs.ErrorThreadAlreadyExist
		}
	}

	var id int32
	err = s.conn.QueryRow(cont, queries.CreateThreadCommand, thread.Title, thread.Author, thread.Message, thread.Created, thread.Slug, thread.Forum).Scan(&id)
	if err != nil {
		threadAlreadyExist, _ := s.Get(cont, thread.Slug)
		return threadAlreadyExist, errs.ErrorThreadAlreadyExist
	}

	threadToReturn := &models.Thread{
		Id:      id,
		Title:   thread.Title,
		Author:  thread.Author,
		Forum:   thread.Forum,
		Message: thread.Message,
		Votes:   0,
		Slug:    thread.Slug,
		Created: thread.Created,
	}

	return threadToReturn, nil
}

func (s *ThreadRepository) Get(cont context.Context, slugOrId string) (*models.Thread, error) {
	var thread models.Thread
	id, err := strconv.Atoi(slugOrId)

	if err != nil {
		err = s.conn.QueryRow(cont, queries.GetThreadBySlugCommand, slugOrId).Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
	} else {
		err = s.conn.QueryRow(cont, queries.GetThreadByIdCommand, id).Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
	}

	if err != nil {
		return nil, errs.ErrorThreadDoesNotExist
	}

	return &thread, nil
}

func (s *ThreadRepository) Update(cont context.Context, slugOrId string, updateData *models.ThreadUpdate) (*models.Thread, error) {
	thread, err := s.Get(cont, slugOrId)
	if err != nil {
		return nil, errs.ErrorThreadDoesNotExist
	}

	if updateData.Message == "" {
		updateData.Message = thread.Message
	} else {
		thread.Message = updateData.Message
	}
	if updateData.Title == "" {
		updateData.Title = thread.Title
	} else {
		thread.Title = updateData.Title
	}

	_, _ = s.conn.Exec(cont, queries.UpdateThreadByIdCommand, updateData.Title, updateData.Message, thread.Id)

	return thread, nil
}

func (s *ThreadRepository) GetPosts(cont context.Context, slugOrId string, getSettings *models.GetPostsForm) (*[]models.Post, error) {
	thread, err := s.Get(cont, slugOrId)
	if err != nil {
		return nil, errs.ErrorThreadDoesNotExist
	}
	var rows pgx.Rows
	if getSettings.Sort == models.Flat {
		if getSettings.Desc {
			if getSettings.Since != -1 {
				rows, _ = s.conn.Query(cont, queries.GetPostsOnThreadFlatDescCommand, thread.Id, getSettings.Since, getSettings.Limit)
			} else {
				rows, _ = s.conn.Query(cont, queries.GetPostsOnThreadFlatDescWithoutSinceCommand, thread.Id, getSettings.Limit)
			}
		} else {
			if getSettings.Since != -1 {
				rows, _ = s.conn.Query(cont, queries.GetPostsOnThreadFlatCommand, thread.Id, getSettings.Since, getSettings.Limit)
			} else {
				rows, _ = s.conn.Query(cont, queries.GetPostsOnThreadFlatWithoutSinceCommand, thread.Id, getSettings.Limit)
			}
		}
	} else if getSettings.Sort == models.Tree {
		if getSettings.Desc {
			if getSettings.Since != -1 {
				rows, _ = s.conn.Query(cont, queries.GetPostsOnThreadTreeDescCommand, thread.Id, getSettings.Since, getSettings.Limit)
			} else {
				rows, _ = s.conn.Query(cont, queries.GetPostsOnThreadTreeDescWithoutSinceCommand, thread.Id, getSettings.Limit)
			}
		} else {
			if getSettings.Since != -1 {
				rows, _ = s.conn.Query(cont, queries.GetPostsOnThreadTreeCommand, thread.Id, getSettings.Since, getSettings.Limit)
			} else {
				rows, _ = s.conn.Query(cont, queries.GetPostsOnThreadTreeWithoutSinceCommand, thread.Id, getSettings.Limit)
			}
		}
	} else if getSettings.Sort == models.ParentTree {
		if getSettings.Desc {
			if getSettings.Since > 0 {
				rows, _ = s.conn.Query(cont, queries.GetPostsOnThreadParentTreeDescWithSinceCommand, thread.Id, getSettings.Since, getSettings.Limit)
			} else {
				rows, _ = s.conn.Query(cont, queries.GetPostsOnThreadParentTreeDescWithoutSinceCommand, thread.Id, getSettings.Limit)
			}
		} else {
			if getSettings.Since != -1 {
				rows, _ = s.conn.Query(cont, queries.GetPostsOnThreadParentTreeCommand, thread.Id, getSettings.Since, getSettings.Limit)
			} else {
				rows, _ = s.conn.Query(cont, queries.GetPostsOnThreadParentTreeWithoutSinceCommand, thread.Id, getSettings.Limit)
			}
		}
	}
	defer rows.Close()
	posts := make([]models.Post, 0, rows.CommandTag().RowsAffected())
	for rows.Next() {
		post := models.Post{}
		_ = rows.Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created)
		posts = append(posts, post)
	}

	return &posts, nil
}
