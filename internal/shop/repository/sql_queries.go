package repository

const (
	getWalletQuery = `SELECT wallet_id, user_id, balance, updated_at FROM wallets WHERE user_id = $1 FOR UPDATE`

	updateWalletQuery = `UPDATE wallets SET balance = $1, updated_at = NOW() WHERE wallet_id = $2 RETURNING wallet_id, user_id, balance, updated_at`

	createPurchaseQuery = `
        INSERT INTO purchases (user_id, wallet_id, item, price, created_at)
        VALUES ($1, $2, $3, $4, NOW())
        RETURNING purchase_id, user_id, wallet_id, item, price, created_at
    `
)
