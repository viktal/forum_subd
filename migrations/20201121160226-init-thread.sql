-- +migrate Up
set search_path to main;

create table thread
(
    thread_id SERIAL not null
        constraint thread_pkey primary key,
    forum_id SERIAL not null
            references main.forum(forum_id),
    forum text,
    user_id SERIAL not null
            references main.users(user_id),
    nickname text,
    title text,
    message text,
    slug text unique,
    create_date timestamp,
    votes numeric default 0
);

-- +migrate Down
drop table main.thread;
