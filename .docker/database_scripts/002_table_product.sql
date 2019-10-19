BEGIN;

-- CREATE SEQUENCE "affiliate_order_invoice_id_seq" ------------
CREATE SEQUENCE IF NOT EXISTS "public"."product_id_seq"
INCREMENT 1
MINVALUE 1
MAXVALUE 9223372036854775807
START 1
CACHE 1;
-- -------------------------------------------------------------

COMMIT;

BEGIN;

-- CREATE SEQUENCE "affiliate_order_invoice_aff_inv" -----------
CREATE SEQUENCE IF NOT EXISTS "public"."product_num"
INCREMENT 1
MINVALUE 1
MAXVALUE 9223372036854775807
START 1
CACHE 1;
-- -------------------------------------------------------------

COMMIT;


BEGIN;

-- CREATE TABLE "affiliate_order_invoice" ----------------------
CREATE TABLE IF NOT EXISTS "public"."product" (
	"id" Bigint DEFAULT nextval('product_id_seq'::regclass) NOT NULL,
    "name" Character Varying COLLATE "pg_catalog"."default",
	"product_number" Character Varying( 2044 ) COLLATE "pg_catalog"."default" NOT NULL,
	"create_time" Timestamp Without Time Zone NOT NULL,
	"update_time" Timestamp Without Time Zone,
	"create_by" Bigint,
	"update_by" Bigint,
	PRIMARY KEY ( "id" ),
	CONSTRAINT "unique_product_number" UNIQUE( "product_number", "id" ) );
 ;
-- -------------------------------------------------------------

COMMIT;

