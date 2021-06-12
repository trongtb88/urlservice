package tests


import (
"fmt"
"log"
"os"
"testing"

"github.com/joho/godotenv"
"github.com/trongtb88/urlservice/api/controllers"
"github.com/trongtb88/urlservice/api/entity"

"gorm.io/driver/mysql"
"gorm.io/gorm"
)

var server = controllers.Server{}
var urlModel = entity.Url{}

func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load(os.ExpandEnv("./../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}
	Database()
	os.Exit(m.Run())

}

func Database() {

	var err error

	TestDbDriver := os.Getenv("TEST_DB_DRIVER")
	if TestDbDriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			os.Getenv("TEST_DB_USER"), os.Getenv("TEST_DB_PASSWORD"), os.Getenv("TEST_DB_HOST"), os.Getenv("TEST_DB_PORT"), os.Getenv("TEST_DB_NAME"))
		server.DB, err = gorm.Open(mysql.Open(DBURL))
		if err != nil {
			fmt.Printf("Cannot connect to %s database", TestDbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database", TestDbDriver)
		}
	}

	server.DB.Debug().AutoMigrate(&entity.Url{})
}

