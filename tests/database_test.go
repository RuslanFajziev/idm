package database

import (
	"idm/inner/common"
	"idm/inner/database"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type DataMigration struct {
	Id      int64     `db:"id"`
	Version int64     `db:"version_id"`
	Applied bool      `db:"is_applied"`
	Create  time.Time `db:"tstamp"`
}

func ClearEnv() {
	driverName := "DB_DRIVER_NAME"
	dsnName := "DB_DSN"
	os.Unsetenv(driverName)
	os.Unsetenv(dsnName)
}

// 1 - в проекте нет .env  файла (должны получить конфигурацию из пременных окружения)
func TestConnectionDbСase1(t *testing.T) {
	ClearEnv()
	as := assert.New(t)
	driverName := "DB_DRIVER_NAME"
	driverVal := "postgres"
	dsnName := "DB_DSN"
	dsnVal := "host=127.0.0.1 port=5430 user=postgres_user password=postgres_password dbname=idm_go sslmode=disable"

	conf, err := common.GetConfig(".env_not_found")
	if err != nil {
		os.Setenv(driverName, driverVal)
		os.Setenv(dsnName, dsnVal)

		conf = common.Config{
			DbDriverName: os.Getenv(driverName),
			Dsn:          os.Getenv(dsnName),
		}
	}

	as.Equal(driverVal, conf.DbDriverName)
	as.Equal(dsnVal, conf.Dsn)
	as.NotEqual(nil, err)
}

// 2 - в проекте есть .env  файл, но в нём нет нужных переменных и в переменных окружения их тоже нет (должны получить пустую структуру idm.inner.common.Config)
func TestConnectionDbСase2(t *testing.T) {
	ClearEnv()
	as := assert.New(t)
	conf, err := common.GetConfig(".env_2")

	as.Equal("", conf.DbDriverName)
	as.Equal("", conf.Dsn)
	as.NotEqual(nil, err)
}

// 3 - в проекте есть .env  файл и в нём нет нужных переменных, но в переменных окружения они есть (должны получить заполненную структуру  idm.inner.common.Config с данными из пременных окружения)
func TestConnectionDbСase3(t *testing.T) {
	ClearEnv()
	as := assert.New(t)
	driverName := "DB_DRIVER_NAME"
	driverVal := "postgres"
	dsnName := "DB_DSN"
	dsnVal := "host=127.0.0.1 port=5430 user=postgres_user password=postgres_password dbname=idm_go sslmode=disable"

	os.Setenv(driverName, driverVal)
	os.Setenv(dsnName, dsnVal)

	conf, err := common.GetConfig(".env_2")

	as.Equal(driverVal, conf.DbDriverName)
	as.Equal(dsnVal, conf.Dsn)
	as.Equal(nil, err)
}

// 4 - в проекте есть корректно заполненный .env файл, в переменных окружения нет конфликтующих с ним переменных  (должны получить структуру  idm.inner.common.Config, заполненную данными из .env файла)
func TestConnectionDbСase4(t *testing.T) {
	ClearEnv()
	as := assert.New(t)
	driverName := "DB_DRIVER_NAME"
	driverVal := "postgres"
	dsnVal := "host=127.0.0.1 port=5430 user=postgres_user password=postgres_password dbname=idm_go sslmode=disable"

	os.Setenv(driverName, driverVal)

	conf, err := common.GetConfig(".env_4")

	as.Equal(driverVal, conf.DbDriverName)
	as.Equal(dsnVal, conf.Dsn)
	as.Equal(nil, err)
}

// 5 - в проекте есть .env  файл и в нём есть нужные переменные, но в переменных окружения они тоже есть (с другими значениями) - должны получить структуру  idm.inner.common.Config, заполненную данными. Нужно проверить, какими значениями она будет заполнена (из .env файла или из переменных окружения)
func TestConnectionDbСase5(t *testing.T) {
	ClearEnv()
	as := assert.New(t)
	driverName := "DB_DRIVER_NAME"
	driverVal := "mssql"
	dsnName := "DB_DSN"
	dsnVal := "host=localhost port=5445"

	os.Setenv(driverName, driverVal)
	os.Setenv(dsnName, dsnVal)

	conf, err := common.GetConfig(".env_2")

	as.Equal(driverVal, conf.DbDriverName)
	as.Equal(dsnVal, conf.Dsn)
	as.Equal(nil, err)

	conf, err = common.GetConfig(".env_5")
	as.NotEqual("", conf.DbDriverName)
	as.NotEqual("", conf.Dsn)
	as.Equal(driverVal, conf.DbDriverName)
	as.Equal(dsnVal, conf.Dsn)
	as.Equal(nil, err)

	ClearEnv()
	conf, err = common.GetConfig(".env_5")
	as.NotEqual("", conf.DbDriverName)
	as.NotEqual("", conf.Dsn)
	as.NotEqual(driverVal, conf.DbDriverName)
	as.NotEqual(dsnVal, conf.Dsn)
	as.Equal(nil, err)
}

// 6 - приложение не может подключиться к базе данных с некорректным конфигом (например, неправильно указан: хост, порт, имя базы данных, логин или пароль)
func TestConnectionDbСase6(t *testing.T) {
	ClearEnv()
	as := assert.New(t)
	_, err := database.ConnectDb(".env_6")
	as.NotEqual(nil, err)
}

// 7 - приложение может подключиться к базе данных с корректным конфигом
func TestConnectionDbСase7(t *testing.T) {
	ClearEnv()
	as := assert.New(t)
	db, err := database.ConnectDb(".env_5")

	if err == nil {
		defer db.Close()
	}

	as.Equal(nil, err)
	var dataMigration DataMigration
	var verId int64 = 20250603224245
	err = db.Get(&dataMigration, "SELECT * FROM goose_db_version WHERE version_id = $1", verId)
	as.Equal(nil, err)
	as.Equal(dataMigration.Version, verId)
	as.Equal(dataMigration.Applied, true)
}
