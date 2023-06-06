CREATE TABLE public.audits (
	created_at timestamptz NULL,
	last_modified_at timestamptz NULL,
	is_deleted bool NULL,
	id bigserial NOT NULL
);