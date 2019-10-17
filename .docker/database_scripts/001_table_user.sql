BEGIN;

-- CREATE SEQUENCE "affiliate_id_seq" --------------------------
CREATE SEQUENCE IF NOT EXISTS "public"."user_id_seq"
INCREMENT 1
MINVALUE 1
MAXVALUE 9223372036854775807
START 1
CACHE 1;
-- -------------------------------------------------------------

COMMIT;

BEGIN;

-- CREATE TABLE "affiliate" ------------------------------------
CREATE TABLE IF NOT EXISTS "public"."user" (
	"id" Bigint DEFAULT nextval('user_id_seq'::regclass) NOT NULL,
	"name" Character Varying( 2044 ) COLLATE "pg_catalog"."default" NOT NULL,
	"status" Smallint DEFAULT 1,
	"create_by" Bigint,
	"create_time" Timestamp Without Time Zone NOT NULL,
	"update_by" Bigint,
	"update_time" Timestamp Without Time Zone,
	PRIMARY KEY ( "id" ) );
;
-- -------------------------------------------------------------

COMMIT;
