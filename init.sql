create schema if not exists main;
set search_path to main;

CREATE EXTENSION IF NOT EXISTS citext;
------------------------------------------------------------------------------------
drop table if exists forum cascade;
drop table if exists post cascade;
drop table if exists thread cascade;
drop table if exists users cascade;
drop table if exists vote cascade;

drop type if exists voice_types;

create type voice_types as enum ('-1', '1');
------------------------------------------------------------------------------------
create unlogged table users
(
    user_id      SERIAL       not null
        constraint users_pkey primary key,
    email        citext       not null unique,
    nickname     citext       not null unique,
    nickname_byt bytea GENERATED ALWAYS AS (lower(nickname)::bytea) STORED,
    fullname     varchar(128) not null,
    about        text
);

create unlogged table forum
(
    forum_id SERIAL not null
        constraint forum_pkey primary key,
    slug     citext unique,
    user_id  SERIAL not null
        references users (user_id),
    author   text,
    title    text   not null,
    threads  numeric default 0,
    posts    numeric default 0

);

create unlogged table thread
(
    thread_id   SERIAL not null
        constraint thread_pkey primary key,
    forum_id    SERIAL not null
        references main.forum (forum_id),
    forum       text,
    user_id     SERIAL not null
        references main.users (user_id),
    nickname    text,
    title       text,
    message     text,
    slug        citext,
    create_date timestamp WITH TIME ZONE,
    votes       numeric default 0
);
create unique index uniq_slug on main.thread (slug) where slug is not null;


create unlogged table post
(
    post_id   SERIAL not null
        constraint post_pkey primary key,
    forum_id  SERIAL not null references forum (forum_id),
    forum     text,
    user_id   SERIAL not null references users (user_id),
    author    text,
    thread_id SERIAL not null references thread (thread_id),
    thread    text,
    message   text,
    is_edited bool  default false,
    created   timestamp WITH TIME ZONE,

    parent    int   default 0,
    path      int[] default array []::int[]
);


create unlogged table vote
(
    vote_id   SERIAL not null
        constraint vote_pkey primary key,
    user_id   SERIAL not null
        references users (user_id),
    thread_id SERIAL not null
        references thread (thread_id),
    voice     voice_types
);
------------------------------------------------------------------------------------

CREATE INDEX if not exists uniq_user_nickname ON main.users (lower(nickname) text_pattern_ops);
CREATE INDEX if not exists uniq_user_email ON main.users (lower(email) text_pattern_ops);
CREATE INDEX if not exists uniq_user_nickname_byt ON main.users (nickname_byt);

CREATE INDEX if not exists uniq_forum_slug ON main.forum (lower(slug) text_pattern_ops);
CREATE INDEX if not exists uniq_forum_auth ON main.forum (lower(author) text_pattern_ops);

CREATE INDEX if not exists uniq_thr_slug ON main.thread (lower(slug) text_pattern_ops);
CREATE INDEX if not exists uniq_thr_forum ON main.thread (lower(forum) text_pattern_ops);
CREATE INDEX if not exists uniq_thr_nick ON main.thread (lower(nickname) text_pattern_ops);
CREATE INDEX if not exists uniq_thr_date ON main.thread (create_date);

CREATE INDEX if not exists uniq_pos_date ON main.post (created);
CREATE INDEX if not exists uniq_pos_thr_id ON main.post (thread_id);
CREATE INDEX if not exists uniq_pos_forum_id ON main.post (forum_id);
CREATE INDEX if not exists uniq_pos_forum ON main.post (lower(forum) text_pattern_ops);
CREATE INDEX if not exists uniq_pos_author ON main.post (lower(author) text_pattern_ops);
CREATE INDEX if not exists uniq_pos_thread ON main.post (lower(thread) text_pattern_ops);

create unique index if not exists uniq_vote on main.vote (user_id, thread_id);
CREATE INDEX if not exists vote_voice ON main.vote (voice);

------------------------------------------------------------------------------------

create or replace function update_cnt_vote_thread()
    returns trigger
    language plpgsql as
$BODY$
begin
    if TG_OP = 'INSERT' then
        if new.voice = '1' then
            update main.thread set votes = votes + 1 where thread_id = new.thread_id;
        else
            update main.thread set votes = votes - 1 where thread_id = new.thread_id;
        end if;
    elseif TG_OP = 'UPDATE' then
        if new.voice = '1' then
            update main.thread set votes = votes + 2 where thread_id = new.thread_id;
        else
            update main.thread set votes = votes - 2 where thread_id = new.thread_id;
        end if;
    end if;
    return new;
end
$BODY$;

CREATE TRIGGER add_vote_to_thread
    AFTER INSERT OR UPDATE
    ON main.vote
    FOR EACH ROW
EXECUTE PROCEDURE update_cnt_vote_thread();

CREATE OR REPLACE FUNCTION update_post_path()
    returns trigger
    language plpgsql as
$BODY$
BEGIN
    IF (new.parent = 0) THEN
        new.path = new.path || new.post_id;
    ELSE
        new.path = (select path from main.post where post_id = new.parent) || new.post_id;
    END IF;
    RETURN new;
END
$BODY$;

CREATE TRIGGER add_post_to_thread
    BEFORE INSERT
    ON main.post
    FOR EACH ROW
EXECUTE PROCEDURE update_post_path();

create or replace function update_cnt_thread()
    returns trigger
    language plpgsql as
$BODY$
begin
    if TG_OP = 'INSERT' then
        update main.forum set threads = threads + 1 where forum_id = new.forum_id;
    end if;
    return new;
end
$BODY$;

CREATE TRIGGER count_thread_forum
    AFTER INSERT
    ON main.thread
    FOR EACH ROW
EXECUTE PROCEDURE update_cnt_thread();

create or replace function update_cnt_post()
    returns trigger
    language plpgsql as
$BODY$
begin
    if TG_OP = 'INSERT' then
        update main.forum set posts = posts + 1 where forum_id = new.forum_id;
    end if;
    return new;
end
$BODY$;

CREATE TRIGGER count_post_forum
    AFTER INSERT
    ON main.post
    FOR EACH ROW
EXECUTE PROCEDURE update_cnt_post();
