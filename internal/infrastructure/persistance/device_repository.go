package persistance

import (
	"context"
	"fmt"
	"strings"

	"github.com/Deve-Lite/DashboardX-API-PoC/internal/domain"
	"github.com/Deve-Lite/DashboardX-API-PoC/internal/domain/repository"
	ae "github.com/Deve-Lite/DashboardX-API-PoC/pkg/errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type deviceRepository struct {
	db *sqlx.DB
}

func NewDeviceRepository(db *sqlx.DB) repository.DeviceRepository {
	return &deviceRepository{db: db}
}

func (r *deviceRepository) Get(ctx context.Context, deviceID uuid.UUID, userID uuid.UUID) (*domain.Device, error) {
	device := &domain.Device{}

	sql := `
		SELECT "id", "broker_id", "name", "icon_name", "icon_background_color",
			"placing", "base_path", "created_at", "updated_at"
		FROM "devices" WHERE "id" = $1 AND "user_id" = $2
	`

	if err := r.db.GetContext(ctx, device, sql, deviceID, userID); err != nil {
		return nil, errors.Wrap(err, "deviceRepository.Get.GetContext")
	}

	return device, nil
}

func (r *deviceRepository) List(ctx context.Context, filters *domain.ListDeviceFilters) ([]*domain.Device, error) {
	var devices []*domain.Device

	sql := `
		SELECT "id", "broker_id", "name", "icon_name", "icon_background_color",
			"placing", "base_path", "created_at", "updated_at"
		FROM "devices" WHERE "user_id" = $1
	`

	if filters.BrokerID.Valid {
		sql = fmt.Sprintf(`%s AND "broker_id" = '%s'`, sql, filters.BrokerID.UUID.String())
	}

	if err := r.db.SelectContext(ctx, &devices, sql, filters.UserID); err != nil {
		return nil, errors.Wrap(err, "deviceRepository.List.SelectContext")
	}

	return devices, nil
}

func (r *deviceRepository) Create(ctx context.Context, device *domain.CreateDevice) (uuid.UUID, error) {
	var f strings.Builder
	f.WriteString(`"user_id", "name", "icon_name", "icon_background_color", "broker_id"`)
	p := fmt.Sprintf("'%s', '%s', '%s', '%s'",
		device.UserID,
		device.Name,
		device.IconName,
		device.IconBackgroundColor)

	if device.BrokerID.Valid {
		p = fmt.Sprintf("%s, '%s'", p, device.BrokerID.UUID)
	} else {
		p = fmt.Sprintf("%s, NULL", p)
	}

	if device.BasePath.Set {
		f.WriteString(`, "base_path"`)
		if device.BasePath.Null {
			p = fmt.Sprintf("%s, NULL", p)
		} else {
			p = fmt.Sprintf("%s, '%s'", p, device.BasePath.String)
		}
	}

	if device.Placing.Set {
		f.WriteString(`, "placing"`)
		if device.Placing.Null {
			p = fmt.Sprintf("%s, NULL", p)
		} else {
			p = fmt.Sprintf("%s, '%s'", p, device.Placing.String)
		}
	}

	created := &domain.Device{}

	sql := fmt.Sprintf(`INSERT INTO "devices" (%s) VALUES (%s) RETURNING id`, f.String(), p)

	if err := r.db.GetContext(ctx, created, sql); err != nil {
		return uuid.Nil, errors.Wrap(err, "deviceRepository.Create.GetContext")
	}

	return created.ID, nil
}

func (r *deviceRepository) Update(ctx context.Context, device *domain.UpdateDevice) error {
	p := []string{}

	if device.BrokerID.Set {
		if device.BrokerID.Null {
			p = append(p, `"broker_id" = NULL`)
		} else {
			p = append(p, fmt.Sprintf(`"broker_id" = '%s'`, device.BrokerID.Value))
		}
	}

	if device.Name.Set && !device.Name.Null {
		p = append(p, fmt.Sprintf(`"name" = '%s'`, device.Name.String))
	}

	if device.IconName.Set && !device.IconName.Null {
		p = append(p, fmt.Sprintf(`"icon_name" = '%s'`, device.IconName.String))
	}

	if device.IconBackgroundColor.Set && !device.IconBackgroundColor.Null {
		p = append(p, fmt.Sprintf(`"icon_background_color" = '%s'`, device.IconBackgroundColor.String))
	}

	if device.BasePath.Set {
		if device.BasePath.Null {
			p = append(p, `"base_path" = NULL`)
		} else {
			p = append(p, fmt.Sprintf(`"base_path" = '%s'`, device.BasePath.String))
		}
	}

	if device.Placing.Set {
		if device.Placing.Null {
			p = append(p, `"placing" = NULL`)
		} else {
			p = append(p, fmt.Sprintf(`"placing" = '%s'`, device.Placing.String))
		}
	}

	if len(p) == 0 {
		return ae.ErrMissingParams
	}

	p = append(p, `"updated_at" = now()`)

	sql := fmt.Sprintf(`UPDATE "devices" SET %s WHERE "id" = $1 AND "user_id" = $2`, strings.Join(p, ","))

	sr, err := r.db.ExecContext(ctx, sql, device.ID, device.UserID)
	if err != nil {
		return errors.Wrap(err, "deviceRepository.Update.ExecContext")
	}

	if af, _ := sr.RowsAffected(); af == 0 {
		return ae.ErrDeviceNotFound
	}
	return nil
}

func (r *deviceRepository) Delete(ctx context.Context, deviceID uuid.UUID, userID uuid.UUID) error {
	sql := `
		DELETE FROM "devices"
		WHERE "id" = $1 AND "user_id" = $2
	`

	sr, err := r.db.ExecContext(ctx, sql, deviceID, userID)
	if err != nil {
		return errors.Wrap(err, "deviceRepository.Delete.ExecContext")
	}

	if af, _ := sr.RowsAffected(); af == 0 {
		return ae.ErrDeviceNotFound
	}
	return nil
}
