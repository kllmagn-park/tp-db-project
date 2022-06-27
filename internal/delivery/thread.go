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

type ThreadService struct {
	threadRepo repository.ThreadRepo
}

func MakeThreadService(threadRepo repository.ThreadRepo) ThreadService {
	return ThreadService{threadRepo: threadRepo}
}


func (s *ThreadService) Create(cont *fasthttp.RequestCtx) {
	cont.SetContentType("application/json")
	uctx := cont.UserValue("cont").(context.Context)
	slug := cont.UserValue("slug").(string)

	var threadCreate models.CreateThreadForm
	_ = json.Unmarshal(cont.PostBody(), &threadCreate)

	thread, err := s.threadRepo.Create(uctx, slug, &threadCreate)
	if err != nil {
		if errors.Is(err, errs.ErrorNoAuthorOrForum) {
			body, _ := json.Marshal(GetErrorMessage(err))
			cont.SetBody(body)
			cont.SetStatusCode(fasthttp.StatusNotFound)
			return
		} else if errors.Is(err, errs.ErrorThreadAlreadyExist) {
			body, _ := json.Marshal(thread)
			cont.SetBody(body)
			cont.SetStatusCode(fasthttp.StatusConflict)
			return
		}
	}

	body, _ := json.Marshal(thread)
	cont.SetBody(body)
	cont.SetStatusCode(fasthttp.StatusCreated)
}


func (s *ThreadService) Get(cont *fasthttp.RequestCtx) {
	cont.SetContentType("application/json")
	uctx := cont.UserValue("cont").(context.Context)
	slugOrId := cont.UserValue("slug_or_id").(string)

	thread, err := s.threadRepo.Get(uctx, slugOrId)
	if err != nil {
		body, _ := json.Marshal(GetErrorMessage(err))
		cont.SetBody(body)
		cont.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	body, _ := json.Marshal(thread)
	cont.SetBody(body)
	cont.SetStatusCode(fasthttp.StatusOK)
}


func (s *ThreadService) Update(cont *fasthttp.RequestCtx) {
	cont.SetContentType("application/json")
	uctx := cont.UserValue("cont").(context.Context)
	slugOrId := cont.UserValue("slug_or_id").(string)

	var threadUpdate models.ThreadUpdate
	_ = json.Unmarshal(cont.PostBody(), &threadUpdate)

	thread, err := s.threadRepo.Update(uctx, slugOrId, &threadUpdate)
	if err != nil {
		body, _ := json.Marshal(GetErrorMessage(err))
		cont.SetBody(body)
		cont.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	body, _ := json.Marshal(thread)
	cont.SetBody(body)
	cont.SetStatusCode(fasthttp.StatusOK)
}


func (s *ThreadService) GetPosts(cont *fasthttp.RequestCtx) {
	cont.SetContentType("application/json")
	uctx := cont.UserValue("cont").(context.Context)
	slugOrId := cont.UserValue("slug_or_id").(string)

	limit, err := strconv.Atoi(string(cont.QueryArgs().Peek("limit")))
	if err != nil {
		limit = 100
	}

	since, err := strconv.Atoi(string(cont.QueryArgs().Peek("since")))
	if err != nil {
		since = -1
	}

	sort := string(cont.QueryArgs().Peek("sort"))
	if sort == "" {
		sort = models.Flat
	}

	desc, err := strconv.ParseBool(string(cont.QueryArgs().Peek("desc")))
	if err != nil {
		desc = false
	}

	threadGetPosts := &models.GetPostsForm{
		Limit: int32(limit),
		Since: int64(since),
		Sort:  sort,
		Desc:  desc,
	}

	posts, err := s.threadRepo.GetPosts(uctx, slugOrId, threadGetPosts)

	if err != nil {
		body, _ := json.Marshal(GetErrorMessage(err))
		cont.SetBody(body)
		cont.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	body, _ := json.Marshal(posts)
	cont.SetBody(body)
	cont.SetStatusCode(fasthttp.StatusOK)
}
