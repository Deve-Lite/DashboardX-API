CREATE TYPE "public"."theme" AS ENUM('inherit', 'light', 'dark');

CREATE TABLE "users" (
    "id" uuid NOT NULL DEFAULT gen_random_uuid(),
    "name" text NOT NULL,
    "email" text NOT NULL,
    "password" text NOT NULL,
    "is_admin" boolean NOT NULL DEFAULT false,
    "language" text NOT NULL DEFAULT 'pl',
    "theme" "public"."theme" NOT NULL DEFAULT 'inherit',
    CONSTRAINT "users_email_key" UNIQUE ("email"),
    CONSTRAINT "users_id_pkey" PRIMARY KEY ("id")
);