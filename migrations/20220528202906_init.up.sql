create table if not exists image (
    uuid uuid primary key not null,
    path varchar(64) unique not null,
    ts  bigint not null
);

create table if not exists image_blurred (
    uuid uuid primary key not null,
    image_uuid uuid not null,
    x_0 int not null,
    y_0 int not null,
    x_1 int not null,
    y_1 int not null,
    ts  bigint not null
);