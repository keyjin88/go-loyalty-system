BEGIN TRANSACTION;
create table users
(
    id         serial
        primary key                    not null,
    created_at timestamp default now() not null,
    updated_at timestamp default now() not null,
    user_name  varchar                 not null
        unique,
    password   varchar                 not null,
    is_deleted boolean   default false not null
);

create table orders
(
    id         serial
        primary key                    not null,
    created_at timestamp default now() not null,
    updated_at timestamp default now() not null,
    number     varchar                 not null
        unique,
    user_id    integer                 not null,
    is_deleted boolean   default false not null
);
COMMIT;