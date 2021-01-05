-- +migrate Up
set search_path to main;

-- +migrate StatementBegin

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
    end if;
    return new;
end
$BODY$;

-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE TRIGGER add_vote_to_thread
    AFTER INSERT ON main.vote FOR EACH ROW EXECUTE PROCEDURE update_cnt_vote_thread();

-- +migrate StatementEnd

-- +migrate Down
drop trigger if exists add_vote_to_thread on main.thread;
drop function if exists main.update_cnt_vote_thread();