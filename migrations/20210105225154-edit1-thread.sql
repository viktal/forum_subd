-- +migrate Up
set search_path to main;

alter table main.thread drop column slug;
alter table main.thread add column slug text;
create unique index uniq_slug on main.thread (slug) where slug is not null;

-- +migrate Down

drop index if exists uniq_slug;