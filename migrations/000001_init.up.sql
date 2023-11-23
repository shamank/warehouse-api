create table products
(
    uuid    uuid primary key default gen_random_uuid(),
    name    varchar,
    size    varchar,
    article varchar,

    unique (article)
);

create index idx_article on products (article);

create table warehouses
(
    uuid         uuid primary key default gen_random_uuid(),
    name         varchar,
    is_available boolean
);

create table warehouse_products
(
    warehouse_uuid    uuid,
    product_uuid      uuid,
    quantity          int,
    reserved_quantity int,

    primary key (warehouse_uuid, product_uuid),
    foreign key (warehouse_uuid) references warehouses (uuid),
    foreign key (product_uuid) references products (uuid),

    constraint check_quantity check (quantity >= 0),
    constraint check_reserved_quantity check (reserved_quantity >= 0)
);





