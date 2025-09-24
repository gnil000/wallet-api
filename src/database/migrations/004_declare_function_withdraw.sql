-- +goose up
-- +goose statementbegin
create or replace function withdraw_balance(w_id uuid, w_amount bigint)
returns void as $$
declare
    current_balance bigint;
    balance_after bigint;
begin
    select balance into current_balance from wallet where id = w_id;

    if not found then
        raise exception 'wallet with id % not found', w_id
            using errcode = 'P0002';
    end if;

    if current_balance - w_amount < 0 then 
        raise exception 'not enough balance';
    end if;

    update wallet set balance = balance - $2, updated = now(), last_operation = 'withdraw' where id = w_id;
    return;
end;
$$ language plpgsql;
-- +goose statementend