create database productdb;
use productdb;

create table productdb.Products
(
    id       int auto_increment primary key,
    model    varchar(30) not null,
    company  varchar(30) not null,
    quantity int         default 0,
    price    decimal         not null
);

insert into productdb.Products (model, company, quantity, price)
values ('iPhone X', 'Apple', 74, 10000),
       ('Pixel 2', 'Google', 62, 22000),
       ('Galaxy S9', 'Samsung', 65, 22000),
       ('Xaiomi', 'redmi', 37, 23000),
       ('S21','Samsung', 22, 21222);

alter table productdb.Products add column quantity int not null;

create table productdb.ProductsFeatures
(
    id              int auto_increment primary key,
    product_id      int not null,
    cpu             int not null,
    memory          int not null,
    display_size    int         default 0,
    camera          decimal     not null,
    UNIQUE KEY      (product_id),
    CONSTRAINT      fk_product_id FOREIGN KEY (product_id) REFERENCES Products(id) ON DELETE CASCADE
);