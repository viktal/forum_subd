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
    slug citext,
    create_date timestamp WITH TIME ZONE,
    votes numeric default 0
);

create unique index uniq_slug on main.thread (slug) where slug is not null;

-- +migrate Down
drop table main.thread;
drop index if exists uniq_slug;
