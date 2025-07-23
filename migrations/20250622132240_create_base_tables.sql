-- +goose Up
create type user_role as enum('ADMIN', 'USER');

create table "user" (
    id uuid unique not null
    , email varchar (100) not null
    , password varchar (255) not null
    , role user_role default 'USER'::user_role
    , created_at timestamp not null default now()
    , updated_at timestamp not null default now()
    , deleted_at timestamp
    , primary key (id)
    , constraint email_unique unique (email)
);


-- +goose Down
drop table "user";
drop type user_role;
