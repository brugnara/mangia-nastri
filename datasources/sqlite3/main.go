package sqlite3

import (
	"mangia_nastri/datasources"
	"mangia_nastri/logger"
	"sync"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SQLiteDataSource struct {
	log logger.Logger
	db  *gorm.DB
}

type dbStruct struct {
	gorm.Model
	datasources.Payload
	ID datasources.Hash `gorm:"primaryKey"`
}

func New(log *logger.Logger, fileLocation string) *SQLiteDataSource {
	db, err := gorm.Open(sqlite.Open(fileLocation), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error opening SQLite database: %v", err)
	}

	db.AutoMigrate(&dbStruct{})

	return &SQLiteDataSource{
		log: log.CloneWithPrefix("sqlite"),
		db:  db,
	}
}

func (ds *SQLiteDataSource) Ready() <-chan bool {
	// do a PING and check the response is PONG
	ch := make(chan bool)
	wg := sync.WaitGroup{}

	// try 5 times then fail
	go func() {
		connected := false
		for i := 0; i < 5 && !connected; i++ {
			wg.Add(1)

			go func() {
				err := ds.db.Exec("SELECT 1").Error

				wg.Done()

				if err != nil {
					ds.log.Printf("Error pinging DB: %v", err)
					return
				}

				connected = true
			}()

			wg.Wait()

			time.Sleep(1 * time.Second)
		}

		if !connected {
			ds.log.Fatalf("DB is not ready after 5 attempts")
		} else {
			ds.log.Info("DB is ready")
			close(ch)
		}
	}()

	return ch
}

func (ds *SQLiteDataSource) Set(key datasources.Hash, value datasources.Payload) error {
	ds.log.Printf("SET: Setting value %v for key %v", value, key)

	dbData := dbStruct{
		ID:      key,
		Payload: value,
	}

	result := ds.db.Create(&dbData)

	if result.Error != nil {
		ds.log.Printf("SET: Error setting value for key %v: %v", key, result.Error)
		return result.Error
	}

	return nil
}

func (ds *SQLiteDataSource) Get(key datasources.Hash) (payload datasources.Payload, err error) {
	ds.log.Printf("GET: Getting value for key %v", key)

	var dbData dbStruct
	result := ds.db.First(&dbData, "id = ?", key)

	if result.Error != nil {
		ds.log.Printf("GET: Error getting value for key %v: %v", key, result.Error)
		return payload, result.Error
	}

	return dbData.Payload, nil
}
