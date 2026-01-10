package services

import (
	"github.com/Mangrover007/banking-backend/app/internals/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type AccountService interface {
	GetUserAccounts(ctx *gin.Context, user repository.User) ([]repository.FindAllAccountsByIDRow, error)
}

type accountService struct {
	conn  *pgx.Conn
	query *repository.Queries
}

func NewAccountService(conn *pgx.Conn, query *repository.Queries) AccountService {
	return &accountService{
		conn:  conn,
		query: query,
	}
}

func (s *accountService) GetUserAccounts(ctx *gin.Context, user repository.User) ([]repository.FindAllAccountsByIDRow, error) {
	res, err := s.query.FindAllAccountsByID(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	return res, nil
}
