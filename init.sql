create schema if not exists main;
set search_path to main;

CREATE EXTENSION IF NOT EXISTS CITEXT WITH SCHEMA main;
------------------------------------------------------------------------------------
drop table if exists forum cascade;
drop table if exists post cascade;
drop table if exists thread cascade;
drop table if exists users cascade;
drop table if exists vote cascade;
drop table if exists forum_users cascade;

drop type if exists voice_types;

create type voice_types as enum ('-1', '1');
------------------------------------------------------------------------------------
create unlogged table users
(
    user_id  SERIAL             not null
        constraint users_pkey primary key,
    email    citext collate "C" not null unique,
    nickname citext collate "C" not null unique,
    fullname text               not null,
    about    text
);

create unlogged table forum
(
    forum_id SERIAL not null
        constraint forum_pkey primary key,
    slug     citext collate "C" unique,
    user_id  SERIAL not null
        references users (user_id),
    author   citext collate "C",
    title    text   not null,
    threads  numeric default 0,
    posts    numeric default 0

);

create unlogged table forum_users
(
    user_id SERIAL references users (user_id) on delete cascade         not null,
    forum   citext collate "C" references forum (slug) on delete cascade not null
);

create unlogged table thread
(
    thread_id   SERIAL             not null
        constraint thread_pkey primary key,
    forum_id    SERIAL             not null
        references main.forum (forum_id),
    forum       CITEXT collate "C" not null,
    user_id     SERIAL             not null
        references main.users (user_id),
    nickname    CITEXT collate "C" not null,
    title       text,
    message     text,
    slug        citext collate "C",
    create_date timestamp WITH TIME ZONE,
    votes       numeric default 0
);
create unique index slug on main.thread (slug) where slug is not null;


create unlogged table post
(
    post_id   SERIAL             not null
        constraint post_pkey primary key,
    forum_id  SERIAL             not null references forum (forum_id),
    forum     CITEXT collate "C" not null,
    user_id   SERIAL             not null references users (user_id),
    author    CITEXT collate "C" not null,
    thread_id SERIAL             not null references thread (thread_id),
    thread    citext collate "C",
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

CREATE INDEX if not exists user_nickname ON main.users using hash (nickname);
CREATE INDEX if not exists user_email ON main.users using hash (email);

CREATE INDEX if not exists forum_slug ON main.forum using hash (slug);

create unique index if not exists forum_users_unique on forum_users (forum, user_id);
cluster forum_users using forum_users_unique;

CREATE INDEX if not exists thr_slug ON main.thread using hash (slug);
CREATE INDEX if not exists thr_date ON main.thread (create_date);
CREATE INDEX if not exists thr_forum ON main.thread using hash (forum);
CREATE INDEX if not exists thr_forum_date ON main.thread (forum, create_date);

create index if not exists post_id_path on main.post (post_id, (path[1]));
create index if not exists post_thread_id_path1_parent on main.post (thread, post_id, (path[1]), parent);
create index if not exists post_thread_path_id on main.post (thread, path, post_id);
create index if not exists post_path1 on main.post ((path[1]));
create index if not exists post_thread_id on main.post (thread, post_id);
CREATE INDEX if not exists post_thr_id ON main.post (thread_id);

create unique index if not exists vote_unique on main.vote (user_id, thread_id);

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

------------------------------------------------------------------------------------

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

------------------------------------------------------------------------------------

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

------------------------------------------------------------------------------------

create or replace function add_forum_user_on_thread_create()
    returns trigger
    language plpgsql as
$BODY$
begin
    insert into main.forum_users (user_id, forum) values (new.user_id, new.forum) on conflict do nothing;
    return new;
end
$BODY$;

CREATE TRIGGER update_forum_users_new_thread
    AFTER INSERT
    ON main.thread
    FOR EACH ROW
EXECUTE PROCEDURE add_forum_user_on_thread_create();