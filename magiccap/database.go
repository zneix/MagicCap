package magiccap

import (
	"database/sql"
	"encoding/json"
	"path"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var (
	// Database defines the database which MagicCap uses.
	Database, _ = sql.Open("sqlite3", path.Join(ConfigPath, "magiccap.db"))

	// DatabaseLock defines the database lock.
	DatabaseLock = sync.Mutex{}

	// ConfigItems defines all the config options which have been set.
	ConfigItems = map[string]interface{}{}

	// ConfigItemsLock is the R/W thread lock for the config.
	ConfigItemsLock = sync.RWMutex{}
)

// LoadDatabase loads in the database schemas.
func LoadDatabase() {
	// Creates the config table.
	_, err := Database.Exec("CREATE TABLE IF NOT EXISTS `config` (`key` TEXT NOT NULL, `value` TEXT NOT NULL)")
	if err != nil {
		panic(err)
	}

	// Creates the captures table.
	_, err = Database.Exec("CREATE TABLE IF NOT EXISTS `captures` (`filename` TEXT NOT NULL, `success` INTEGER NOT NULL, `timestamp` INTEGER NOT NULL, `url` TEXT, `file_path` TEXT)")
	if err != nil {
		panic(err)
	}
	_, err = Database.Exec("CREATE INDEX IF NOT EXISTS TimestampIndex ON captures(timestamp)")
	if err != nil {
		panic(err)
	}

	// Creates the tokens table.
	_, err = Database.Exec("CREATE TABLE IF NOT EXISTS tokens (token TEXT NOT NULL, expires INTEGER NOT NULL, uploader TEXT NOT NULL)")
	if err != nil {
		panic(err)
	}

	// Gets all the config items.
	rows, err := Database.Query("SELECT * FROM config")
	if err != nil {
		panic(err)
	}
	ConfigItemsLock.Lock()
	for rows.Next() {
		var Key string
		var Value string
		err = rows.Scan(&Key, &Value)
		if err != nil {
			panic(err)
		}
		var GenericInterface interface{}
		err = json.Unmarshal([]byte(Value), &GenericInterface)
		if err != nil {
			panic(err)
		}
		ConfigItems[Key] = GenericInterface
	}
	ConfigItemsLock.Unlock()

	// Log that the database is initialised.
	println("Database initialised.")
}

// UpdateConfig is used to update the config in the database.
func UpdateConfig() {
	DatabaseLock.Lock()
	_, err := Database.Exec("DELETE FROM config")
	if err != nil {
		panic(err)
	}
	Statement, err := Database.Prepare("INSERT INTO config VALUES (?, ?)")
	if err != nil {
		panic(err)
	}
	ConfigItemsLock.RLock()
	for k, v := range ConfigItems {
		b, err := json.Marshal(&v)
		if err != nil {
			panic(err)
		}
		_, err = Statement.Exec(k, string(b))
		if err != nil {
			panic(err)
		}
	}
	ConfigItemsLock.RUnlock()
	DatabaseLock.Unlock()
}