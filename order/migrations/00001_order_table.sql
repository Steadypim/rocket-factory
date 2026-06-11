-- +goose Up
CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    part_ids UUID[] NOT NULL,
    payment_method TEXT,
    status TEXT NOT NULL,
    total_price NUMERIC(12, 2) NOT NULL,
    transaction_id UUID,

    CONSTRAINT orders_status_check
        CHECK (status IN ('PENDING_PAYMENT', 'PAID', 'CANCELLED')),

    CONSTRAINT orders_part_ids_check
        CHECK (cardinality(part_ids) > 0)
);

CREATE UNIQUE INDEX IF NOT EXISTS orders_transaction_id_idx
    ON orders (transaction_id)
    WHERE transaction_id IS NOT NULL;

-- +goose Down
DROP TABLE IF EXISTS orders;