-- +migrate Up

create table users (
    id bigserial primary key,
    email text not null,
    password text not null
);

INSERT INTO users (email, password)
VALUES ('max@mail', '$2b$10$ggulBRryhFGQEbaPX76oGeZ1EgduENOtSZWSe3d693z27X33Zt4Xe');

create table refresh_tokens (
    token text primary key,
    owner_id int references users (id) on delete cascade,
    valid_date int
);

-- +migrate Down

drop table refresh_tokens;
drop table users;
