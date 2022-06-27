package delivery

import (
	"context"
	"encoding/json"
	"log"
	"tp-db-project/repository"

	"github.com/valyala/fasthttp"
)

type StatusService struct {
	statusRepo repository.StatusRepo
}

func MakeServiceService(statusRepo repository.StatusRepo) StatusService {
	return StatusService{statusRepo: statusRepo}
}

func (s *StatusService) GetStatus(cont *fasthttp.RequestCtx) {
	log.Println("Получаю статус...")
	cont.SetContentType("application/json")
	uctx := cont.UserValue("cont").(context.Context)
	result, _ := s.statusRepo.GetStatus(uctx)

	body, _ := json.Marshal(result)
	cont.SetBody(body)
	cont.SetStatusCode(fasthttp.StatusOK)
}

func (s *StatusService) Clear(cont *fasthttp.RequestCtx) {
	cont.SetContentType("application/json")
	uctx := cont.UserValue("cont").(context.Context)

	_ = s.statusRepo.Clear(uctx)

	cont.SetStatusCode(fasthttp.StatusOK)
}
