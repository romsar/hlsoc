create table if not exists users
(
    id          uuid    not null primary key,
    password    varchar not null,
    first_name  varchar not null,
    second_name varchar not null,
    birth_date  date    not null,
    gender      int     not null,
    biography   text,
    city        varchar not null
);

