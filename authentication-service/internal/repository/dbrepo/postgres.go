package dbrepo

import (
	"authentication-service/internal/models"
	"context"
	"database/sql"
	"errors"
	"time"
)

func (m *postgresDBRepo) Authenticate(email, username, password string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var query string
	var queryParam string
	if email != "" {
		queryParam = email
		query = `select username, email, password, created_at, updated_at from users where email = $1 and password = $2`
	} else {
		queryParam = username
		query = `select username, email, password, created_at, updated_at from users where username = $1 and password = $2`
	}

	var newUserFromDB models.User

	err := m.DB.QueryRowContext(ctx, query, queryParam, password).Scan(&newUserFromDB)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return 0, errors.New("there is now user with this credentials")
	case err != nil:
		return 0, err
	default:
		return newUserFromDB.ID, nil
	}

}

func (m *postgresDBRepo) AddUser(user models.User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	queue := `select email, username, password from users where email = $1`
	rows, err := m.DB.QueryContext(ctx, queue, user.Email)
	defer rows.Close()
	if rows.Next() {
		return 0, errors.New("this email has already been used")
	}

	queue = `select email, username, password from users where username = $1`
	rows, err = m.DB.QueryContext(ctx, queue, user.Username)
	defer rows.Close()
	if rows.Next() {
		return 0, errors.New("this username has already been used")
	}

	var newID int
	stmt := `insert into users (username, email, password,
                   created_at, updated_at)
                   values ($1, $2, $3, $4, $5) returning id`

	err = m.DB.QueryRowContext(ctx, stmt,
		user.Username,
		user.Email,
		user.Password,
		time.Now(),
		time.Now(),
	).Scan(&newID)
	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (m *postgresDBRepo) UpdateUser(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
			update users set username = $1, email = $2, password = $3, updated_at = $4 where id = $5`

	_, err := m.DB.ExecContext(ctx, query,
		user.Username,
		user.Email,
		user.Password,
		time.Now(),
		user.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) DeleteUser(user models.User) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `delete from users where id = $1`

	_, err := m.DB.ExecContext(ctx, query, user.ID)
	if err != nil {
		return err
	}
	return nil
}

//TODO implement encrypting password to DB
//TODO implement decrypting password from DB
