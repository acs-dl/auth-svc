-- +migrate Up

create type users_status_enum as enum ('super_admin', 'admin', 'user');

create table if not exists users (
    id bigint primary key,
    email text not null,
    password text not null,
    status users_status_enum not null
);

insert into users (id, email, password, status)
values (1, 'serhii.pomohaiev@distributedlab.com', '$2b$10$ggulBRryhFGQEbaPX76oGeZ1EgduENOtSZWSe3d693z27X33Zt4Xe', 'super_admin');

create table if not exists refresh_tokens (
    token text primary key,
    owner_id int not null references users (id) on delete cascade,
    valid_till int not null
);

create index if not exists refresh_tokens_token_idx on refresh_tokens(token);

create table if not exists modules (
    id bigserial primary key,
    name text unique not null
);

create index if not exists module_name_idx on modules(name);

create table if not exists permissions (
    id bigserial primary key,
    module_id bigint not null,
    name text not null,
    status users_status_enum not null,

    unique(module_id, name, status),
    foreign key(module_id) references modules(id) on delete cascade
);

create index if not exists permissions_moduleid_idx on permissions(module_id);
create index if not exists permissions_status_idx on permissions(status);


-- +migrate Down

drop index if exists permissions_moduleid_idx;
drop index if exists permissions_status_idx;

drop table if exists permissions;

drop index if exists module_name_idx;

drop table if exists modules;

drop index if exists refresh_tokens_token_idx;

drop table if exists refresh_tokens;
drop table if exists users;

drop type if exists users_status_enum;