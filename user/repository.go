package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type User struct {
	ID      int       `json:"user_id"`
	Name    int       `json:"user_name"`
	Coin    int       `json:"coin"`
	Jewel   int       `json:"jewel"`
	Created time.Time `json:"created"`
}

type Repository struct {
	db  *sql.DB
	rdb *redis.Client
}

func NewRepository(db *sql.DB, rdb *redis.Client) *Repository {
	return &Repository{db: db, rdb: rdb}
}

func (r *Repository) List(ctx context.Context) ([]User, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT user_id, user_name, coin, jewel, created
		FROM users
		ORDER BY user_id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Coin, &u.Jewel, &u.Created); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *Repository) Get(ctx context.Context, id int) (*User, error) {
	cacheKey := fmt.Sprintf("user:%d", id)

	if r.rdb != nil {
		if b, err := r.rdb.Get(ctx, cacheKey).Bytes(); err == nil {
			var u User
			if err := json.Unmarshal(b, &u); err == nil {
				return &u, nil
			}
		}
	}

	var u User
	err := r.db.QueryRowContext(ctx, `
		SELECT user_id, user_name, coin, jewel, created
		FROM users
		WHERE user_id = ?
	`, id).Scan(&u.ID, &u.Name, &u.Coin, &u.Jewel, &u.Created)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if r.rdb != nil {
		if b, err := json.Marshal(&u); err == nil {
			_ = r.rdb.Set(ctx, cacheKey, b, time.Minute).Err()
		}
	}

	return &u, nil
}

func (r *Repository) Create(ctx context.Context, name, coin, jewel int) (int64, error) {
	res, err := r.db.ExecContext(ctx, `
		INSERT INTO users (user_name, coin, jewel)
		VALUES (?, ?, ?)
	`, name, coin, jewel)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *Repository) Update(ctx context.Context, id, name, coin, jewel int) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE users
		SET user_name = ?, coin = ?, jewel = ?
		WHERE user_id = ?
	`, name, coin, jewel, id)
	return err
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM users
		WHERE user_id = ?
	`, id)
	return err
}
