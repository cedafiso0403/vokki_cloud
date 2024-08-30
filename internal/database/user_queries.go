package database

import (
	"database/sql"
	"log"
)

var (
	preparedGetUserQuery      *sql.Stmt
	preparedCreateUserQuery   *sql.Stmt
	preparedActivateUserQuery *sql.Stmt
	preparedGetUserEmailQuery *sql.Stmt
	preparedGetUserProfile    *sql.Stmt
	preparedUpdateUserProfile *sql.Stmt
)

func initPreparedUserStaments() error {
	var err error

	log.Println("Preparing user statements")

	if preparedGetUserQuery, err = db.Prepare("SELECT id, created_at, email, updated_at FROM users WHERE email=$1"); err != nil {
		log.Println("Error preparing get user query: ", err)
		return err
	}

	if preparedCreateUserQuery, err = db.Prepare("INSERT INTO users (email, hashed_password) VALUES ($1, $2) RETURNING id, email"); err != nil {
		log.Println("Error preparing create user query: ", err)
		return err
	}

	if preparedActivateUserQuery, err = db.Prepare("UPDATE users SET activated=true WHERE id=$1"); err != nil {
		log.Println("Error preparing activate user query: ", err)
		return err
	}

	if preparedGetUserEmailQuery, err = db.Prepare("SELECT id, email FROM users WHERE email=$1"); err != nil {
		log.Println("Error preparing get user email query: ", err)
		return err
	}

	if preparedGetUserProfile, err = db.Prepare("SELECT users.id, users.email, user_profiles.first_name, user_profiles.last_name FROM users INNER JOIN user_profiles ON users.id = user_profiles.user_id WHERE users.id=$1"); err != nil {
		log.Println("Error preparing get user profile query: ", err)
		return err
	}

	if preparedUpdateUserProfile, err = db.Prepare("WITH updated AS (UPDATE user_profiles SET first_name = COALESCE($1, first_name), last_name = COALESCE($2, last_name) WHERE user_id = $3 RETURNING user_id, first_name, last_name) SELECT users.id, users.email, updated.first_name, updated.last_name FROM users INNER JOIN updated ON users.id = updated.user_id WHERE users.id = $3"); err != nil {
		log.Println("Error preparing update user profile query: ", err)
		return err
	}

	log.Println("All user statements prepared successfully")

	return nil
}

func closePreparedUserStatements() {
	log.Println("Closing user prepared statements...")

	if preparedGetUserQuery != nil {
		preparedGetUserQuery.Close()
	}

	if preparedCreateUserQuery != nil {
		preparedCreateUserQuery.Close()
	}

	if preparedActivateUserQuery != nil {
		preparedActivateUserQuery.Close()
	}

	if preparedGetUserEmailQuery != nil {
		preparedGetUserEmailQuery.Close()
	}

	if preparedGetUserProfile != nil {
		preparedGetUserProfile.Close()
	}

	if preparedUpdateUserProfile != nil {
		preparedUpdateUserProfile.Close()
	}

	log.Println("User prepared statements closed")
}

func GetPreparedGetUserQuery() *sql.Stmt {
	return preparedGetUserQuery
}

func GetPreparedCreateUserQuery() *sql.Stmt {
	return preparedCreateUserQuery
}

func GetPreparedActivateUserQuery() *sql.Stmt {
	return preparedActivateUserQuery
}

func GetPreparedGetUserEmailQuery() *sql.Stmt {
	return preparedGetUserEmailQuery
}

func GetPreparedGetUserProfile() *sql.Stmt {
	return preparedGetUserProfile
}

func GetPreparedUpdateUserProfile() *sql.Stmt {
	return preparedUpdateUserProfile
}
