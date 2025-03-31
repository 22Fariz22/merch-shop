package migrate

import (
	"fmt"

	"github.com/22Fariz22/merch-shop/internal/models"
	"github.com/22Fariz22/merch-shop/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Migrate applies database migrations
func Migrate(logger logger.Logger, dsn string) error {
	// Инициализация GORM с использованием только для миграций
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		logger.Debugf("Error in pkg/db/migrate/migrate.go")
		return err
	}

	// Выполнение миграций в правильном порядке
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		logger.Debugf("Error in AutoMigrate User: %v", err)
		return err
	}

	err = db.AutoMigrate(&models.Wallet{})
	if err != nil {
		logger.Debugf("Error in AutoMigrate Wallet: %v", err)
		return err
	}

	err = db.AutoMigrate(&models.Purchase{})
	if err != nil {
		logger.Debugf("Error in AutoMigrate Purchase: %v", err)
		return err
	}

	err = db.AutoMigrate(&models.Transfer{})
	if err != nil {
		logger.Debugf("Error in AutoMigrate Transfer: %v", err)
		return err
	}

	// Проверка существования таблиц
	var usersTableExists, walletsTableExists, purchasesTableExists, transfersTableExists bool
	err = db.Raw("SELECT EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'users')").
		Scan(&usersTableExists).Error
	if err != nil {
		logger.Debugf("Error checking 'users' table existence: %v", err)
		return err
	}

	err = db.Raw("SELECT EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'wallets')").
		Scan(&walletsTableExists).Error
	if err != nil {
		logger.Debugf("Error checking 'wallets' table existence: %v", err)
		return err
	}

	err = db.Raw("SELECT EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'purchases')").
		Scan(&purchasesTableExists).Error
	if err != nil {
		logger.Debugf("Error checking 'purchases' table existence: %v", err)
		return err
	}

	err = db.Raw("SELECT EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'transfers')").
		Scan(&transfersTableExists).Error
	if err != nil {
		logger.Debugf("Error checking 'transfers' table existence: %v", err)
		return err
	}

	if usersTableExists {
		err = db.Exec(`ALTER TABLE users ALTER COLUMN user_id SET DEFAULT gen_random_uuid();`).Error
		if err != nil {
			logger.Debugf("Error setting DEFAULT for user_id: %v", err)
			return err
		}
	} else {
		logger.Debugf("Table 'users' does not exist after AutoMigrate")
		return fmt.Errorf("table 'users' was not created")
	}

	if walletsTableExists {
		err = db.Exec(`ALTER TABLE wallets ALTER COLUMN wallet_id SET DEFAULT gen_random_uuid();`).Error
		if err != nil {
			logger.Debugf("Error setting DEFAULT for wallet_id: %v", err)
			return err
		}
	} else {
		logger.Debugf("Table 'wallets' does not exist after AutoMigrate")
		return fmt.Errorf("table 'wallets' was not created")
	}

	if purchasesTableExists {
		err = db.Exec(`ALTER TABLE purchases ALTER COLUMN purchase_id SET DEFAULT gen_random_uuid();`).Error
		if err != nil {
			logger.Debugf("Error setting DEFAULT for purchase_id: %v", err)
			return err
		}
	} else {
		logger.Debugf("Table 'purchases' does not exist after AutoMigrate")
		return fmt.Errorf("table 'purchases' was not created")
	}

	if transfersTableExists {
		err = db.Exec(`ALTER TABLE transfers ALTER COLUMN transfer_id SET DEFAULT gen_random_uuid();`).Error
		if err != nil {
			logger.Debugf("Error setting DEFAULT for transfer_id: %v", err)
			return err
		}
	} else {
		logger.Debugf("Table 'transfers' does not exist after AutoMigrate")
		return fmt.Errorf("table 'transfers' was not created")
	}

	logger.Debugf("Database migrations applied successfully")
	return nil
}
