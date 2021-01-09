
-- +migrate Up
set search_path to main;

CREATE UNIQUE INDEX ON main.users (lower(nickname) text_pattern_ops);
CREATE UNIQUE INDEX ON main.users (lower(email) text_pattern_ops);

CREATE INDEX ON main.forum (lower(slug) text_pattern_ops);
CREATE INDEX ON main.forum (lower(author) text_pattern_ops);

CREATE INDEX ON main.thread (lower(slug) text_pattern_ops);
CREATE INDEX ON main.thread (lower(forum) text_pattern_ops);
CREATE INDEX ON main.thread (lower(nickname) text_pattern_ops);
CREATE INDEX ON main.thread (create_date);

CREATE INDEX ON main.post (created);
CREATE INDEX ON main.post (thread_id);
CREATE INDEX ON main.post (forum_id);
CREATE INDEX ON main.post (lower(forum) text_pattern_ops);
CREATE INDEX ON main.post (lower(author) text_pattern_ops);
CREATE INDEX ON main.post (lower(thread) text_pattern_ops);

CREATE UNIQUE INDEX ON main.vote (user_id, thread_id);
CREATE INDEX ON main.vote (voice);

-- +migrate Down
