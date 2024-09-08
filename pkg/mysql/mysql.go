package mysql

import (
	"database/sql"
	"fmt"
	"github.com/barancanatbas/messaging/config"
	"github.com/rs/zerolog/log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func NewMysqlClient(config config.DatabaseConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", config.User, config.Password, config.Host, config.Port, config.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Error().Err(err).Msg("Could not open DB connection")
		return nil, fmt.Errorf("could not open db connection: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Error().Err(err).Msg("Could not ping MySQL DB")
		return nil, fmt.Errorf("could not ping db: %w", err)
	}

	log.Info().Str("host", config.Host).Str("db", config.DBName).Msg("Successfully connected to MySQL database")

	if err := runMigration(db); err != nil {
		return nil, fmt.Errorf("failed to run migration: %w", err)
	}

	return db, nil
}

func runMigration(db *sql.DB) error {
	checkTableQuery := `
		SELECT COUNT(*) 
		FROM information_schema.tables 
		WHERE table_name = 'messages' AND table_schema = DATABASE();
	`

	var count int
	err := db.QueryRow(checkTableQuery).Scan(&count)
	if err != nil {
		log.Error().Err(err).Msg("Error checking if table exists")
		return err
	}

	if count == 0 {
		log.Info().Msg("Table 'messages' does not exist, creating...")

		createTableQuery := `
			CREATE TABLE messages (
				id INT AUTO_INCREMENT PRIMARY KEY,
				content TEXT NOT NULL,
				phone_number VARCHAR(15) NOT NULL,
				status ENUM('PENDING', 'FAILED', 'DELIVERED') NOT NULL DEFAULT 'PENDING',
				sent_at DATETIME NULL,
				uuid VARCHAR(36)
			);
		`

		_, err := db.Exec(createTableQuery)
		if err != nil {
			log.Error().Err(err).Msg("Error creating 'messages' table")
			return err
		}

		log.Info().Msg("Table 'messages' created successfully")

		insertDummyDataQuery := `
			INSERT INTO messages (content, phone_number, status) VALUES 
			('Hello, World!', '+1234567890', 'PENDING'),
			('Test Message', '+0987654321', 'PENDING');
		`

		_, err = db.Exec(insertDummyDataQuery)
		if err != nil {
			log.Error().Err(err).Msg("Error inserting dummy data into 'messages' table")
			return err
		}

		log.Info().Msg("Dummy data inserted successfully")
	} else {
		log.Info().Msg("Table 'messages' already exists")
	}

	return nil
}
