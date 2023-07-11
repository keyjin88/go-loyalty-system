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
    balance    float     default 0.0   not null,
    withdrawn  float4    default 0.0   not null,
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
    status     varchar   default 'NEW' not null,
    accrual    float,
    is_deleted boolean   default false not null
);
COMMIT;