create table user (
    id varchar (40) unique not null
    , email varchar (100) not null
    , password varchar (255)  not null
    , role varchar (20) not null
    , created_at timestamp not null
    , updated_at timestamp not null
    , deleted_at timestamp
    , primary key (id)
    , constraint email_unique unique (email)
);
