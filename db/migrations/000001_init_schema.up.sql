CREATE TABLE "users" (
    "id" UUID NOT NULL DEFAULT (gen_random_uuid()),
    "name" VARCHAR NOT NULL,
    "phone_number" VARCHAR NOT NULL UNIQUE,
    "opt" VARCHAR,
    "opt_expiration_time" TIMESTAMP(3),

    CONSTRAINT "users_pkey" PRIMARY KEY ("id")
);