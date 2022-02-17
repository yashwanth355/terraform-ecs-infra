package repositories

import (
	"fmt"
	"net/url"
	Config "usermvc/config"

	"github.com/jinzhu/gorm"
)

func NewDb() (*gorm.DB, error) {
	conf := Config.LoadConfig()

	var (
		host     = conf.DbHost
		port     = conf.DbPort
		user     = conf.DbUserName
		password = conf.DbPassword
		dbname   = conf.DbName
	)

	fmt.Println("printing the user and password", conf.DbUserName, conf.DbPassword)
	dsn := url.URL{
		User:     url.UserPassword(user, password),
		Scheme:   "postgres",
		Host:     fmt.Sprintf("%s:%d", host, port),
		Path:     dbname,
		RawQuery: (&url.Values{"sslmode": []string{"disable"}}).Encode(),
	}
	fmt.Println(password)
	db, err := gorm.Open("postgres", dsn.String())
	if err != nil {
		return nil, err
	}
	return db, nil
}
