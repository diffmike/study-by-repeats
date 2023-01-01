CREATE SEQUENCE IF NOT EXISTS sessions_id_seq;

CREATE TABLE "public"."sessions"
(
    "id"         numeric   NOT NULL DEFAULT nextval('sessions_id_seq'::regclass),
    "tg_id"      numeric   NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT now(),
    PRIMARY KEY ("id")
);

CREATE SEQUENCE IF NOT EXISTS cards_id_seq;

CREATE TABLE "public"."cards"
(
    "id"           numeric   NOT NULL DEFAULT nextval('cards_id_seq'::regclass),
    "front"        varchar   NOT NULL,
    "back"         varchar,
    "tg_id"        numeric   NOT NULL,
    "created_at"   timestamp NOT NULL DEFAULT now(),
    "repeat_after" timestamp,
    PRIMARY KEY ("id")
);

CREATE SEQUENCE IF NOT EXISTS id_seq;

CREATE TABLE "public"."users"
(
    "id"         numeric   NOT NULL DEFAULT nextval('id_seq'::regclass),
    "username"   varchar   NOT NULL,
    "tg_id"      numeric   NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT now(),
    PRIMARY KEY ("id")
);
CREATE SEQUENCE IF NOT EXISTS repeats_id_seq;

CREATE TABLE "public"."repeats"
(
    "id"         numeric NOT NULL DEFAULT nextval('repeats_id_seq'::regclass),
    "card_id"    numeric,
    "session_id" numeric NOT NULL,
    "repeat_in"  numeric,
    CONSTRAINT "repeats_session_id_fkey" FOREIGN KEY ("session_id") REFERENCES "public"."sessions" ("id") ON DELETE SET NULL ON UPDATE CASCADE,
    CONSTRAINT "repeats_card_id_fkey" FOREIGN KEY ("card_id") REFERENCES "public"."cards" ("id") ON DELETE SET NULL ON UPDATE CASCADE,
    PRIMARY KEY ("id")
);

COMMENT ON COLUMN "public"."repeats"."repeat_in" IS 'In hours';