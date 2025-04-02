-- +goose Up
CREATE TABLE transfers (
    transfer_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_from_id UUID NOT NULL,
    wallet_to_id UUID NOT NULL,
    amount INT NOT NULL CHECK (amount >= 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) DEFAULT 'completed' CHECK (status IN ('pending', 'completed', 'failed')),
    CONSTRAINT fk_wallet_from FOREIGN KEY (wallet_from_id) REFERENCES wallets(wallet_id) ON DELETE CASCADE,
    CONSTRAINT fk_wallet_to FOREIGN KEY (wallet_to_id) REFERENCES wallets(wallet_id) ON DELETE CASCADE
);

CREATE INDEX idx_transfers_wallet_from_id ON transfers(wallet_from_id);
CREATE INDEX idx_transfers_wallet_to_id ON transfers(wallet_to_id);

-- +goose Down
DROP TABLE transfers;