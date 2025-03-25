package repository

const (
	createUserQuery = `INSERT INTO users (username, password, created_at)
						VALUES ($1, $2, now()) 
						RETURNING *`


	deleteUserQuery = `DELETE FROM users WHERE user_id = $1`

	getUserQuery = `SELECT user_id, first_name, last_name, email, role, about, avatar, phone_number, 
       				 address, city, gender, postcode, birthday, created_at, updated_at, login_date  
					 FROM users 
					 WHERE user_id = $1`

	getTotal = `SELECT COUNT(user_id) FROM users`
	

	findUserByUsername = `SELECT user_id, username, created_at, password
				 		FROM users 
				 		WHERE username = $1`
)
