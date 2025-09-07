-- +goose Up
CREATE TYPE order_status AS ENUM (
    'PENDING_PAYMENT',
    'PAID',
    'CANCELLED'
);

CREATE TYPE payment_method AS ENUM (
    'PAYMENT_METHOD_UNSPECIFIED',
    'PAYMENT_METHOD_CARD',
    'PAYMENT_METHOD_SBP',
    'PAYMENT_METHOD_CREDIT_CARD',
    'PAYMENT_METHOD_INVESTOR_MONEY'
);

CREATE TABLE orders (
    order_uuid UUID PRIMARY KEY,
    user_uuid UUID NOT NULL,
    part_uuids JSONB NOT NULL DEFAULT '[]'::jsonb,
    total_price REAL NOT NULL,
    transaction_uuid UUID,
    payment_method payment_method NOT NULL DEFAULT 'PAYMENT_METHOD_UNSPECIFIED',
    status order_status NOT NULL DEFAULT 'PENDING_PAYMENT',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    cancelled_at TIMESTAMPTZ
);

COMMENT ON TABLE orders IS 'Таблица заказов';
COMMENT ON COLUMN orders.order_uuid IS 'Уникальный идентификатор заказа';
COMMENT ON COLUMN orders.user_uuid IS 'Уникальный идентификатор пользователя';
COMMENT ON COLUMN orders.part_uuids IS 'Массив UUID деталей в формате JSON';
COMMENT ON COLUMN orders.total_price IS 'Общая стоимость заказа';
COMMENT ON COLUMN orders.transaction_uuid IS 'Идентификатор транзакции оплаты (заполняется при оплате)';
COMMENT ON COLUMN orders.payment_method IS 'Способ оплаты';
COMMENT ON COLUMN orders.status IS 'Статус заказа';
COMMENT ON COLUMN orders.created_at IS 'Дата и время создания заказа';
COMMENT ON COLUMN orders.updated_at IS 'Дата и время последнего обновления заказа';
COMMENT ON COLUMN orders.cancelled_at IS 'Дата и время отмены заказа (заполняется при отмене)';

CREATE INDEX idx_orders_user_uuid ON orders(user_uuid);

-- +goose Down
DROP INDEX IF EXISTS idx_orders_user_uuid;

DROP TABLE IF EXISTS orders;

DROP TYPE IF EXISTS payment_method;
DROP TYPE IF EXISTS order_status;
