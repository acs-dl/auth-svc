-- +migrate Up

create table users (
    id bigserial primary key,
    email text not null,
    password text not null
);

INSERT INTO users (email, password)
VALUES ('nik@gmail.com', '$2b$10$ggulBRryhFGQEbaPX76oGeZ1EgduENOtSZWSe3d693z27X33Zt4Xe');

create table refresh_tokens (
    token text primary key,
    owner_id int not null references users (id) on delete cascade,
    valid_date int not null
);

create table amounts (
    access int not null default 0,
    refresh int not null default 0
);
INSERT INTO amounts (access, refresh)
VALUES (0, 0);

create table modules (
    id bigserial primary key,
    module_name text not null,
    permission text not null
);

create table modules_users (
    module_id int not null,
    user_id int not null
);


-- +migrate Down

drop table refresh_tokens;
drop table modules_users;
drop table users;
drop table amounts;
drop table modules;
