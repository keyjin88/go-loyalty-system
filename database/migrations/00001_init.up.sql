create table users
(
    id         serial
        primary key                    not null,
    created_at timestamp default now() not null,
    updated_at timestamp default now() not null,
    username   varchar                 not null
        unique,
    password   varchar                 not null
);