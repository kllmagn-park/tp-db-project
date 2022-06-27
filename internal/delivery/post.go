package delivery

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"tp-db-project/internal/models"
	"tp-db-project/pkg/errs"
	"tp-db-project/repository"

	"github.com/valyala/fasthttp"
)

type PostService struct {
	postRepo repository.PostRepo
}

func MakePostService(postRepo repository.PostRepo) PostService {
	return PostService{postRepo: postRepo}
}


func (s *PostService) Create(cont *fasthttp.RequestCtx) {
	cont.SetContentType("application/json")
	uctx := cont.UserValue("cont").(context.Context)

	slugOrId := cont.UserValue("slug_or_id").(string)
	postsCreate := make([]models.CreatePostForm, 0)

	_ = json.Unmarshal(cont.PostBody(), &postsCreate)

	posts, err := s.postRepo.Create(uctx, slugOrId, &postsCreate)
	if err != nil {
		body, _ := json.Marshal(GetErrorMessage(err))
		cont.SetBody(body)
		if errors.Is(err, errs.ErrorThreadDoesNotExist) {
			cont.SetStatusCode(fasthttp.StatusNotFound)
		} else if errors.Is(err, errs.ErrorParentPostDoesNotExist) {
			cont.SetStatusCode(fasthttp.StatusConflict)
		} else if errors.Is(err, errs.ErrorAuthorDoesNotExist) {
			cont.SetStatusCode(fasthttp.StatusNotFound)
		}
		return
	}

	body, _ := json.Marshal(posts)
	cont.SetBody(body)
	cont.SetStatusCode(fasthttp.StatusCreated)
}


func (s *PostService) Get(cont *fasthttp.RequestCtx) {
	cont.SetContentType("application/json")
	uctx := cont.UserValue("cont").(context.Context)
	id, _ := strconv.Atoi(cont.UserValue("id").(string))

	related := string(cont.QueryArgs().Peek("related"))

	postGet := &models.GetPostReqForm{
		Related: strings.Split(related, ","),
	}

	posts, err := s.postRepo.Get(uctx, int64(id), postGet)

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


func (s *PostService) Update(cont *fasthttp.RequestCtx) {
	cont.SetContentType("application/json")
	uctx := cont.UserValue("cont").(context.Context)
	id, _ := strconv.Atoi(cont.UserValue("id").(string))
	var postUpdate models.PostUpdate

	_ = json.Unmarshal(cont.PostBody(), &postUpdate)

	post, err := s.postRepo.Update(uctx, int64(id), &postUpdate)
	if err != nil {
		body, _ := json.Marshal(GetErrorMessage(err))
		cont.SetBody(body)
		cont.SetStatusCode(fasthttp.StatusNotFound)
		return
	} else {
		body, _ := json.Marshal(post)
		cont.SetBody(body)
		cont.SetStatusCode(fasthttp.StatusOK)
		return
	}
}
