package cybase

import (
	"strings"

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

	dbHost := readConfigFromSetting(setting, "options", "db_host")
	dbPort := readConfigFromSetting(setting, "options", "db_port")
	dbUser := readConfigFromSetting(setting, "options", "db_user")
	dbPassword := readConfigFromSetting(setting, "options", "db_password")
	dbName := readConfigFromSetting(setting, "options", "db_name")
	dbQuery := readConfigFromSetting(setting, "options", "db_query", map[string]string{"canEmpty": "true"})

	dbUrl := "postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "/" + dbName
	if dbQuery != "" {
		dbUrl += "?" + dbQuery
	}
	writeOsm := cydb.Osm("postgres", dbUrl, infoZapLogger, errorZapLogger)

	_, err = writeOsm.SelectKVS("SELECT key, value FROM cy_base_matedata;")(&matedata)
	if err != nil {
		panic(err)
	}

	readonlyOsm := cydb.Osm("postgres", ReadConfig("cy_readonly_db_url", dbUrl), infoZapLogger, errorZapLogger)

	return writeOsm, readonlyOsm, infoZapLogger, errorZapLogger
}

func readConfigFromSetting(setting *config.Config, section, option string, others ...map[string]string) string {
	value, err := setting.String(section, option)
	if err != nil {
		panic(err)
	}
	canEmpty := false
	if len(others) > 0 {
		other := others[0]
		if len(other) > 0 {
			canEmpty = other["canEmpty"] == "true"
		}
	}
	if value == "" && !canEmpty {
		panic(option + " is empty")
	}

	if strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'") {
		value = value[1 : len(value)-1]
	} else if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
		value = value[1 : len(value)-1]
	}

	return value
}
