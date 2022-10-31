-- +goose Up
create table wallets
(
    id          serial primary key,
    user_id     bigint      not null unique,
    balance     bigint               default 0 check ( balance >= 0 ), -- все, что касается денег – в копейках
    reservation bigint               default 0 check ( reservation <= balance ),
    created_at  timestamptz not null default now(),
    updated_at  timestamptz not null default now()                     -- Навешен триггер. На любое изменение ставить now().
);

create table transactions
(
    id         serial primary key,
    wallet_id  integer     not null references wallets (id),
    -- Возможные значения: reservation / write_off / cancel / incoming_transfer
    "type"     text        not null,
    -- Для reservation: { "order_id": <order_id> }
    -- Для write_off: { "order_id": <order_id> }
    -- Для cancel: { "order_id": <order_id> }
    -- Для incoming_transfer:
    --  - { "type": "enrollment" }
    payload    jsonb,
    amount     bigint check ( amount > 0 ),
    created_at timestamptz not null default now()
);

create table orders
(
    id          serial primary key,
    wallet_id   integer     not null references wallets (id),
    external_id bigint      not null unique,
    service_id  integer     not null,
    status      text,
    amount      bigint,
    created_at  timestamptz not null default now(),
    updated_at  timestamptz not null default now()
);

create table order_transactions
(
    id         serial primary key,
    order_id   integer     not null references orders (id),
    -- Возможные значения: reservation / write_off / cancel
    "type"     text        not null,
    created_at timestamptz not null default now()
);

create table report
(
    id            serial primary key,
    "period"      text    not null,
    service_id    integer not null,
    total_revenue bigint default 0 check ( total_revenue >= 0 )
);

create index wallets_user_idx on wallets (user_id);
create index transactions_created_idx on transactions (created_at desc);
create index report_period_idx on report (period desc);

-- +goose StatementBegin
create function trigger_set_timestamp()
    returns trigger as $$
begin
    new.updated_at=now();
    return new;
end;
$$ language 'plpgsql';
-- +goose StatementEnd

create trigger wallets_set_timestamp
    before update
    on wallets
    for each row
    execute procedure trigger_set_timestamp();

create trigger orders_set_timestamp
    before update
    on orders
    for each row
    execute procedure trigger_set_timestamp();

-- +goose Down
drop function trigger_set_timestamp;
drop trigger wallets_set_timestamp;
drop trigger orders_set_timestamp;
drop index transactions_created_idx;
drop index wallets_user_idx;
drop index report_period_idx;

drop table report;
drop table order_transactions;
drop table transactions;
drop table orders;
drop table wallets;