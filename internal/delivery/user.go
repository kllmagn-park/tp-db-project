package delivery

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"tp-db-project/internal/models"
	"tp-db-project/pkg/errs"
	"tp-db-project/repository"

	"github.com/valyala/fasthttp"
)

type UserService struct {
	userRepo repository.UserRepo
}

func MakeUserService(userRepo repository.UserRepo) UserService {
	return UserService{userRepo: userRepo}
}


func (s *UserService) Create(cont *fasthttp.RequestCtx) {
	cont.SetContentType("application/json")
	uctx := cont.UserValue("cont").(context.Context)

	nickname := cont.UserValue("nickname").(string)
	var user models.User
	_ = json.Unmarshal(cont.PostBody(), &user)
	user.Nickname = nickname

	log.Println("Создаю пользователя... [nickname = ", nickname, " , uctx = ", uctx, "]")

	userAfterCreate, err := s.userRepo.Create(uctx, &user)
	if err != nil {
		log.Println("Пользователь уже существует.")
		body, _ := json.Marshal(userAfterCreate)
		cont.SetBody(body)
		cont.SetStatusCode(fasthttp.StatusConflict)
		return
	}

	body, _ := json.Marshal((*userAfterCreate)[0])
	cont.SetBody(body)
	cont.SetStatusCode(fasthttp.StatusCreated)
}


func (s *UserService) Get(cont *fasthttp.RequestCtx) {
	cont.SetContentType("application/json")
	uctx := cont.UserValue("cont").(context.Context)

	nickname := cont.UserValue("nickname").(string)
	user, err := s.userRepo.Get(uctx, nickname)
	if err != nil {
		body, _ := json.Marshal(GetErrorMessage(err))
		cont.SetBody(body)
		cont.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	body, _ := json.Marshal(user)
	cont.SetBody(body)
	cont.SetStatusCode(fasthttp.StatusOK)
}


func (s *UserService) Update(cont *fasthttp.RequestCtx) {
	log.Println("Обновление данных пользователя...")
	cont.SetContentType("application/json")
	uctx := cont.UserValue("cont").(context.Context)

	nickname := cont.UserValue("nickname").(string)
	var updateData models.UserUpdate
	_ = json.Unmarshal(cont.PostBody(), &updateData)

	user, err := s.userRepo.Update(uctx, nickname, &updateData)
	if err != nil {
		log.Println("Обработка ошибки...")
		body, _ := json.Marshal(GetErrorMessage(err))
		log.Println(err, GetErrorMessage(err))
		cont.SetBody(body)
		if errors.Is(err, errs.ErrorUserDoesNotExist) {
			cont.SetStatusCode(fasthttp.StatusNotFound)
		} else if errors.Is(err, errs.ErrorConflictUpdateUser) {
			cont.SetStatusCode(fasthttp.StatusConflict)
		}
		return
	}

	body, _ := json.Marshal(user)
	cont.SetBody(body)
	cont.SetStatusCode(fasthttp.StatusOK)
}
