-- +migrate Up
set search_path to main;

create table post
(
    post_id SERIAL not null
        constraint post_pkey primary key,
    forum_id SERIAL not null
            references forum(forum_id),
    forum text,
    user_id SERIAL not null
            references users(user_id),
    author text,
    thread_id SERIAL not null
        references thread(thread_id),
    thread text,
    message text,
    parent numeric,
    is_edited bool default false,
    created timestamp
);

-- +migrate Down
drop table main.post;
