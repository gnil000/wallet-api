-- +goose Up
insert into wallet (balance, last_operation) values
  (1000, 'deposit'),
  (2500, 'deposit'),
  (500,  'withdraw'),
  (7500, 'deposit'),
  (300,  'withdraw');