package postgres

import (
	"context"
	"strconv"
	"tp-db-project/internal/models"
	"tp-db-project/pkg/errs"
	"tp-db-project/pkg/queries"
	"tp-db-project/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type VoteRepository struct {
	conn *pgxpool.Pool
}

func InitVoteRepository(db *pgxpool.Pool) repository.VoteRepo {
	return &VoteRepository{conn: db}
}

func (s *VoteRepository) Create(cont context.Context, slugOrId string, vote *models.CreateVoteForm) (*models.Thread, error) {
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
	var checkVote models.Vote
	err = s.conn.QueryRow(cont, queries.GetVoteByNicknameAndThreadCommand, vote.Nickname, thread.Id).Scan(&checkVote.Nickname, &checkVote.Thread, &checkVote.Voice)
	if err != nil {
		_, err = s.conn.Exec(cont, queries.CreateVoteCommand, vote.Nickname, thread.Id, vote.Voice)
		if err != nil {
			return nil, errs.ErrorUserDoesNotExist
		}
		thread.Votes += vote.Voice
	} else {
		_, _ = s.conn.Exec(cont, queries.UpdateVoteCommand, vote.Voice, vote.Nickname, thread.Id)
		if vote.Voice != checkVote.Voice {
			thread.Votes += 2 * vote.Voice
		}
	}
	return &thread, nil
}
