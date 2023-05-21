package cybase

import (
	"github.com/chuanmoon/utils/cydb"
	"github.com/yinshuwei/config"
	"github.com/yinshuwei/osm/v2"
	"go.uber.org/zap"

	// justifying
	_ "github.com/lib/pq"
)

var (
	infoZapLogger  *zap.Logger
	errorZapLogger *zap.Logger
	matedata       = map[string]string{}
)

// Init init cybase
// appName: the name of the application
// return writeOsm, readonlyOsm, infoZapLogger, errorZapLogger
func Init(appName string) (*osm.Osm, *osm.Osm, *zap.Logger, *zap.Logger) {
	infoZapLogger, errorZapLogger = ecszapInfoAndErrorLogger(appName)
	setting, err := config.ReadDefault("/etc/chuanmoon/odoo.conf")
	if err != nil {
		panic(err)
	}

	dbHost, err := setting.String("options", "db_host")
	if err != nil {
		panic(err)
	}
	if dbHost == "" {
		panic("db_host is empty")
	}

	dbPort, err := setting.String("options", "db_port")
	if err != nil {
		panic(err)
	}
	if dbPort == "" {
		panic("db_port is empty")
	}

	dbUser, err := setting.String("options", "db_user")
	if err != nil {
		panic(err)
	}
	if dbUser == "" {
		panic("db_user is empty")
	}

	dbPassword, err := setting.String("options", "db_password")
	if err != nil {
		panic(err)
	}
	if dbPassword == "" {
		panic("db_password is empty")
	}

	dbName, err := setting.String("options", "db_name")
	if err != nil {
		panic(err)
	}
	if dbName == "" {
		panic("db_name is empty")
	}

	dbQuery, err := setting.String("options", "db_query")
	if err != nil {
		panic(err)
	}

	dbUrl := "postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "/" + dbName
	if dbQuery != "" {
		dbUrl += "?" + dbQuery
	}
	writeOsm := cydb.Osm("postgres", dbUrl, infoZapLogger, errorZapLogger)

	_, err = writeOsm.SelectKVS("SELECT key, value FROM cy_matedata;")(&matedata)
	if err != nil {
		panic(err)
	}

	readonlyOsm := cydb.Osm("postgres", ReadConfig("cy_readonly_db_url", dbUrl), infoZapLogger, errorZapLogger)

	return writeOsm, readonlyOsm, infoZapLogger, errorZapLogger
}
