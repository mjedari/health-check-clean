package cmd

import (
	"database/sql"
	"fmt"
	"github.com/mjedari/health-checker/app/config"
	"github.com/mjedari/health-checker/domain"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrating database",
	Run: func(cmd *cobra.Command, args []string) {
		migrate()
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}

func migrate() {
	fmt.Println("migration started...")

	conf := config.Config.MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", conf.User, conf.Pass, conf.Host, conf.Port, conf.Database)

	// Connect to MySQL server
	msqldb, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer msqldb.Close()

	// Create the database
	_, err = msqldb.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", conf.Database))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Database created successfully")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}

	// Perform auto migration
	err = db.AutoMigrate(&domain.Endpoint{})
	if err != nil {
		log.Fatal("failed to migrate database: ", err)
	}

	fmt.Println("migration finished.")
}
