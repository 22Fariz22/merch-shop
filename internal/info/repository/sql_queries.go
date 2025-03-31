package repository

const (
	getWalletQuery = `
        SELECT wallet_id, user_id, balance, updated_at 
        FROM wallets 
        WHERE user_id = $1
    `
	getTransfers = `
        SELECT 
            t.transfer_id,
            t.wallet_from_id,
            t.wallet_to_id,
            t.amount,
            t.created_at,
            t.status,
            uf.username AS from_username,
            ut.username AS to_username
        FROM transfers t
        LEFT JOIN wallets wf ON wf.wallet_id = t.wallet_from_id
        LEFT JOIN wallets wt ON wt.wallet_id = t.wallet_to_id
        LEFT JOIN users uf ON uf.user_id = wf.user_id
        LEFT JOIN users ut ON ut.user_id = wt.user_id
        WHERE t.wallet_from_id IN (SELECT wallet_id FROM wallets WHERE user_id = $1)
           OR t.wallet_to_id IN (SELECT wallet_id FROM wallets WHERE user_id = $1)
    `
	getPurchases = `
        SELECT 
            item AS "type",
            COUNT(*) AS quantity
        FROM purchases 
        WHERE user_id = $1
        GROUP BY item
    `
)
