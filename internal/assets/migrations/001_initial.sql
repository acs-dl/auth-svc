-- +migrate Up

CREATE TABLE IF NOT EXISTS users (
    id bigint PRIMARY KEY,
    email text NOT NULL,
    password text NOT NULL
);

INSERT INTO users (id, email, password)
VALUES (1, 'serhii.pomohaiev@distributedlab.com', '$2b$10$ggulBRryhFGQEbaPX76oGeZ1EgduENOtSZWSe3d693z27X33Zt4Xe');

CREATE TABLE IF NOT EXISTS refresh_tokens (
    token text PRIMARY KEY,
    owner_id INT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    valid_date INT NOT NULL
);

CREATE TABLE IF NOT EXISTS modules (
    id bigserial PRIMARY KEY,
    name text UNIQUE NOT NULL
);

CREATE INDEX IF NOT EXISTS module_namex ON modules(name);

CREATE TABLE IF NOT EXISTS permissions (
    id bigserial PRIMARY KEY,
    module_id bigint NOT NULL,
    name text NOT NULL,

    UNIQUE(module_id, name),
    FOREIGN KEY(module_id) REFERENCES modules(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS permissions_moduleid_name_idx ON permissions(module_id, name);

CREATE TABLE IF NOT EXISTS permissions_users (
    permission_id bigint NOT NULL,
    user_id bigint NOT NULL,

    FOREIGN KEY(permission_id) REFERENCES permissions(id) ON DELETE CASCADE,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS user_idx ON permissions_users(user_id);

-- +migrate Down

DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS permissions_users;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS modules;
