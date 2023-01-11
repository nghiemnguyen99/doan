package infrastructure

import (
	"database/sql"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	mysqlDriver "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func GetGormConfig() *gorm.Config {
	return &gorm.Config{
		DisableAutomaticPing: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}
}

func NewDatabase(logger *logrus.Logger) (db *gorm.DB, master *sql.DB, err error) {
	gormConfig := GetGormConfig()
	gormConfig.Logger = NewGormLogger(logger)

	connMaster, err := sql.Open("mysql", os.Getenv("MYSQL_MASTER_CONNECTION_STRING"))
	if err != nil {
		return nil, nil, err
	}

	connMaster.SetConnMaxLifetime(time.Minute * 5)
	connMaster.SetConnMaxIdleTime(time.Minute * 20)
	connMaster.SetMaxIdleConns(50)
	connMaster.SetMaxOpenConns(100)

	driverMaster := mysqlDriver.New(mysqlDriver.Config{
		Conn:                      connMaster,
		SkipInitializeWithVersion: true,
		DontSupportForShareClause: true,
	})

	db, err = gorm.Open(driverMaster, gormConfig)
	if err != nil {
		// The connection string contains the password, so having
		// it in the error could risk leaking the password
		return nil, nil, err
	}
	err = Ping(db)
	if err != nil {
		return nil, nil, err
	}

	return db, connMaster, nil
}

func Ping(db *gorm.DB) error {
	mysqlDB, err := db.DB()
	if err != nil {
		return err
	}

	return mysqlDB.Ping()
}
