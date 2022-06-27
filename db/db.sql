CREATE EXTENSION IF NOT EXISTS citext;

DROP TABLE users;
DROP TABLE forums;
DROP TABLE threads;
DROP TABLE posts;
DROP TABLE forumToUsers;
DROP TABLE votes;

DROP INDEX for_search_by_slug;
DROP INDEX for_search_by_forum;
DROP INDEX for_search_threads_on_forum;
DROP INDEX for_tree_search;
DROP INDEX for_parent_tree_search;
DROP INDEX user_nickname_hash;
DROP INDEX search_user_vote;
DROP INDEX forum_slug_hash;
DROP INDEX forum_users_forum;


CREATE UNLOGGED TABLE if not exists users
(
    nickname citext COLLATE "C" NOT NULL PRIMARY KEY,
    fullname text               NOT NULL,
    about    text,
    email    citext UNIQUE
);

CREATE UNLOGGED TABLE if not exists forums
(
    slug    citext             NOT NULL PRIMARY KEY,
    title   text               NOT NULL,
    "user"  citext COLLATE "C" NOT NULL REFERENCES users (nickname),
    posts   bigint DEFAULT 0,
    threads bigint DEFAULT 0
);

CREATE UNLOGGED TABLE if not exists threads
(
    id      bigserial          NOT NULL PRIMARY KEY,
    title   text               NOT NULL,
    author  citext COLLATE "C" NOT NULL REFERENCES users (nickname),
    forum   citext             NOT NULL REFERENCES forums (slug),
    message text               NOT NULL,
    votes   integer     DEFAULT 0,
    slug    citext             NOT NULL,
    created timestamptz DEFAULT now()
);

CREATE UNLOGGED TABLE if not exists posts
(
    id          bigserial          NOT NULL PRIMARY KEY,
    parent      integer     DEFAULT 0,
    author      citext COLLATE "C" NOT NULL REFERENCES users (nickname),
    message     text               NOT NULL,
    isEdited    boolean     DEFAULT false,
    forum       citext             NOT NULL REFERENCES forums (slug),
    thread      integer REFERENCES threads (id),
    created     timestamptz DEFAULT now(),
    parent_path BIGINT[]    DEFAULT ARRAY []::integer[]
);

CREATE UNLOGGED TABLE IF NOT EXISTS forumToUsers
(
    nickname citext COLLATE "C" NOT NULL REFERENCES users (nickname),
    fullname text               NOT NULL,
    about    text,
    email    citext             NOT NULL,
    forum    citext             NOT NULL REFERENCES forums (slug),
    PRIMARY KEY (nickname, forum)
);

CREATE UNLOGGED TABLE if not exists votes
(
    nickname citext COLLATE "C" NOT NULL REFERENCES users (nickname),
    thread   serial             NOT NULL REFERENCES threads (id),
    voice    integer            NOT NULL,
    PRIMARY KEY (nickname, thread)
);


CREATE OR REPLACE FUNCTION forum_users_update() RETURNS TRIGGER AS
$forum_users_update$
DECLARE
    nickname_param citext;
    fullname_param text;
    about_param    text;
    email_param    citext;
BEGIN
    SELECT t.nickname, t.fullname, t.about, t.email
    FROM users AS t
    WHERE t.nickname = new.author
    INTO nickname_param, fullname_param, about_param, email_param;

    INSERT INTO forumToUsers (nickname, fullname, about, email, forum)
    VALUES (nickname_param, fullname_param, about_param, email_param, new.forum)
    ON CONFLICT DO NOTHING;

    return new;
END;
$forum_users_update$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION set_post_parent_path() RETURNS TRIGGER AS
$set_post_parent_path$
BEGIN
    new.parent_path = (SELECT parent_path FROM posts WHERE id = new.parent) || new.id;
    return new;
END;
$set_post_parent_path$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION add_forum_thread_count() RETURNS TRIGGER AS
$add_forum_thread_count$
BEGIN
    UPDATE forums SET threads = forums.threads + 1 WHERE slug = new.forum;
    return new;
END;
$add_forum_thread_count$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION add_forum_posts_count() RETURNS TRIGGER AS
$add_forum_posts_count$
BEGIN
    UPDATE forums SET posts = forums.posts + 1 WHERE slug = new.forum;
    return new;
END;
$add_forum_posts_count$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION add_thread_vote() RETURNS TRIGGER AS
$add_thread_vote$
BEGIN
    UPDATE threads SET votes = threads.votes + new.voice WHERE id = new.thread;
    return new;
END;
$add_thread_vote$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_thread_vote() RETURNS TRIGGER AS
$update_thread_vote$
BEGIN
    UPDATE threads SET votes = threads.votes + (new.voice - old.voice) WHERE id = new.thread;
    return new;
END;
$update_thread_vote$ LANGUAGE plpgsql;


CREATE TRIGGER forum_users_for_post
    AFTER INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE forum_users_update();

CREATE TRIGGER forum_users_for_thread
    AFTER INSERT
    ON threads
    FOR EACH ROW
EXECUTE PROCEDURE forum_users_update();

CREATE TRIGGER set_post_parent_path_trigger
    BEFORE INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE set_post_parent_path();

CREATE TRIGGER add_forum_thread_count_trigger
    AFTER INSERT
    ON threads
    FOR EACH ROW
EXECUTE PROCEDURE add_forum_thread_count();

CREATE TRIGGER add_forum_posts_count_trigger
    AFTER INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE add_forum_posts_count();

CREATE TRIGGER add_thread_vote_trigger
    AFTER INSERT
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE add_thread_vote();

CREATE TRIGGER update_thread_vote_trigger
    AFTER UPDATE
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE update_thread_vote();



CREATE INDEX IF NOT EXISTS for_search_by_slug ON threads USING hash (slug);
CREATE INDEX IF NOT EXISTS for_search_by_forum ON threads USING hash (forum);
CREATE INDEX IF NOT EXISTS for_search_threads_on_forum ON threads (forum, created);

CREATE INDEX IF NOT EXISTS for_search_users_on_forum_posts ON posts (forum, author);
CREATE INDEX IF NOT EXISTS for_flat_search ON posts (thread, id);
CREATE INDEX IF NOT EXISTS for_tree_search ON posts (thread, parent_path);
CREATE INDEX IF NOT EXISTS for_parent_tree_search ON posts ((parent_path[1]), parent_path);
CREATE INDEX IF NOT EXISTS post_id_hash ON posts using hash (id);
CREATE INDEX IF NOT EXISTS post_thread_hash ON posts using hash (thread);

CREATE INDEX IF NOT EXISTS user_nickname_hash ON users using hash (nickname);
CREATE INDEX IF NOT EXISTS  user_nickname_email ON users (nickname, email);

CREATE INDEX IF NOT EXISTS search_user_vote ON votes (nickname, thread, voice);

CREATE INDEX IF NOT EXISTS forum_slug_hash ON forums using hash (slug);

CREATE INDEX IF NOT EXISTS forum_users_forum ON forumToUsers (forum, nickname);

VACUUM ANALYZE;
