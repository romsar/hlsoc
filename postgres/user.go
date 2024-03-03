package postgres

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/romsar/hlsoc"
	"strings"
)

func (db *DB) CreateUser(ctx context.Context, user *hlsoc.User) error {
	query := `
		INSERT INTO users (id, password, first_name, second_name, birth_date, gender, biography, city) 
		VALUES (@id, @password, @firstName, @secondName, @birthDate, @gender, @biography, @city)
		RETURNING id
	`
	id := uuid.New()
	args := pgx.NamedArgs{
		"id":         id,
		"password":   user.Password,
		"firstName":  user.FirstName,
		"secondName": user.SecondName,
		"birthDate":  user.BirthDate.Format("2006-01-02"),
		"gender":     user.Gender,
		"biography":  user.Biography,
		"city":       user.City,
	}
	_, err := db.db.ExecContext(ctx, query, args)
	if err != nil {
		return fmt.Errorf("unable to insert user row: %w", err)
	}

	user.ID = id

	return nil
}

func (db *DB) GetUser(ctx context.Context, filter hlsoc.UserFilter) (*hlsoc.User, error) {
	filter.Limit = 1

	users, err := db.GetUsers(ctx, filter)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, hlsoc.ErrUserNotFound
	}

	return users[0], nil
}

func (db *DB) GetUsers(ctx context.Context, filter hlsoc.UserFilter) ([]*hlsoc.User, error) {
	where, args := []string{"1 = 1"}, pgx.NamedArgs{}
	if filter.ID != uuid.Nil {
		where = append(where, "id = @id")
		args["id"] = filter.ID.String()
	}

	query := `
		SELECT id, password, first_name, second_name, birth_date, gender, biography, city 
		FROM users
		WHERE ` + strings.Join(where, " AND ") + `
		` + FormatLimitOffset(filter.Limit, filter.Offset) + `
	`

	rows, err := db.db.QueryContext(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("unable to query users: %w", err)
	}
	defer rows.Close()

	var users []*hlsoc.User
	for rows.Next() {
		user := hlsoc.User{}

		err = rows.Scan(
			&user.ID,
			&user.Password,
			&user.FirstName,
			&user.SecondName,
			&user.BirthDate,
			&user.Gender,
			&user.Biography,
			&user.City,
		)
		if err != nil {
			return nil, fmt.Errorf("unable to scan row: %w", err)
		}
		users = append(users, &user)
	}

	return users, nil
}
