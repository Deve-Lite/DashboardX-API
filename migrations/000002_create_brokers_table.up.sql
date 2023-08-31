CREATE TABLE "brokers" (
    "id" uuid NOT NULL DEFAULT gen_random_uuid(),
    "user_id" uuid NOT NULL,
    "name" text NOT NULL, 
    "server" text NOT NULL,
    "port" integer NOT NULL,
    "keep_alive" integer NOT NULL, 
    "icon_name" text NOT NULL,
    "icon_background_color" text NOT NULL,
    "is_ssl" boolean NOT NULL DEFAULT false,
    "username" text,
    "password" text, 
    "client_id" text, 
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    CONSTRAINT "brokers_id_pkey" PRIMARY KEY ("id"),
    CONSTRAINT "brokers_user_id_fkey" FOREIGN KEY ("user_id")
        REFERENCES "users"("id")
        ON DELETE CASCADE
        ON UPDATE NO ACTION
);

CREATE INDEX "brokers_user_id_idx" ON "brokers"("user_id");