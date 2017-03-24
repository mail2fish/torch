package utilities

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"reflect"

	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// *mUser + ":" + *mPassword + "@" + *mHost + "/" + "?charset=utf8&parseTime=True&loc=Local"
func BuildConnectionStr(mHost *string, mUser *string, mPassword *string, mDatabase *string) string {
	var buffer bytes.Buffer
	buffer.WriteString(*mUser)
	if len(*mPassword) > 0 {
		buffer.WriteString(":")
		buffer.WriteString(*mPassword)
	}

	buffer.WriteString("@")
	buffer.WriteString(*mHost)
	buffer.WriteString("/")
	buffer.WriteString(*mDatabase)
	buffer.WriteString("?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=True&loc=UTC")

	return buffer.String()

}

func InitMysqlDB() *gorm.DB {
	mHost := flag.String("h", "unix(/tmp/mysql.sock)", "Specify the mysql server address. Default:unix(/tmp/mysql.sock)")
	mUser := flag.String("u", "root", "The username of connecting the mysql server. Default: root")
	mPassword := flag.String("p", "", "The password of connecting the mysql server.")
	mDatabase := flag.String("d", "poem_combat", "Specify the database of connecting the mysql server. Default: poem_combat")

	flag.Parse()

	connString := BuildConnectionStr(mHost, mUser, mPassword, mDatabase)

	db, err := gorm.Open("mysql", connString)
	if err != nil {
		switch err.(type) {
		case *mysql.MySQLError:
			mErr := err.(*mysql.MySQLError)
			if mErr.Number == 1049 {
				fmt.Println("Creating your database firstly")
				fmt.Println("create database poem_combat DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;")
				os.Exit(0)
			}
		}
		fmt.Println("DB connection string is :", connString)
		fmt.Println("Error is :", reflect.TypeOf(err), err)
		panic("Connecting mysql server failed.")

	}
	return db
}
