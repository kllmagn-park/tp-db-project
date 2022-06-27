package delivery

import (
	"context"
	"encoding/json"
	"tp-db-project/internal/models"
	"tp-db-project/repository"

	"github.com/valyala/fasthttp"
)

type VoteService struct {
	voteRepo repository.VoteRepo
}

func MakeVoteService(voteRepo repository.VoteRepo) VoteService {
	return VoteService{voteRepo: voteRepo}
}


func (s *VoteService) Create(cont *fasthttp.RequestCtx) {
	cont.SetContentType("application/json")
	uctx := cont.UserValue("cont").(context.Context)

	slugOrId := cont.UserValue("slug_or_id").(string)
	var voteCreate models.CreateVoteForm
	_ = json.Unmarshal(cont.PostBody(), &voteCreate)

	thread, err := s.voteRepo.Create(uctx, slugOrId, &voteCreate)

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
