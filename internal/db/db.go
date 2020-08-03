package db

import (
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const sleepDuration = 5

//Db ...
type Db struct {
	Database *gorm.DB
	Adress   string
}

//New ...
func New(adress string) (*Db, error) {
	var err error
	db := &Db{}
	db.Adress = adress
	db.Database, err = gorm.Open("postgres", adress)
	if err != nil {
		return nil, fmt.Errorf("db.New: %v", err)
	}

	err = db.ping()
	if err != nil {
		return nil, fmt.Errorf("db.New: %s", err.Error())
	}
	db.checkConnection()

	return db, nil
}

func (db *Db) Close() error {
	err := db.Database.Close()
	if err != nil {
		return fmt.Errorf("db.Close: %v", err)
	}
	return nil
}

//ping (internal)
func (db *Db) ping() error {
	//db.Database.LogMode(true)
	result := struct {
		Result int
	}{}

	err := db.Database.Raw("select 1+1 as result").Scan(&result).Error
	if err != nil {
		return fmt.Errorf("db.ping: %v", err)
	}
	if result.Result != 2 {
		return fmt.Errorf("db.ping: incorrect result!=2 (%d)", result.Result)
	}
	return nil
}

//checkConnection (internal)
func (db *Db) checkConnection() {
	go func() {
		for {
			err := db.ping()
			if err != nil {

				log.Printf("db.checkConnection: no connection: %s", err.Error())
				tempDb, err := gorm.Open("postgres", db.Adress)

				if err != nil {
					log.Printf("db.checkConnection: could not establish connection: %v", err)
				} else {
					db.Database = tempDb
				}
			}
			time.Sleep(sleepDuration * time.Second)
		}
	}()
}
