package postgres

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
	"tp-db-project/internal/models"
	"tp-db-project/pkg/errs"
	"tp-db-project/pkg/queries"
	"tp-db-project/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ()

type PostRepository struct {
	conn *pgxpool.Pool
}

func InitPostRepository(db *pgxpool.Pool) repository.PostRepo {
	return &PostRepository{conn: db}
}

func isIn(arr *[]string, find string) bool {
	for _, str := range *arr {
		if str == find {
			return true
		}
	}

	return false
}

func (s *PostRepository) Get(cont context.Context, id int64, getSettings *models.GetPostReqForm) (*models.GetPostForm, error) {
	var post models.Post
	err := s.conn.QueryRow(cont, queries.GetPostCommand, id).Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created)
	if err != nil {
		return nil, errs.ErrorPostDoesNotExist
	}

	var postResult models.GetPostForm
	postResult.Post = &post

	if isIn(&getSettings.Related, models.RelatedUser) {
		author := &models.User{}
		_ = s.conn.QueryRow(cont, queries.GetPostAuthorCommand, post.Author).Scan(&author.Nickname, &author.Fullname, &author.About, &author.Email)
		postResult.Author = author
	}
	if isIn(&getSettings.Related, models.RelatedThread) {
		thread := &models.Thread{}
		_ = s.conn.QueryRow(cont, queries.GetPostThreadCommand, post.Thread).Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
		postResult.Thread = thread
	}
	if isIn(&getSettings.Related, models.RelatedForum) {
		forum := &models.Forum{}
		_ = s.conn.QueryRow(cont, queries.GetPostForumCommand, post.Forum).Scan(&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads)
		postResult.Forum = forum
	}

	return &postResult, nil
}

func (s *PostRepository) Update(cont context.Context, id int64, updateDate *models.PostUpdate) (*models.Post, error) {
	var post models.Post
	err := s.conn.QueryRow(cont, queries.GetPostCommand, id).Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created)
	if err != nil {
		return nil, errs.ErrorPostDoesNotExist
	}

	if updateDate.Message == "" || updateDate.Message == post.Message {
		return &post, nil
	}

	_ = s.conn.QueryRow(cont, queries.UpdatePostCommand, updateDate.Message, id)
	post.Message = updateDate.Message
	post.IsEdited = true

	return &post, nil
}

func (s *PostRepository) CheckParentAndAuthor(cont context.Context, post *models.CreatePostForm) error {
	if post.Parent != 0 {
		_, err := s.conn.Exec(cont, queries.GetPostCommand, post.Parent)
		if err != nil {
			return err
		}
	}
	_, err := s.conn.Exec(cont, queries.GetUserByNicknameCommand, post.Author)
	return err
}

func (s *PostRepository) Create(cont context.Context, slugOrId string, posts *[]models.CreatePostForm) (*[]models.Post, error) {
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
	if len(*posts) == 0 {
		postsToRet := make([]models.Post, 0)
		return &postsToRet, nil
	}
	if len(*posts) == 0 {
		emptyPosts := make([]models.Post, 0)
		return &emptyPosts, nil
	}
	if s.CheckParentAndAuthor(cont, &(*posts)[0]) != nil {
		return nil, errs.ErrorParentPostDoesNotExist
	}
	command := strings.Builder{}
	command.WriteString("INSERT INTO posts (parent, author, message, forum, thread, created) VALUES ")
	argsForCommand := make([]interface{}, 0, len(*posts))
	postsToReturn := make([]models.Post, 0, len(*posts))
	createdTime := time.Unix(0, time.Now().UnixNano()/1e6*1e6)
	for ind, post := range *posts {
		if post.Parent != 0 {
			var parentPost models.Post
			err = s.conn.QueryRow(cont, queries.GetPostCommand, post.Parent).Scan(&parentPost.Id, &parentPost.Parent, &parentPost.Author, &parentPost.Message, &parentPost.IsEdited, &parentPost.Forum, &parentPost.Thread, &parentPost.Created)
			if err != nil || parentPost.Thread != thread.Id {
				return nil, errs.ErrorParentPostDoesNotExist
			}
		}
		var author models.User
		err = s.conn.QueryRow(cont, queries.GetUserByNicknameCommand, post.Author).Scan(&author.Nickname, &author.Fullname, &author.About, &author.Email)
		if err != nil {
			return nil, errs.ErrorAuthorDoesNotExist
		}
		sixInd := ind * 6
		postsToReturn = append(postsToReturn, models.Post{Parent: post.Parent, Author: post.Author, Message: post.Message, Forum: thread.Forum, Thread: thread.Id, Created: createdTime})
		fmt.Fprintf(&command, "($%d, $%d, $%d, $%d, $%d, $%d),", sixInd+1, sixInd+2, sixInd+3, sixInd+4, sixInd+5, sixInd+6)
		argsForCommand = append(argsForCommand, post.Parent, post.Author, post.Message, thread.Forum, thread.Id, createdTime)
	}
	qs := command.String()
	qs = qs[:len(qs)-1] + " RETURNING id"
	rows, err := s.conn.Query(cont, qs, argsForCommand...)
	if err != nil {
		return nil, errs.ErrorAuthorDoesNotExist
	}
	defer rows.Close()
	for ind := 0; rows.Next(); ind++ {
		err = rows.Scan(&postsToReturn[ind].Id)
		if err != nil {
			return nil, errs.ErrorAuthorDoesNotExist
		}
	}
	return &postsToReturn, nil
}
