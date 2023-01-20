package daog

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
	"time"
)

type DbConf struct {
	DbUrl    string
	Size     int
	Life     int
	IdleCons int
	IdleTime int
	LogSQL   bool
}

func NewDatasource(conf *DbConf) (Datasource, error) {
	dbUrl := conf.DbUrl
	if -1 == strings.Index(conf.DbUrl, "interpolateParams") {
		if strings.Index(conf.DbUrl, "?") != -1 {

			dbUrl = dbUrl + "&interpolateParams=true"
		} else {
			dbUrl = dbUrl + "?interpolateParams=true"
		}
	}
	db, err := sql.Open("mysql", dbUrl)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if conf.Size > 0 {
		db.SetMaxOpenConns(conf.Size)
	}
	if conf.IdleCons > 0 {
		db.SetMaxIdleConns(conf.IdleCons)
	}
	if conf.IdleTime > 0 {
		db.SetConnMaxIdleTime(time.Duration(int64(conf.IdleTime) * 1e9))
	}
	if conf.Life > 0 {
		db.SetConnMaxLifetime(time.Duration(int64(conf.Life) * 1e9))
	}

	return &database{db, conf.LogSQL}, nil
}

type Datasource interface {
	getDB(ctx context.Context) *sql.DB
	Shutdown()
	IsLogSQL() bool
}

type database struct {
	db     *sql.DB
	logSQL bool
}

func (db *database) getDB(ctx context.Context) *sql.DB {
	return db.db
}
func (db *database) Shutdown() {
	db.db.Close()
}

func (db *database) IsLogSQL() bool {
	return db.logSQL
}
