create table if not exists user_friends
(
    user_id       uuid    not null references users(id) on delete cascade on update cascade,
    friend_id     uuid    not null references users(id) on delete cascade on update cascade,

    PRIMARY KEY (user_id, friend_id)
);

