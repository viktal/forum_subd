-- +migrate Up
create schema if not exists main;
set search_path to main;

CREATE EXTENSION IF NOT EXISTS citext;

create table users
(
    user_id SERIAL not null
        constraint users_pkey primary key,
    email citext not null unique,
    nickname citext not null unique,
    fullname varchar(128) not null,
    about text
);

-- +migrate Down
drop table main.users;
