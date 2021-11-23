package utils

import (
	"github.com/gregoflash05/gradely/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
	Db *gorm.DB
)

func ConnectToDB(databaseUrl string) (*gorm.DB, error) {
	dbConnection, err := gorm.Open(mysql.Open(databaseUrl), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db, Db = dbConnection, dbConnection
	return db, nil

}

func Migrate() {
	db.AutoMigrate(&models.User{}, &models.Product{}, &models.Session{})
}
func DropTables() {
	db.Migrator().DropTable(&models.User{}, &models.Product{}, &models.Session{})
}

func GetItemsByField(model interface{}, field string, value interface{}) *gorm.DB {
	return db.Find(model, field+" = ?", value)
}
func GetItemByField(model interface{}, field string, value interface{}) *gorm.DB {
	return db.First(model, field+" = ?", value)
}

func GetItemByPrimaryKey(model interface{}, key interface{}) *gorm.DB {
	return db.First(model, key)
}

func CreateItem(data interface{}) *gorm.DB {
	return db.Create(data)
}
