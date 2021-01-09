-- +migrate Up
set search_path to main;

-- +migrate StatementBegin

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

-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE TRIGGER count_thread_forum
    AFTER INSERT ON main.thread FOR EACH ROW EXECUTE PROCEDURE update_cnt_thread();

CREATE TRIGGER count_post_forum
    AFTER INSERT ON main.post FOR EACH ROW EXECUTE PROCEDURE update_cnt_post();

-- +migrate StatementEnd

-- +migrate Down
drop trigger if exists count_thread_forum on main.forum;
drop trigger if exists count_post_forum on main.forum;
drop function if exists main.update_cnt_thread();
drop function if exists main.update_cnt_post();
