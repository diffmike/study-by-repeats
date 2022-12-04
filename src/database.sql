CREATE TABLE "public"."users"
(
    "id"         numeric   NOT NULL,
    "username"   varchar   NOT NULL,
    "uuid"       numeric   NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT now(),
    PRIMARY KEY ("id")
);