package infrastructure

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func ConnectDatabase() *sqlx.DB {
	// Get config with viper
	dbHost := viper.GetString("database.host")
	dbUser := viper.GetString("database.user")
	dbPass := viper.GetString("database.pass")
	dbName := viper.GetString("database.name")
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true", dbUser, dbPass, dbHost, dbName)
	databaseConnection, err := sqlx.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}

	err = databaseConnection.Ping()
	if err != nil {
		panic(err)
	}
	logrus.Info("MySQL is connected")

	if viper.GetString("env") == "test" {
		if err = migrateUpTest(databaseConnection); err != nil {
			panic(err)
		}
	} else if viper.GetString("env") == "dev" {
		if err = migrateUpDevelop(databaseConnection); err != nil {
			panic(err)
		}
	}

	return databaseConnection
}

func migrateUpDevelop(db *sqlx.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS user (
			id int NOT NULL AUTO_INCREMENT,
			address varchar(48) NOT NULL,
			dob date NOT NULL,
			name varchar(256) NOT NULL,
			updated_at timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			created_at timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  			PRIMARY KEY (id),
  			INDEX(name ASC)
		)
		ENGINE=InnoDB
		DEFAULT CHARSET=utf8mb4
		COLLATE=utf8mb4_0900_ai_ci;
		`)

	return err
}

func migrateUpTest(db *sqlx.DB) error {

	if _, err := db.Exec(`DROP TABLE IF EXISTS user`); err != nil {
		return err
	}

	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS user (
			id int NOT NULL AUTO_INCREMENT,
			address varchar(256) NOT NULL,
			dob date NOT NULL,
			name varchar(48) NOT NULL,
			updated_at timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			created_at timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  			PRIMARY KEY (id),
  			INDEX(name ASC)
		)
		ENGINE=InnoDB
		DEFAULT CHARSET=utf8mb4
		COLLATE=utf8mb4_0900_ai_ci;
		`); err != nil {
		return err
	}

	bulkInsert(db)
	return nil
}

func bulkInsert(db *sqlx.DB) {
	if _, err := db.Exec("INSERT INTO user (name, address, dob) VALUES (?,?,?)", "Jhon", "5th avenue", "1991-01-10"); err != nil {
		panic(err)
	}

	if _, err := db.Exec("INSERT INTO user (name, address, dob) VALUES (?,?,?)", "Jhon", "5th avenue", "1991-01-10"); err != nil {
		panic(err)
	}

	if _, err := db.Exec("INSERT INTO user (name, address, dob) VALUES (?,?,?)", "Jhon", "5th avenue", "1991-01-10"); err != nil {
		panic(err)
	}

	if _, err := db.Exec("INSERT INTO user (name, address, dob) VALUES (?,?,?)", "Jhon", "5th avenue", "1991-01-10"); err != nil {
		panic(err)
	}

	if _, err := db.Exec("INSERT INTO user (name, address, dob) VALUES (?,?,?)", "Jhon", "5th avenue", "1991-01-10"); err != nil {
		panic(err)
	}

	if _, err := db.Exec("INSERT INTO user (name, address, dob) VALUES (?,?,?)", "Jhon", "5th avenue", "1991-01-10"); err != nil {
		panic(err)
	}

	if _, err := db.Exec("INSERT INTO user (name, address, dob) VALUES (?,?,?)", "Jhon", "5th avenue", "1991-01-10"); err != nil {
		panic(err)
	}

	if _, err := db.Exec("INSERT INTO user (name, address, dob) VALUES (?,?,?)", "Jhon", "5th avenue", "1991-01-10"); err != nil {
		panic(err)
	}
}
