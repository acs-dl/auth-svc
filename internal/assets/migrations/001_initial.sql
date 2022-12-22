-- +migrate Up

CREATE TABLE IF NOT EXISTS users (
    id bigserial primary key,
    email text not null,
    password text not null
);

INSERT INTO users (email, password)
VALUES ('nik@gmail.com', '$2b$10$ggulBRryhFGQEbaPX76oGeZ1EgduENOtSZWSe3d693z27X33Zt4Xe');

CREATE TABLE IF NOT EXISTS refresh_tokens (
    token text primary key,
    owner_id int not null references users (id) on delete cascade,
    valid_date int not null
);

CREATE TABLE IF NOT EXISTS modules (
    id bigserial PRIMARY KEY,
    name text UNIQUE NOT NULL
);

CREATE INDEX IF NOT EXISTS module_namex ON modules(name);

CREATE TABLE IF NOT EXISTS permissions (
    id bigserial PRIMARY KEY,
    module_id bigint NOT NULL,
    name text UNIQUE NOT NULL,

    foreign key(module_id) references modules(id) on delete cascade
);

CREATE TABLE IF NOT EXISTS permissions_users (
    permission_id bigint not null,
    user_id bigint not null,

    foreign key(permission_id) references permissions(id) on delete cascade,
    foreign key(user_id) references users(id) on delete cascade
);

CREATE INDEX IF NOT EXISTS user_idx on permissions_users(user_id);

-- +migrate Down

DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS permissions_users;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS modules;
