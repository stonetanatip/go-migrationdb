package main

import (
	"database/sql"
	"fmt"
	"go-migrationdb/services/migrate"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

func main() {
	db, err := initLocalDB()
	if err != nil {
		fmt.Println("could not connect to the MySQL database... ", err)
		panic(err)
	}

	if err := db.Ping(); err != nil {
		fmt.Println("could not ping DB... ", err)
		panic(err)
	}

	service := migrate.NewService(db)
	migrate := migrate.NewHandler(service)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/health", makeHealthHandler(db))
	e.POST("/up", migrate.UpDB)
	e.POST("/down", migrate.DownDB)

	e.Logger.Fatal(e.Start(":1323"))
}

func initLocalDB() (*sql.DB, error) {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v",
		viper.GetString("mysql.username"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.database.bank"),
	)

	fmt.Println("Start DB :", dsn)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func makeHealthHandler(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := db.Ping()
		if err != nil {
			fmt.Println("db ping error: ", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"status": "unhealty",
				"msg":    fmt.Sprintf("db ping error: %s", err.Error()),
			})
		}

		return c.JSON(http.StatusOK, map[string]string{"status": "healthy"})
	}
}
