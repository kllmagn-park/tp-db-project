package delivery

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"tp-db-project/internal/models"
	"tp-db-project/pkg/errs"
	"tp-db-project/repository"

	"github.com/valyala/fasthttp"
)

type ForumService struct {
	forumRepo repository.ForumRepo
}

func MakeForumService(forumRepo repository.ForumRepo) ForumService {
	return ForumService{forumRepo: forumRepo}
}


func (s *ForumService) Create(cont *fasthttp.RequestCtx) {
	cont.SetContentType("application/json")
	uctx := cont.UserValue("cont").(context.Context)

	var forumCreate models.CreateForumForm
	_ = json.Unmarshal(cont.PostBody(), &forumCreate)

	forum, err := s.forumRepo.Create(uctx, &forumCreate)

	if err != nil {
		if errors.Is(err, errs.ErrorUserDoesNotExist) {
			body, _ := json.Marshal(GetErrorMessage(err))
			cont.SetBody(body)
			cont.SetStatusCode(fasthttp.StatusNotFound)
		} else if errors.Is(err, errs.ErrorForumAlreadyExist) {
			body, _ := json.Marshal(forum)
			cont.SetBody(body)
			cont.SetStatusCode(fasthttp.StatusConflict)
		}

		return
	}

	body, _ := json.Marshal(forum)
	cont.SetBody(body)
	cont.SetStatusCode(fasthttp.StatusCreated)
}


func (s *ForumService) Get(cont *fasthttp.RequestCtx) {
	cont.SetContentType("application/json")
	uctx := cont.UserValue("cont").(context.Context)
	slug := cont.UserValue("slug").(string)

	forum, err := s.forumRepo.Get(uctx, slug)
	if err != nil {
		body, _ := json.Marshal(GetErrorMessage(err))
		cont.SetBody(body)
		cont.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	body, _ := json.Marshal(forum)
	cont.SetBody(body)
	cont.SetStatusCode(fasthttp.StatusOK)
}


func (s *ForumService) GetUsers(cont *fasthttp.RequestCtx) {
	cont.SetContentType("application/json")
	uctx := cont.UserValue("cont").(context.Context)
	slug := cont.UserValue("slug").(string)

	limit, err := strconv.Atoi(string(cont.QueryArgs().Peek("limit")))
	if err != nil {
		limit = 100
	}

	since := string(cont.QueryArgs().Peek("since"))

	desc, err := strconv.ParseBool(string(cont.QueryArgs().Peek("desc")))
	if err != nil {
		desc = false
	}

	forumUsers := &models.GetUsersForm{
		Slug:  slug,
		Limit: int32(limit),
		Since: since,
		Desc:  desc,
	}

	users, err := s.forumRepo.GetUsers(uctx, forumUsers)
	if err != nil {
		body, _ := json.Marshal(GetErrorMessage(err))
		cont.SetBody(body)
		cont.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	body, _ := json.Marshal(users)
	cont.SetStatusCode(fasthttp.StatusOK)
	cont.SetBody(body)
}


func (s *ForumService) GetThreads(cont *fasthttp.RequestCtx) {
	cont.SetContentType("application/json")
	uctx := cont.UserValue("cont").(context.Context)
	slug := cont.UserValue("slug").(string)

	limit, err := strconv.Atoi(string(cont.QueryArgs().Peek("limit")))
	if err != nil {
		limit = 100
	}

	desc, err := strconv.ParseBool(string(cont.QueryArgs().Peek("desc")))
	if err != nil {
		desc = false
	}

	forumThreads := &models.GetThreadsForm{
		Limit: int32(limit),
		Since: string(cont.QueryArgs().Peek("since")),
		Desc:  desc,
	}

	threads, err := s.forumRepo.GetThreads(uctx, slug, forumThreads)
	if err != nil {
		body, _ := json.Marshal(GetErrorMessage(err))
		cont.SetBody(body)
		cont.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	body, _ := json.Marshal(threads)
	cont.SetBody(body)
	cont.SetStatusCode(fasthttp.StatusOK)
}
