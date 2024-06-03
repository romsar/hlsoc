create table if not exists posts
(
    id          uuid       not null primary key,
    text        text       not null,
    created_by  uuid       not null references users(id) on delete cascade on update cascade,
    created_at  timestamp  not null
);

