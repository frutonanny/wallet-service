-- +goose Up
-- Таблица wallets содержит информцию по каждому пользователю.
create table wallets
(
    id          serial primary key,
    user_id     bigint      not null unique,
    balance     bigint               default 0 check ( balance >= 0 ), -- все, что касается денег – в копейках
    reservation bigint               default 0 check ( reservation >= 0),
    created_at  timestamptz not null default now(),
    updated_at  timestamptz not null default now()                     -- Навешен триггер. На любое изменение ставить now().
);

create index wallets_user_idx on wallets (user_id);

-- В таблицу transactions заносятся абсолютно все денежные операции.
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

create index transactions_created_idx on transactions (created_at desc);

-- В таблицу orders заносится информация по полученным внешним заказам для покупки услуг.
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

create index orders_external_idx on orders (external_id);

-- В таблицу order_transactions заносятся все операции, связанные с заказами.
create table order_transactions
(
    id         serial primary key,
    order_id   integer     not null references orders (id),
    -- Возможные значения: reservation / write_off / cancel
    "type"     text        not null,
    created_at timestamptz not null default now()
);

-- В таблицу report заносится информация о купленных услугах за месяц.
create table report
(
    id            serial primary key,
    "period"      text    not null,
    service_id    integer not null,
    total_revenue bigint default 0 check ( total_revenue >= 0 )
);

create index report_period_idx on report (period desc);

-- Таблица содержит информацию о существующих услугах
create table services
(
    id     serial primary key,
    "name" text not null
);

-- Заполняем таблицу services видами услуг.
insert into services("name")
values ('more views'),
       ('highlighting'),
       ('XL-advert');

-- +goose StatementBegin
create function trigger_set_timestamp()
    returns trigger as $$
begin
    new.updated_at
=now();
return new;
end;
$$
language 'plpgsql';
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
drop trigger wallets_set_timestamp on wallets;
drop trigger orders_set_timestamp on orders;
drop function trigger_set_timestamp;
drop index transactions_created_idx;
drop index wallets_user_idx;
drop index orders_external_idx;
drop index report_period_idx;
drop table report;
drop table order_transactions;
drop table transactions;
drop table orders;
drop table services;
drop table wallets;