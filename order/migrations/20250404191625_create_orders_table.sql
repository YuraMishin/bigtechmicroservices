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
    cancelled_at TIMESTAMPTZ,
    
    CONSTRAINT check_positive_total_price CHECK (total_price > 0),
    CONSTRAINT check_part_uuids_format CHECK (
        jsonb_typeof(part_uuids) = 'array'
    ),
    CONSTRAINT check_cancelled_at_logic CHECK (
        (status = 'CANCELLED' AND cancelled_at IS NOT NULL) OR
        (status != 'CANCELLED' AND cancelled_at IS NULL)
    ),
    CONSTRAINT check_transaction_uuid_logic CHECK (
        (status = 'PAID' AND transaction_uuid IS NOT NULL) OR
        (status != 'PAID')
    )
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
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_created_at ON orders(created_at);
CREATE INDEX idx_orders_transaction_uuid ON orders(transaction_uuid) WHERE transaction_uuid IS NOT NULL;

CREATE INDEX idx_orders_part_uuids ON orders USING GIN (part_uuids);

CREATE OR REPLACE FUNCTION update_updated_at_column() RETURNS TRIGGER AS $func$ BEGIN NEW.updated_at = NOW(); RETURN NEW; END; $func$ LANGUAGE plpgsql;

CREATE TRIGGER update_orders_updated_at BEFORE UPDATE ON orders FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE OR REPLACE FUNCTION set_cancelled_at() RETURNS TRIGGER AS $func$ BEGIN IF NEW.status = 'CANCELLED' AND OLD.status != 'CANCELLED' THEN NEW.cancelled_at = NOW(); END IF; RETURN NEW; END; $func$ LANGUAGE plpgsql;

CREATE TRIGGER set_orders_cancelled_at BEFORE UPDATE ON orders FOR EACH ROW EXECUTE FUNCTION set_cancelled_at();

-- +goose Down
DROP TRIGGER IF EXISTS set_orders_cancelled_at ON orders;
DROP TRIGGER IF EXISTS update_orders_updated_at ON orders;

DROP FUNCTION IF EXISTS set_cancelled_at();
DROP FUNCTION IF EXISTS update_updated_at_column();

DROP INDEX IF EXISTS idx_orders_part_uuids;
DROP INDEX IF EXISTS idx_orders_transaction_uuid;
DROP INDEX IF EXISTS idx_orders_created_at;
DROP INDEX IF EXISTS idx_orders_status;
DROP INDEX IF EXISTS idx_orders_user_uuid;

DROP TABLE IF EXISTS orders;

DROP TYPE IF EXISTS payment_method;
DROP TYPE IF EXISTS order_status;
