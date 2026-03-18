package order

import (
	"context"
	"database/sql"
	"time"
)

type Order struct {
	ID      int       `json:"order_id"`
	UserID  int       `json:"user_id"`
	Amount  int       `json:"amount"`
	Created time.Time `json:"created"`
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List(ctx context.Context) ([]Order, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT order_id, user_id, amount, created
		FROM orders
		ORDER BY order_id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.Amount, &o.Created); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, rows.Err()
}

func (r *Repository) Get(ctx context.Context, id int) (*Order, error) {
	var o Order
	err := r.db.QueryRowContext(ctx, `
		SELECT order_id, user_id, amount, created
		FROM orders
		WHERE order_id = ?
	`, id).Scan(&o.ID, &o.UserID, &o.Amount, &o.Created)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *Repository) Create(ctx context.Context, userID, amount int) (int64, error) {
	res, err := r.db.ExecContext(ctx, `
		INSERT INTO orders (user_id, amount)
		VALUES (?, ?)
	`, userID, amount)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *Repository) Update(ctx context.Context, id, amount int) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE orders
		SET amount = ?
		WHERE order_id = ?
	`, amount, id)
	return err
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM orders
		WHERE order_id = ?
	`, id)
	return err
}
