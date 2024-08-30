package database

import (
	"database/sql"
	"log"
)

var (
	preparedTokenExistsQuery  *sql.Stmt
	preparedCurrentTokenQuery *sql.Stmt
	preparedUpdateTokenQuery  *sql.Stmt
)

func initPreparedTokenStatements() error {

	var err error

	log.Println("Preparing token statements")

	if preparedTokenExistsQuery, err = db.Prepare("SELECT COUNT(*) FROM user_tokens WHERE verification_token=$1 AND revoked_at IS NULL"); err != nil {
		log.Println("Error preparing token exists query: ", err)
		return err
	}

	if preparedCurrentTokenQuery, err = db.Prepare("UPDATE user_tokens set revoked_at = now() WHERE user_id=$1 AND token_type=$2 AND revoked_at IS NULL RETURNING verification_token"); err != nil {
		log.Println("Error preparing current token query: ", err)
		return err
	}

	if preparedUpdateTokenQuery, err = db.Prepare("UPDATE user_tokens SET revoked_at=$1, user_id=$2 WHERE verification_token=$3"); err != nil {
		log.Println("Error preparing update token query: ", err)
		return err
	}

	log.Println("All token statements prepared successfully")

	return nil
}

func closePreparedTokenStatements() {

	log.Println("Closing token prepared statements...")

	if preparedTokenExistsQuery != nil {
		preparedTokenExistsQuery.Close()
	}
	if preparedCurrentTokenQuery != nil {
		preparedCurrentTokenQuery.Close()
	}
	if preparedUpdateTokenQuery != nil {
		preparedUpdateTokenQuery.Close()
	}

	log.Println("Token prepared statements closed")
}

func GetPreparedTokenExistsQuery() *sql.Stmt {
	return preparedTokenExistsQuery
}

func GetPreparedCurrentTokenQuery() *sql.Stmt {
	return preparedCurrentTokenQuery
}

func GetPreparedUpdateTokenQuery() *sql.Stmt {
	return preparedUpdateTokenQuery
}
