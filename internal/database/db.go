package database

import (
	"database/sql"
	"log"
	"os"
)

var (
	db                        *sql.DB
	preparedTokenExistsQuery  *sql.Stmt
	preparedCurrentTokenQuery *sql.Stmt
	preparedGetUserQuery      *sql.Stmt
	preparedCreateUserQuery   *sql.Stmt
	preparedUpdateTokenQuery  *sql.Stmt
	preparedActivateUserQuery *sql.Stmt
	preparedGetUserEmailQuery *sql.Stmt
	preparedGetUserProfile    *sql.Stmt
)

func Connect() {

	var err error
	db, err = sql.Open("postgres", os.Getenv("DB_URL"))

	if err != nil {
		log.Fatal("Cannot connect to the database: ", err)
	}

	log.Println("Connected to the database")

	err = initPreparedStatements()
	if err != nil {
		log.Fatal("Error initializing prepared statements: ", err)
	}

}

func GetDB() *sql.DB {
	return db
}

func initPreparedStatements() error {
	var err error

	log.Println("Preparing statements...")

	if preparedTokenExistsQuery, err = db.Prepare("SELECT COUNT(*) FROM user_tokens WHERE verification_token=$1 AND revoked_at IS NULL"); err != nil {
		log.Println("Error preparing token exists query: ", err)
		return err
	}

	if preparedCurrentTokenQuery, err = db.Prepare("UPDATE user_tokens set revoked_at = now() WHERE user_id=$1 AND token_type=$2 AND revoked_at IS NULL RETURNING verification_token"); err != nil {
		log.Println("Error preparing current token query: ", err)
		return err
	}

	if preparedGetUserQuery, err = db.Prepare("SELECT id, created_at, email, updated_at FROM users WHERE email=$1"); err != nil {
		log.Println("Error preparing get user query: ", err)
		return err
	}

	if preparedCreateUserQuery, err = db.Prepare("INSERT INTO users (email, hashed_password) VALUES ($1, $2) RETURNING id, email"); err != nil {
		log.Println("Error preparing create user query: ", err)
		return err
	}

	if preparedUpdateTokenQuery, err = db.Prepare("UPDATE user_tokens SET revoked_at=$1, user_id=$2 WHERE verification_token=$3"); err != nil {
		log.Println("Error preparing update token query: ", err)
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

	log.Println("All statements prepared successfully")
	return nil
}

func Close() {
	log.Println("Closing database connections and prepared statements...")

	if preparedActivateUserQuery != nil {
		preparedActivateUserQuery.Close()
	}

	if preparedUpdateTokenQuery != nil {
		preparedUpdateTokenQuery.Close()
	}

	if preparedCreateUserQuery != nil {
		preparedCreateUserQuery.Close()
	}

	if preparedGetUserQuery != nil {
		preparedGetUserQuery.Close()
	}

	if preparedCurrentTokenQuery != nil {
		preparedCurrentTokenQuery.Close()
	}

	if preparedTokenExistsQuery != nil {
		preparedTokenExistsQuery.Close()
	}

	if preparedGetUserEmailQuery != nil {
		preparedGetUserEmailQuery.Close()
	}

	if preparedGetUserProfile != nil {
		preparedGetUserProfile.Close()
	}

	if db != nil {
		db.Close()
	}
}

func GetPreparedTokenExistsQuery() *sql.Stmt {
	return preparedTokenExistsQuery
}

func GetPreparedCurrentTokenQuery() *sql.Stmt {
	return preparedCurrentTokenQuery
}

func GetPreparedGetUserQuery() *sql.Stmt {
	return preparedGetUserQuery
}

func GetPreparedCreateUserQuery() *sql.Stmt {
	return preparedCreateUserQuery
}

func GetPreparedUpdateTokenQuery() *sql.Stmt {
	return preparedUpdateTokenQuery
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
