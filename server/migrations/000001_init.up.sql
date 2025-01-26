CREATE TABLE users
(
    id            serial      not null unique,
    username      varchar(32) not null unique,
    password_hash varchar(96) not null
);

CREATE TABLE lists
(
    id          serial      not null unique,
    title       varchar(32) not null
);

CREATE TABLE users_lists
(
    id      serial                                      not null unique,
    user_id int references users (id) on delete cascade not null,
    list_id int references lists (id) on delete cascade not null
);

CREATE TABLE tasks
(
    id          serial      not null unique,
    title       varchar(32) not null,
    done        boolean     not null default false
);


CREATE TABLE lists_tasks
(
    id      serial                                      not null unique,
    task_id int references tasks (id) on delete cascade not null,
    list_id int references lists (id) on delete cascade not null
);

CREATE TABLE tokens
(
    id            serial      not null unique,
    refresh_token varchar(64) not null,
    expires_at    timestamptz not null
);

CREATE TABLE users_tokens 
(
    id       serial                                       not null unique,
    user_id  int references users (id) on delete cascade  not null,
    token_id int references tokens (id) on delete cascade not null
);