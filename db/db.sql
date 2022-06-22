DROP SCHEMA IF EXISTS dbproject CASCADE;
CREATE SCHEMA dbproject;
CREATE EXTENSION IF NOT EXISTS citext;

DROP TABLE IF EXISTS dbproject."User" CASCADE;
DROP TABLE IF EXISTS dbproject."Post" CASCADE;
DROP TABLE IF EXISTS dbproject."Thread" CASCADE;
DROP TABLE IF EXISTS dbproject."Forum" CASCADE;
DROP TABLE IF EXISTS dbproject."Vote" CASCADE;
DROP TABLE IF EXISTS dbproject."Users_by_Forum" CASCADE;

CREATE UNLOGGED TABLE dbproject."User"
(
    Id SERIAL PRIMARY KEY,
    NickName CITEXT UNIQUE NOT NULL,
    FullName TEXT NOT NULL,
    About TEXT,
    Email CITEXT UNIQUE NOT NULL
);

CREATE UNLOGGED TABLE dbproject."Forum"
(
    Id SERIAL PRIMARY KEY,
    Title TEXT NOT NULL,
    "user" CITEXT REFERENCES dbproject."User"(NickName) NOT NULL,
    Slug CITEXT UNIQUE NOT NULL,
    Posts INT,
    Threads INT
);

CREATE UNLOGGED TABLE dbproject."Thread"
(
    Id SERIAL PRIMARY KEY,
    Title TEXT NOT NULL,
    Author CITEXT REFERENCES dbproject."User"(NickName) NOT NULL,
    Forum CITEXT REFERENCES dbproject."Forum"(Slug) NOT NULL,
    Message TEXT NOT NULL,
    Votes INT,
    Slug CITEXT UNIQUE DEFAULT citext(1),
    Created TIMESTAMP WITH TIME ZONE
);


CREATE UNLOGGED TABLE dbproject."Post"
(
    Id SERIAL PRIMARY KEY,
    Parent INT DEFAULT 0,
    Author CITEXT REFERENCES dbproject."User"(NickName) NOT NULL,
    Message TEXT NOT NULL,
    IsEdited bool NOT NULL DEFAULT FALSE,
    Forum CITEXT REFERENCES dbproject."Forum"(Slug) NOT NULL,
    Thread INT REFERENCES dbproject."Thread"(Id) NOT NULL,
    Created TIMESTAMP WITH TIME ZONE DEFAULT now(),
    Path INT[] DEFAULT ARRAY []::INTEGER[]
);

CREATE UNLOGGED TABLE dbproject."Users_by_Forum"
(
    Id SERIAL PRIMARY KEY,
    Forum CITEXT NOT NULL,
    "user" CITEXT REFERENCES dbproject."User"(NickName) NOT NULL,
    CONSTRAINT onlyOneUser UNIQUE (Forum, "user")
);

CREATE UNLOGGED TABLE dbproject."Vote"
(
    Id SERIAL PRIMARY KEY,
    ThreadId INT REFERENCES dbproject."Thread"(id) NOT NULL,
    "user" CITEXT REFERENCES dbproject."User"(NickName) NOT NULL,
    Value INT NOT NULL,
    CONSTRAINT onlyOneVote UNIQUE (ThreadId, "user")
);

-- adding a new voice
CREATE OR REPLACE FUNCTION addNewVoice() RETURNS TRIGGER AS $$
BEGIN
    UPDATE dbproject."Thread" t SET votes = t.votes + NEW.Value WHERE t.Id = NEW.threadid;
    RETURN NULL;
END
$$ LANGUAGE 'plpgsql';

CREATE TRIGGER voiceTrigger
    AFTER INSERT ON dbproject."Vote"
    FOR EACH ROW EXECUTE PROCEDURE addNewVoice();

-- changing voice
CREATE OR REPLACE FUNCTION changeVoice() RETURNS TRIGGER AS $$
BEGIN
    IF OLD.value <> NEW.value
    THEN UPDATE dbproject."Thread" t SET votes = (t.votes + NEW.value * 2) WHERE t.Id = NEW.threadid;
    END IF;
    RETURN NEW;
END
$$ LANGUAGE 'plpgsql';

CREATE TRIGGER voiceUpdateTrigger
    AFTER UPDATE ON dbproject."Vote"
    FOR EACH ROW EXECUTE PROCEDURE changeVoice();

-- add new thread
CREATE OR REPLACE FUNCTION incForumThreads() RETURNS TRIGGER AS $$
BEGIN
    UPDATE dbproject."Forum" SET threads = threads + 1 WHERE NEW.Forum = slug;
    INSERT INTO dbproject."Users_by_Forum" (forum, "user") VALUES (NEW.Forum, NEW.Author)
    ON CONFLICT DO NOTHING;
    RETURN NULL;
END
$$ LANGUAGE 'plpgsql';

CREATE TRIGGER createThreadTrigger
    AFTER INSERT ON dbproject."Thread"
    FOR EACH ROW EXECUTE PROCEDURE incForumThreads();

-- adding a post
CREATE OR REPLACE FUNCTION addPost() RETURNS TRIGGER AS $$
BEGIN
    --  increase counter
    UPDATE dbproject."Forum" SET posts = posts + 1 WHERE Slug = NEW.forum;
--  add user to table forum-user
    INSERT INTO dbproject."Users_by_Forum" (forum, "user") VALUES (NEW.forum, NEW.author)
    ON CONFLICT DO NOTHING;
--  write path
    NEW.path = (SELECT P.path FROM dbproject."Post" P WHERE P.id = NEW.parent LIMIT 1) || NEW.id;
    RETURN NEW;
END
$$ LANGUAGE 'plpgsql';

CREATE TRIGGER addPost
    BEFORE INSERT ON dbproject."Post"
    FOR EACH ROW EXECUTE PROCEDURE addPost();

CREATE INDEX IF NOT EXISTS postPath ON dbproject."Post" (path);
CREATE INDEX IF NOT EXISTS postPath1 ON dbproject."Post" ((path[1]));
CREATE INDEX IF NOT EXISTS postIdPath1 ON dbproject."Post" (id, (path[1]));
CREATE INDEX IF NOT EXISTS postForum ON dbproject."Post" (forum);
CREATE INDEX IF NOT EXISTS postThread ON dbproject."Post" (thread);

CREATE INDEX IF NOT EXISTS userNick ON dbproject."User" USING hash (nickname);
CREATE INDEX IF NOT EXISTS userEmail ON dbproject."User" USING hash(email);
CREATE INDEX IF NOT EXISTS forumUsersUser ON dbproject."Users_by_Forum" USING hash ("user");

CREATE INDEX IF NOT EXISTS forumSlug ON dbproject."Forum" USING hash(slug);
CREATE INDEX IF NOT EXISTS threadSlug ON dbproject."Thread" USING hash(slug);
CREATE INDEX IF NOT EXISTS threadForum ON dbproject."Thread" (forum);
CREATE INDEX IF NOT EXISTS threadCreated ON dbproject."Thread" (created);
CREATE INDEX IF NOT EXISTS threadCreatedForum ON dbproject."Thread" (forum, created);

CREATE UNIQUE INDEX IF NOT EXISTS votesNicknameThreadNickname ON dbproject."Vote" (ThreadId, "user");