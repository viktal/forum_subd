-- +migrate Up
set search_path to main;

create type voice_types as enum ('-1', '1');

create table vote
(
    vote_id SERIAL not null
        constraint vote_pkey primary key,
    user_id SERIAL not null
            references users(user_id),
    thread_id SERIAL not null
        references thread(thread_id),
    voice voice_types
);

create unique index uniq_vote on main.vote (user_id, thread_id);

-- +migrate Down
drop type if exists voice_types;
drop table main.vote;
drop index if exists uniq_vote;


