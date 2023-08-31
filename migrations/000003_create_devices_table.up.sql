CREATE TABLE "devices" (
    "id" uuid NOT NULL DEFAULT gen_random_uuid(),
    "user_id" uuid NOT NULL,
    "broker_id" uuid,
    "name" text NOT NULL,
    "icon_name" text NOT NULL,
    "icon_background_color" text NOT NULL,
    "base_path" text,
    "placing" text,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    CONSTRAINT "devices_id_pkey" PRIMARY KEY ("id"),
    CONSTRAINT "devices_user_id_fkey" FOREIGN KEY ("user_id")
        REFERENCES "users"("id")
        ON DELETE CASCADE
        ON UPDATE NO ACTION,
    CONSTRAINT "devices_broker_id_fkey" FOREIGN KEY ("broker_id")
        REFERENCES "brokers"("id")
        ON DELETE SET NULL
        ON UPDATE NO ACTION
);

CREATE INDEX "devices_user_id_idx" ON "devices"("user_id");

CREATE INDEX "devices_broker_id_idx" ON "devices"("broker_id");