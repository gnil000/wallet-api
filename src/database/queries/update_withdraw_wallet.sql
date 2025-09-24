select withdraw_balance($1, $2);

-- select balance from wallet where id = $1 for update;
-- update wallet set balance = balance - $2, updated = now(), last_operation = 'withdraw' where id = $1;