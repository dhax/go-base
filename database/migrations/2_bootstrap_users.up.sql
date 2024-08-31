INSERT INTO accounts (id, email, name, active, roles)
VALUES (DEFAULT, 'admin@example.com', 'Admin Example', true, '{admin}');

--bun:split

INSERT INTO accounts (id, email, name, active)
VALUES (DEFAULT, 'user@example.com', 'User Example', true);

