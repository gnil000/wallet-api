select balance from wallet where id = $1 for update;
update wallet set balance = balance + $2, updated = now(), last_operation = 'deposit' where id = $1;