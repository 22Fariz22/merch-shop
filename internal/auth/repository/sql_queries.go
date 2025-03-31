package repository

const (
	createUserQuery = `INSERT INTO users (username, password, created_at)
						VALUES ($1, $2, now()) 
						RETURNING *`

	createWalletQuery = `
        INSERT INTO wallets (user_id, balance, updated_at)
        VALUES ($1, $2,  NOW())
        RETURNING  user_id, balance, updated_at
    `

	deleteUserQuery = `DELETE FROM users WHERE user_id = $1`

	getUserQuery = `SELECT user_id, username 
					 FROM users 
					 WHERE user_id = $1`

	findUserByUsername = `SELECT user_id, username,password, created_at
				 		FROM users 
				 		WHERE username = $1`
)
