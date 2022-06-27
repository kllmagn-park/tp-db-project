package postgres

import (
	"context"
	"tp-db-project/internal/models"
	"tp-db-project/pkg/errs"
	"tp-db-project/pkg/queries"
	"tp-db-project/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	conn *pgxpool.Pool
}

func InitUserRepository(db *pgxpool.Pool) repository.UserRepo {
	return &UserRepository{conn: db}
}

func (s *UserRepository) getUserByNicknameOrEmail(cont context.Context, nickname string, email string) (*[]models.User, error) {
	rows, err := s.conn.Query(cont, queries.GetUserByNicknameOrEmailCommand, nickname, email)
	if err != nil {
		return nil, err
	}

	users := make([]models.User, 0, rows.CommandTag().RowsAffected())
	for rows.Next() {
		user := models.User{}
		err = rows.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return &users, nil
}

func (s *UserRepository) Create(cont context.Context, user *models.User) (*[]models.User, error) {
	_, err := s.conn.Exec(cont, queries.CreateUserCommand, user.Nickname, user.Fullname, user.About, user.Email)
	if err != nil {
		checkAlreadyExist, err := s.getUserByNicknameOrEmail(cont, user.Nickname, user.Email)
		if err == nil && len(*checkAlreadyExist) > 0 {
			return checkAlreadyExist, errs.ErrorUserAlreadyExist
		} else {
			return nil, errs.ErrorUserAlreadyExist
		}
	}
	userToReturn := make([]models.User, 0, 1)
	userToReturn = append(userToReturn, *user)
	return &userToReturn, nil
}

func (s *UserRepository) Update(cont context.Context, nickname string, updateData *models.UserUpdate) (*models.User, error) {
	user, err := s.Get(cont, nickname)
	if err != nil {
		return nil, errs.ErrorUserDoesNotExist
	}

	if updateData.Fullname == "" {
		updateData.Fullname = user.Fullname
	} else {
		user.Fullname = updateData.Fullname
	}
	if updateData.About == "" {
		updateData.About = user.About
	} else {
		user.About = updateData.About
	}
	if updateData.Email == "" {
		updateData.Email = user.Email
	} else {
		user.Email = updateData.Email
	}

	_, err = s.conn.Exec(cont, queries.UpdateUserCommand, updateData.Fullname, updateData.About, updateData.Email, nickname)

	if err != nil {
		return nil, errs.ErrorConflictUpdateUser
	}

	return user, nil
}

func (s *UserRepository) Get(cont context.Context, nicknameOrEmail string) (*models.User, error) {
	var user models.User
	err := s.conn.QueryRow(cont, queries.GetUserByNicknameCommand, nicknameOrEmail).Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
	if err != nil {
		err = s.conn.QueryRow(cont, queries.GetUserByEmailCommand, nicknameOrEmail).Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
	}

	if err != nil {
		return nil, errs.ErrorUserDoesNotExist
	}

	return &user, nil
}
