package persistance

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Deve-Lite/DashboardX-API-PoC/internal/domain"
	"github.com/Deve-Lite/DashboardX-API-PoC/internal/domain/repository"
	ae "github.com/Deve-Lite/DashboardX-API-PoC/pkg/errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type brokerRepository struct {
	db *sqlx.DB
}

func NewBrokerRepository(db *sqlx.DB) repository.BrokerRepository {
	return &brokerRepository{db}
}

func (r *brokerRepository) Get(ctx context.Context, brokerID uuid.UUID, userID uuid.UUID) (*domain.Broker, error) {
	broker := &domain.Broker{}

	sqls := `
		SELECT "id", "user_id", "name", "server", "port", "keep_alive", "icon_name", "icon_background_color", "is_ssl", "username",
			"password", "client_id", "created_at", "updated_at"
		FROM "brokers"
		WHERE "id" = $1 AND "user_id" = $2
	`

	if err := r.db.GetContext(ctx, broker, sqls, brokerID, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ae.ErrBrokerNotFound
		}

		return nil, errors.Wrap(err, "brokerRepository.Get.GetContext")
	}

	return broker, nil
}

func (r *brokerRepository) List(ctx context.Context, userID uuid.UUID) ([]*domain.Broker, error) {
	var brokers []*domain.Broker

	sql := `
		SELECT "id", "user_id", "name", "server", "port", "keep_alive", "icon_name", "icon_background_color", "is_ssl", "username",
			"password", "client_id", "created_at", "updated_at"
		FROM "brokers"
		WHERE "user_id" = $1
	`

	if err := r.db.SelectContext(ctx, &brokers, sql, userID); err != nil {
		return nil, errors.Wrap(err, "brokerRepository.List.SelectContext")
	}

	return brokers, nil
}

func (r *brokerRepository) Create(ctx context.Context, broker *domain.CreateBroker) (uuid.UUID, error) {
	var f strings.Builder
	f.WriteString(`"user_id", "name", "server", "port", "keep_alive", "icon_name", "icon_background_color", "is_ssl"`)
	p := fmt.Sprintf("'%s', '%s', '%s', %d, %d, '%s', '%s', %v",
		broker.UserID.String(),
		broker.Name,
		broker.Server,
		broker.Port,
		broker.KeepAlive,
		broker.IconName,
		broker.IconBackgroundColor,
		broker.IsSSL)

	if broker.ClientID.Set {
		f.WriteString(`, client_id`)
		if broker.ClientID.Null {
			p = fmt.Sprintf("%s, NULL", p)
		} else {
			p = fmt.Sprintf("%s, '%s'", p, broker.ClientID.String)
		}
	}

	if broker.Username.Set {
		f.WriteString(`, "username"`)
		if broker.Username.Null {
			p = fmt.Sprintf("%s, NULL", p)
		} else {
			p = fmt.Sprintf("%s, '%s'", p, broker.Username.String)
		}
	}

	if broker.Password.Set {
		f.WriteString(`, "password"`)
		if broker.Password.Null {
			p = fmt.Sprintf("%s, NULL", p)
		} else {
			p = fmt.Sprintf("%s, '%s'", p, broker.Password.String)
		}
	}

	created := &domain.Broker{}

	sql := fmt.Sprintf(`INSERT INTO "brokers" (%s) VALUES (%s) RETURNING "id"`, f.String(), p)

	if err := r.db.GetContext(ctx, created, sql); err != nil {
		return uuid.Nil, errors.Wrap(err, "brokerRepository.Create.GetContext")
	}

	return created.ID, nil
}

func (r *brokerRepository) Update(ctx context.Context, broker *domain.UpdateBroker) error {
	var p []string

	if broker.Name.Set {
		p = append(p, fmt.Sprintf(`"name" = '%s'`, broker.Name.String))
	}

	if broker.Server.Set {
		p = append(p, fmt.Sprintf(`"server" = '%s'`, broker.Server.String))
	}

	if broker.IconName.Set {
		p = append(p, fmt.Sprintf(`"icon_name" = '%s'`, broker.IconName.String))
	}

	if broker.IconBackgroundColor.Set {
		p = append(p, fmt.Sprintf(`"icon_background_color" = '%s'`, broker.IconBackgroundColor.String))
	}

	if broker.KeepAlive.Set {
		p = append(p, fmt.Sprintf(`"keep_alive" = %d`, broker.KeepAlive.Uint16))
	}

	if broker.Port.Set {
		p = append(p, fmt.Sprintf(`"port" = %d`, broker.Port.Uint16))
	}

	if broker.IsSSL.Set {
		p = append(p, fmt.Sprintf(`"is_ssl" = %v`, broker.IsSSL.Bool))
	}

	if broker.Username.Set {
		if broker.Username.Null {
			p = append(p, `"username" = NULL`)
		} else {
			p = append(p, fmt.Sprintf(`"username" = '%s'`, broker.Username.String))
		}
	}

	if broker.Password.Set {
		if broker.Password.Null {
			p = append(p, `"password" = NULL`)
		} else {
			p = append(p, fmt.Sprintf(`"password" = '%s'`, broker.Password.String))
		}
	}

	if broker.ClientID.Set {
		if broker.ClientID.Null {
			p = append(p, `"client_id" = NULL`)
		} else {
			p = append(p, fmt.Sprintf(`"client_id" = '%s'`, broker.ClientID.String))
		}
	}

	if len(p) == 0 {
		return ae.ErrMissingParams
	}

	p = append(p, `"updated_at" = now()`)

	sql := fmt.Sprintf(`UPDATE "brokers" SET %s WHERE "id" = '%s' AND "user_id" = '%s'`, strings.Join(p, ","), broker.ID.String(), broker.UserID.String())

	sr, err := r.db.ExecContext(ctx, sql)
	if err != nil {
		return errors.Wrap(err, "brokerRepository.Update.ExecContext")
	}

	if af, _ := sr.RowsAffected(); af == 0 {
		return ae.ErrBrokerNotFound
	}

	return nil
}

func (r *brokerRepository) Delete(ctx context.Context, brokerID uuid.UUID, userID uuid.UUID) error {
	sql := `
		DELETE FROM "brokers"
		WHERE "id" = $1 AND "user_id" = $2 
	`

	sr, err := r.db.ExecContext(ctx, sql, brokerID, userID)
	if err != nil {
		return errors.Wrap(err, "brokerRepository.Delete.ExecContext")
	}

	if af, _ := sr.RowsAffected(); af == 0 {
		return ae.ErrBrokerNotFound
	}

	return nil
}
