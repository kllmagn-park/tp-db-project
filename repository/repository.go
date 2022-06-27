package repository

import (
	"context"
	"tp-db-project/internal/models"
)

type UserRepo interface {
	Create(cont context.Context, user *models.User) (*[]models.User, error)
	Update(cont context.Context, nickname string, updateData *models.UserUpdate) (*models.User, error)
	Get(cont context.Context, nicknameOrEmail string) (*models.User, error)
}

type ForumRepo interface {
	Create(cont context.Context, forum *models.CreateForumForm) (*models.Forum, error)
	Get(cont context.Context, slug string) (*models.Forum, error)
	GetUsers(cont context.Context, getSettings *models.GetUsersForm) (*[]models.User, error)                   
	GetThreads(cont context.Context, slug string, getSettings *models.GetThreadsForm) (*[]models.Thread, error)
}

type PostRepo interface {
	Get(cont context.Context, id int64, getSettings *models.GetPostReqForm) (*models.GetPostForm, error)
	Update(cont context.Context, id int64, updateDate *models.PostUpdate) (*models.Post, error)
	Create(cont context.Context, slugOrId string, posts *[]models.CreatePostForm) (*[]models.Post, error)
}

type ThreadRepo interface {
	Create(cont context.Context, forumSlug string, thread *models.CreateThreadForm) (*models.Thread, error)
	Get(cont context.Context, slugOrId string) (*models.Thread, error)
	Update(cont context.Context, slugOrId string, updateData *models.ThreadUpdate) (*models.Thread, error)
	GetPosts(cont context.Context, slugOrId string, getSettings *models.GetPostsForm) (*[]models.Post, error) 
}

type VoteRepo interface {
	Create(cont context.Context, slugOrId string, vote *models.CreateVoteForm) (*models.Thread, error) 
}

type StatusRepo interface {
	GetStatus(cont context.Context) (*models.Status, error)
	Clear(cont context.Context) error
}
