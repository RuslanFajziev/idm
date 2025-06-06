package main

import (
	"fmt"
	"idm/inner/database"
	"time"
)

type DataMigration struct {
	Id      int64     `db:"id"`
	Version int64     `db:"version_id"`
	Applied bool      `db:"is_applied"`
	Create  time.Time `db:"tstamp"`
}

func main() {
	db := database.ConnectDb()
	defer db.Close()

	fmt.Println("**************************")
	fmt.Println("Check connetion to DB: Started")
	fmt.Println("**************************")

	var dataMigration DataMigration

	err := db.Get(&dataMigration, "SELECT * FROM goose_db_version WHERE id = $1", 4)

	if err != nil {
		panic(fmt.Errorf("failed to query database: %v", err))
	}

	fmt.Printf("version:%v applied:%v\n", dataMigration.Version, dataMigration.Applied)

	if dataMigration.Applied {
		fmt.Println("**************************")
		fmt.Println("Check connetion to DB: Passed")
		fmt.Println("**************************")
	}
}
