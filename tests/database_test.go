package database

import (
	"idm/inner/common"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 1 - в проекте нет .env  файла (должны получить конфигурацию из пременных окружения)
func TestConnectionDbСase1(t *testing.T) {
	as := assert.New(t)
	driverName := "DB_DRIVER_NAME"
	driverVal := "postgres"
	dsnName := "DB_DSN"
	dsnVal := "host=127.0.0.1 port=5430 user=postgres_user password=postgres_password dbname=idm_go sslmode=disable"
	os.Unsetenv(driverName)
	os.Unsetenv(dsnName)

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
	as := assert.New(t)
	driverName := "DB_DRIVER_NAME"
	dsnName := "DB_DSN"
	os.Unsetenv(driverName)
	os.Unsetenv(dsnName)
	conf, err := common.GetConfig(".env_2")

	as.Equal("", conf.DbDriverName)
	as.Equal("", conf.Dsn)
	as.NotEqual(nil, err)
}

// 3 - в проекте есть .env  файл и в нём нет нужных переменных, но в переменных окружения они есть (должны получить заполненную структуру  idm.inner.common.Config с данными из пременных окружения)
func TestConnectionDbСase3(t *testing.T) {
	as := assert.New(t)
	driverName := "DB_DRIVER_NAME"
	driverVal := "postgres"
	dsnName := "DB_DSN"
	dsnVal := "host=127.0.0.1 port=5430 user=postgres_user password=postgres_password dbname=idm_go sslmode=disable"
	os.Unsetenv(driverName)
	os.Unsetenv(dsnName)

	os.Setenv(driverName, driverVal)
	os.Setenv(dsnName, dsnVal)

	conf, err := common.GetConfig(".env_2")

	as.Equal(driverVal, conf.DbDriverName)
	as.Equal(dsnVal, conf.Dsn)
	as.Equal(nil, err)
}

// 4 - в проекте есть корректно заполненный .env файл, в переменных окружения нет конфликтующих с ним переменных  (должны получить структуру  idm.inner.common.Config, заполненную данными из .env файла)
func TestConnectionDbСase4(t *testing.T) {
	as := assert.New(t)
	driverName := "DB_DRIVER_NAME"
	driverVal := "postgres"
	dsnVal := "host=127.0.0.1 port=5430 user=postgres_user password=postgres_password dbname=idm_go sslmode=disable"
	dsnName := "DB_DSN"
	os.Unsetenv(driverName)
	os.Unsetenv(dsnName)

	os.Setenv(driverName, driverVal)

	conf, err := common.GetConfig(".env_4")

	as.Equal(driverVal, conf.DbDriverName)
	as.Equal(dsnVal, conf.Dsn)
	as.Equal(nil, err)
}

// 5 - в проекте есть .env  файл и в нём есть нужные переменные, но в переменных окружения они тоже есть (с другими значениями) - должны получить структуру  idm.inner.common.Config, заполненную данными. Нужно проверить, какими значениями она будет заполнена (из .env файла или из переменных окружения)
func TestConnectionDbСase5(t *testing.T) {
	as := assert.New(t)
	driverName := "DB_DRIVER_NAME"
	driverVal := "mssql"
	dsnName := "DB_DSN"
	dsnVal := "host=localhost port=5445"
	os.Unsetenv(driverName)
	os.Unsetenv(dsnName)

	os.Setenv(driverName, driverVal)
	os.Setenv(dsnName, dsnVal)

	conf, err := common.GetConfig(".env_2")

	as.Equal(driverVal, conf.DbDriverName)
	as.Equal(dsnVal, conf.Dsn)
	as.Equal(nil, err)

	conf1, err1 := common.GetConfig(".env_5")
	as.NotEqual("", conf1.DbDriverName)
	as.NotEqual("", conf1.Dsn)
	as.Equal(driverVal, conf1.DbDriverName)
	as.Equal(dsnVal, conf1.Dsn)
	as.Equal(nil, err1)

	os.Unsetenv(driverName)
	os.Unsetenv(dsnName)
	conf3, err3 := common.GetConfig(".env_5")
	as.NotEqual("", conf3.DbDriverName)
	as.NotEqual("", conf3.Dsn)
	as.NotEqual(driverVal, conf3.DbDriverName)
	as.NotEqual(dsnVal, conf3.Dsn)
	as.Equal(nil, err3)
}
