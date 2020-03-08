package app

import (
	"fmt"
	"github.com/spf13/viper"
)

// DBConnectionString returns DB credentials
func DBConnectionString() string {
	result := fmt.Sprintf("sslmode=disable host=%v port=%v dbname=%v user=%v",
		viper.GetString("db_host"),
		viper.GetString("db_port"),
		viper.GetString("db_name"),
		viper.GetString("db_user"),
	)
	if password := viper.GetString("db_password"); password != "" {
		result += " password=" + password
	}
	return result
}
