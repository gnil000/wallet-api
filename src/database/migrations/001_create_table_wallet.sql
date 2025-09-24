-- +goose Up
create table if not exists wallet (
    id uuid default gen_random_uuid() primary key,
    balance bigint,
    created timestamp not null default now(),
    updated timestamp not null default now(),
    last_operation text
);