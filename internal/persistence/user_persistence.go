package persistence

import (
	"api/internal/models"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

type ListAllUsersResult struct {
	models.User
	TotalCount int64 `db:"total_count"`
}

type ListAllUsersArgs struct {
	Limit  int32
	Offset int32
}

const listAllUsersQuery = `
	select *, count(*) over() as total_count from "user"
	limit $1
	offset $2;
`

func (p Persistence) ListAllUsers(ctx context.Context, args *ListAllUsersArgs) ([]*ListAllUsersResult, error) {
	result := []*ListAllUsersResult{}
	if err := p.db.Select(&result, listAllUsersQuery, args.Limit, args.Offset); err != nil {
		p.logger.Error("failed to list users", "error", err.Error())
		return nil, errors.New("failed to list users")
	}

	return result, nil
}

const findUserByEmailQuery string = `
	select * from "user"
	where email = $1 and deleted_at is null
	limit 1
`

func (p Persistence) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := p.db.QueryRowxContext(ctx, findUserByEmailQuery, email).StructScan(&user); err != nil {
		p.logger.Error("failed to find user by email", "error", err.Error())
		return nil, errors.New("failed to find user")
	}

	return &user, nil
}

const createUserQuery string = `
	insert into "user" (id, email, role, password, created_at, updated_at)
	values ($1, $2, $3, $4, $5, $5);
`

func (p Persistence) CreateUser(ctx context.Context, user *models.User) error {
	now := time.Now()
	_, err := p.db.ExecContext(ctx, createUserQuery, user.Id, user.Email, user.Role, user.Password, now)
	if err != nil {
		e := err.(*pq.Error)
		p.logger.Error("failed to create user", "error", err.Error())

		switch e.Constraint {
		case models.ConstraintsUserEmailUnique:
			return errors.New("user with the provided email address already exists")

		default:
			return errors.New("failed to create user")
		}
	}

	return nil
}

type SetUserDeleteStatusArgs struct {
	UserId    string
	DeletedAt sql.NullTime
}

const setUserDeleteStatusQuery = `
	update "user"
	set deleted_at = $2
	where id = $1;
`

func (p Persistence) SetUserDeleteStatus(ctx context.Context, args *SetUserDeleteStatusArgs) error {
	if _, err := p.db.ExecContext(ctx, setUserDeleteStatusQuery, args.UserId, args.DeletedAt); err != nil {
		p.logger.Error("failed to set user delete status", "error", err.Error())
		return errors.New("failed to set user delete status")
	}

	return nil
}
