INSERT INTO public.roles
(id, "name", created_at, last_modified_at, is_deleted)
VALUES
    (nextval('roles_id_seq'::regclass), 'Role Name 1', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, false),
    (nextval('roles_id_seq'::regclass), 'Role Name 2', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, false),
    (nextval('roles_id_seq'::regclass), 'Role Name 3', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, false);


DO $$
DECLARE
i INT := 0;
BEGIN
FOR i IN 1..3 LOOP
        INSERT INTO public.users
        (id, "name", email, "password", role_id, created_at, last_modified_at, is_deleted, created_by, updated_by, "version")
        VALUES
        (nextval('users_id_seq'::regclass), 'Fake User ' || i, 'fakeuser' || i || '@example.com', 'password' || i, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, false, 'admin', 'admin', '1');
END LOOP;
END $$;
