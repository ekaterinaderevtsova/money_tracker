CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE spendings (
	"uuid" uuid DEFAULT uuid_generate_v4() NOT NULL,
	"date" date NOT NULL,
	sum int4 NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
	CONSTRAINT spendings_pkey PRIMARY KEY (uuid)
);
