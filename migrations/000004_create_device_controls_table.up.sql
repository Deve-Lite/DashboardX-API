CREATE TYPE "public"."control_type" AS ENUM(
    'button',
    'color',
    'date-time',
    'radio',
    'slider',
    'state',
    'switch',
    'text-out'
);

CREATE TYPE "public"."qos_level" AS ENUM('0', '1', '2');

CREATE TABLE "device_controls" (
    "id" uuid NOT NULL DEFAULT gen_random_uuid(),
    "device_id" uuid NOT NULL,
    "name" text NOT NULL,
    "type" "public"."control_type" NOT NULL,
    "icon_name" text NOT NULL,
    "icon_background_color" text NOT NULL,
    "is_available" boolean NOT NULL,
    "is_confirmation_required" boolean NOT NULL,
    "can_notify_on_publish" boolean NOT NULL,
    "can_display_name" boolean NOT NULL,
    "quality_of_service" "public"."qos_level" NOT NULL,
    "topic" text NOT NULL,
    "attributes" jsonb NOT NULL,
    CONSTRAINT "device_controls_id_pkey" PRIMARY KEY ("id"),
    CONSTRAINT "device_controls_device_id_fkey" FOREIGN KEY ("device_id")
        REFERENCES "devices"("id")
        ON DELETE CASCADE
        ON UPDATE NO ACTION
);

CREATE INDEX "device_controls_device_id_idx" ON "device_controls"("device_id");