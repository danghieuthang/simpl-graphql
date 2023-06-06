

CREATE TABLE public.audits2 (
	created_at timestamptz NULL,
	last_modified_at timestamptz NULL,
	is_deleted bool NULL,
	id bigserial NOT NULL
);