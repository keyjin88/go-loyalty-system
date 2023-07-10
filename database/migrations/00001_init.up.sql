BEGIN TRANSACTION;
create table users
(
    id         serial
        primary key                    not null,
    created_at timestamp default now() not null,
    updated_at timestamp default now() not null,
    user_name   varchar                 not null
        unique,
    password   varchar                 not null,
    is_deleted boolean   default false not null
);
COMMIT;