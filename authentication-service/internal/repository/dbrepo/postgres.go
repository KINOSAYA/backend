package dbrepo

import (
	"authentication-service/internal/models"
	"context"
	"time"
)

func (m *postgresDBRepo) GetUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select username, email, password, created_at, updated_at where id = $1`

	var user models.User
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&user)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (m *postgresDBRepo) AddUser(user models.User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int
	stmt := `insert into users (username, email, password,
                   created_at, updated_at)
                   values ($1, $2, $3, $4, $5) returning id`

	err := m.DB.QueryRowContext(ctx, stmt,
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
