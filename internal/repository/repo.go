package repository

import (
	"context"
	"database/sql"
	"usdt/internal/models"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Repo struct {
	DB     *sql.DB
	Tracer trace.Tracer
	Logger *zap.Logger
}

// SaveRate inserts a new rate into the database.
func (r *Repo) SaveRate(ctx context.Context, rate *models.Rate) error {
	ctx, span := r.Tracer.Start(ctx, "Repo.SaveRate")
	defer span.End()

	_, err := r.DB.ExecContext(ctx, `INSERT INTO rates (ask, bid, timestamp) VALUES ($1, $2, $3)`,
		rate.Ask, rate.Bid, rate.Timestamp)
	if err != nil {
		span.RecordError(err)
		if r.Logger != nil {
			r.Logger.Error("failed to save rate",
				zap.Float64("ask", rate.Ask),
				zap.Float64("bid", rate.Bid),
				zap.Time("timestamp", rate.Timestamp),
				zap.Error(err))
		}
		return err
	}

	if r.Logger != nil {
		r.Logger.Info("rate saved successfully",
			zap.Float64("ask", rate.Ask),
			zap.Float64("bid", rate.Bid),
			zap.Time("timestamp", rate.Timestamp))
	}
	return nil
}
