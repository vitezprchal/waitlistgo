package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/vitezprchal/waitlistgo/internal/models"
)

type Repository interface {
	AddSubscriber(context context.Context, sub *models.Subscriber) error
	Close()
}

type repository struct {
	db *pgx.Conn
}

func NewSubscriberRepository(db *pgx.Conn) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) AddSubscriber(context context.Context, sub *models.Subscriber) error {
	_, err := r.db.Exec(context, "INSERT INTO subscribers (name, email) VALUES ($1, $2)", sub.Name, sub.Email)
	return err
}

func (r *repository) Close() {
	r.db.Close(context.Background())
}
