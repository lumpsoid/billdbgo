CREATE TABLE "invoice" (
	"invoice_id" TEXT NOT NULL UNIQUE,
	"invoice_name" TEXT NOT NULL,
	"invoice_date" TEXT NOT NULL,
	"invoice_price" REAL NOT NULL,
	"invoice_currency" TEXT,
	"invoice_country" TEXT,
	"invoice_link" TEXT,
	"invoice_text" TEXT,
	PRIMARY KEY("invoice_id")
);
CREATE TABLE "item" (
	"item_id" TEXT NOT NULL UNIQUE,
	"invoice_id" TEXT NOT NULL,
	"item_name" TEXT,
	"item_price" REAL,
	"item_price_one" REAL,
	"item_quantity" REAL,
	"item_photo" TEXT,
	PRIMARY KEY("item_id"),
	FOREIGN KEY("invoice_id") REFERENCES "invoice"("invoice_id")
);
CREATE TABLE "tag" (
    "tag_id" INTEGER,
    "tag_name" TEXT NOT NULL,
		PRIMARY KEY ("tag_id" AUTOINCREMENT)
);

CREATE TABLE "item_tag" (
	"item_id"	TEXT NOT NULL UNIQUE,
	"tag_id"	NUMBER NOT NULL,
	PRIMARY KEY("item_id","tag_id"),
	FOREIGN KEY("tag_id") REFERENCES "tag"("tag_id"),
	FOREIGN KEY("item_id") REFERENCES "item"("item_id")
);
CREATE TABLE "invoice_tag" (
    "invoice_id" TEXT NOT NULL UNIQUE,
    "tag_id" NUMBER NOT NULL,
    PRIMARY KEY ("invoice_id", "tag_id"),
    FOREIGN KEY ("invoice_id") REFERENCES invoice("invoice_id"),
    FOREIGN KEY ("tag_id") REFERENCES tag("tag_id")
);
CREATE TABLE "exchange_rate_eur" (
	"exchange_rate_eur_id" INTEGER,
	"exchange_rate_eur_date" TEXT NOT NULL,
	"exchange_rate_eur_currency" TEXT NOT NULL,
	"exchange_rate_eur_value" REAL NOT NULL,
	PRIMARY KEY("exchange_rate_eur_id" AUTOINCREMENT)
);
CREATE TABLE "migration" (
	"name" TEXT UNIQUE,
	PRIMARY KEY("name")
);