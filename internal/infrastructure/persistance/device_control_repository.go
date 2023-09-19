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

type deviceControlRepository struct {
	db *sqlx.DB
}

func NewDeviceControlRepository(db *sqlx.DB) repository.DeviceControlRepository {
	return &deviceControlRepository{db: db}
}

func (r *deviceControlRepository) Exist(ctx context.Context, filters *domain.ExistDeviceControlFilters) (bool, error) {
	var controls []*domain.DeviceControl

	sql := `
		SELECT "id"
		FROM "device_controls"
		WHERE "device_id" = $1 AND "type" = $2
	`

	if err := r.db.SelectContext(ctx, &controls, sql, filters.DeviceID, filters.Type); err != nil {
		return false, errors.Wrap(err, "deviceControlRepository.Exist.SelectContext")
	}

	return len(controls) > 0, nil
}

func (r *deviceControlRepository) List(ctx context.Context, filters *domain.ListDeviceControlFilters) ([]*domain.DeviceControl, error) {
	var controls []*domain.DeviceControl

	sql := `
		SELECT
			"id", "device_id", "name", "type", "quality_of_service", "icon_name", "icon_background_color",
			"is_available", "is_confirmation_required", "can_notify_on_publish", "can_display_name",
			"topic", "attributes"
		FROM "device_controls"
		WHERE "device_id" = $1
	`

	if err := r.db.SelectContext(ctx, &controls, sql, filters.DeviceID); err != nil {
		return nil, errors.Wrap(err, "deviceControlRepository.List.SelectContext")
	}

	return controls, nil
}

func (r *deviceControlRepository) Create(ctx context.Context, control *domain.CreateDeviceControl) (uuid.UUID, error) {
	created := &domain.DeviceControl{}

	attr, err := control.Attributes.Value()
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "deviceControlRepository.Create.Attributes.Value")
	}

	sql := fmt.Sprintf(`
		INSERT INTO "device_controls" (
			"device_id", "name", "type", "icon_name", "icon_background_color", "is_available",
			"is_confirmation_required", "can_notify_on_publish", "can_display_name",
			"quality_of_service", "topic", "attributes")
		VALUES (
			'%s', '%s', '%s'::"control_type", '%s', '%s', %v, %v, %v, %v, '%d'::"qos_level", '%s', $1			 
		) RETURNING "id"`,
		control.DeviceID,
		control.Name,
		control.Type,
		control.IconName,
		control.IconBackgroundColor,
		control.IsAvailable,
		control.IsConfirmationRequired,
		control.CanNotifyOnPublish,
		control.CanDisplayName,
		control.QoS,
		control.Topic,
	)

	if err := r.db.GetContext(ctx, created, sql, attr); err != nil {
		return uuid.Nil, errors.Wrap(err, "deviceControlRepository.Create.GetContext")
	}

	return created.ID, nil
}

func (r *deviceControlRepository) Update(ctx context.Context, control *domain.UpdateDeviceControl) error {
	f := []string{}

	if control.Name != nil {
		f = append(f, fmt.Sprintf(`"name" = '%s'`, *control.Name))
	}

	if control.IconName != nil {
		f = append(f, fmt.Sprintf(`"icon_name" = '%s'`, *control.IconName))
	}

	if control.IconBackgroundColor != nil {
		f = append(f, fmt.Sprintf(`"icon_background_color" = '%s'`, *control.IconBackgroundColor))
	}

	if control.Type != nil {
		f = append(f, fmt.Sprintf(`"type" = '%s'::"control_type"`, *control.Type))
	}

	if control.CanDisplayName != nil {
		f = append(f, fmt.Sprintf(`"can_display_name" = %v`, *control.CanDisplayName))
	}

	if control.CanNotifyOnPublish != nil {
		f = append(f, fmt.Sprintf(`"can_notify_on_publish" = %v`, *control.CanNotifyOnPublish))
	}

	if control.IsAvailable != nil {
		f = append(f, fmt.Sprintf(`"is_available" = %v`, *control.IsAvailable))
	}

	if control.IsConfirmationRequired != nil {
		f = append(f, fmt.Sprintf(`"is_confirmation_required" = %v`, *control.IsConfirmationRequired))
	}

	if control.Attributes != nil {
		attr, err := control.Attributes.Value()
		if err != nil {
			return errors.Wrap(err, "deviceControlRepository.Update.Attributes.Value")
		}

		f = append(f, fmt.Sprintf(`"attributes" = '%v'`, string(attr.([]byte))))
	}

	if control.QoS != nil {
		f = append(f, fmt.Sprintf(`"quality_of_service" = '%d'::"qos_level"`, *control.QoS))
	}

	if control.Topic != nil {
		f = append(f, fmt.Sprintf(`"topic" = '%s'`, *control.Topic))
	}

	sql := fmt.Sprintf(`UPDATE "device_controls" SET %s WHERE "id" = $1 AND "device_id" = $2`, strings.Join(f, ","))

	sr, err := r.db.ExecContext(ctx, sql, control.ID, control.DeviceID)
	if err != nil {
		return errors.Wrap(err, "deviceControlRepository.Update.ExecContext")
	}

	if af, _ := sr.RowsAffected(); af == 0 {
		return ae.ErrDeviceControlNotFound
	}
	return nil
}

func (r *deviceControlRepository) Delete(ctx context.Context, deviceID uuid.UUID, controlID uuid.UUID) error {

	sql := `DELETE FROM "device_controls" WHERE "id" = $1 AND "device_id" = $2`

	sr, err := r.db.ExecContext(ctx, sql, controlID, deviceID)
	if err != nil {
		return errors.Wrap(err, "deviceControlRepository.Delete.ExecContext")
	}

	if af, _ := sr.RowsAffected(); af == 0 {
		return ae.ErrDeviceControlNotFound
	}
	return nil
}
