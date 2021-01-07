
-- +migrate Up
set search_path to main;

create table forum
(
    forum_id SERIAL not null
        constraint forum_pkey primary key,
    slug citext unique,
    user_id SERIAL not null
            references users(user_id),
    author text,
    title text not null

);
-- +migrate Down
drop table main.forum;
