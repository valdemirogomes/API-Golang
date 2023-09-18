package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/melisource/fury_go-dev-base-3-v2/cmd/api/config"
	"github.com/melisource/fury_go-dev-base-3-v2/cmd/api/lib/helperdb"
	"github.com/mercadolibre/fury_go-core/pkg/telemetry"
	"github.com/mercadolibre/go-meli-toolkit/gomelipass"
	"github.com/mercadolibre/go-meli-toolkit/goutils/logger"
)

type DB interface {
	WithoutTransaction(ctx context.Context, exec func(helperdb.Tx) error) (err error)
	WithTransaction(ctx context.Context, txFunc func(*sql.Tx) error) (err error)
	Close() error
}

const maxConnLifetimeMinutes = 10

type MySQL struct {
	*sql.DB
}

func NewMySQL(config config.Environment) (DB, error) {
	var passwd string
	var host string

	if passwd = gomelipass.GetEnv(config.MySQLConfig.Password); passwd == "" {
		passwd = config.MySQLConfig.Password
	}

	if host = gomelipass.GetEnv(config.MySQLConfig.Host); host == "" {
		host = config.MySQLConfig.Host
	}

	connectionString := buildConnectionString(config.MySQLConfig.User, passwd, host, config.MySQLConfig.Database)

	db, err := sql.Open(config.MySQLConfig.Drive, connectionString)
	if err != nil {
		logger.Errorf("[event: fail_db_init][service: db_service] Could not start DB connection %s", err, err.Error())

		return nil, err
	}
	db.SetMaxOpenConns(config.MySQLConfig.PoolSizeMax)
	db.SetMaxIdleConns(config.MySQLConfig.PoolSizeIddle)
	db.SetConnMaxLifetime(maxConnLifetimeMinutes * time.Minute)

	return &MySQL{db}, nil
}

func buildConnectionString(user string, password string, host string, database string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true", user, password, host, database)
}

func (db MySQL) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return db.DB.QueryContext(ctx, query, args...)
}

func (db MySQL) WithoutTransaction(ctx context.Context, txFunc func(helperdb.Tx) error) error {
	ctx, span := telemetry.StartSpan(ctx, "mysql_without_transaction")
	defer span.Finish()

	return txFunc(db)
}

func (db MySQL) WithTransaction(ctx context.Context, txFunc func(*sql.Tx) error) (err error) {
	spanTransaction := "mysql_with_transaction"
	ctx, span := telemetry.StartSpan(ctx, spanTransaction)
	span.SetLabel(spanTransaction, "begin")

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			span.SetLabel(spanTransaction, "rollback")
			err = tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			logger.Debugf("Error preventing transaction commit %+v", err)
			span.SetLabel(spanTransaction, "rollback")
			_ = tx.Rollback()
		} else {
			span.SetLabel(spanTransaction, "commit")
			err = tx.Commit() // err is nil; if Commit returns error update err
		}
		span.Finish()
	}()

	err = txFunc(tx)
	return err
}
