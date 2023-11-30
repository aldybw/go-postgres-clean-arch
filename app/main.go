package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"

	"go-postgres-clean-arch/config"
	_tagHttpDelivery "go-postgres-clean-arch/tag/delivery/http"
	_tagHttpDeliveryMiddleware "go-postgres-clean-arch/tag/delivery/http/middleware"
	_tagRepo "go-postgres-clean-arch/tag/repository/postgresql"
	_tagUcase "go-postgres-clean-arch/tag/usecase"
)

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	// sql/database golang connection (for raw queries)
	// dbHost := viper.GetString(`database.host`)
	// dbPort := viper.GetString(`database.port`)
	// dbUser := viper.GetString(`database.user`)
	// dbPass := viper.GetString(`database.pass`)
	// dbName := viper.GetString(`database.name`)
	// connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	// connection := fmt.Sprintf("postgres://%s:%s@%s/%s:%s?sslmode=verify-full", dbUser, dbPass, dbHost, dbPort, dbName)
	// connection := "postgres://postgres@localhost/clean_arch_test?sslmode=verify-full"
	dsn := "host=localhost port=5432 user=postgres dbname=clean_arch_test sslmode=disable"
	// val := url.Values{}
	// val.Add("parseTime", "1")
	// val.Add("loc", "Asia/Jakarta")
	// dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := sql.Open(`postgres`, dsn)

	// gorm connection (for ORM)
	db := config.DatabaseConnection()

	if err != nil {
		log.Fatal(err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	e := echo.New()
	middL := _tagHttpDeliveryMiddleware.InitMiddleware()
	e.Use(middL.CORS)
	authorRepo := _tagRepo.NewPostgresqlTagRepository(dbConn, db)

	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	var validate *validator.Validate
	tu := _tagUcase.NewTagUsecase(authorRepo, timeoutContext, validate)
	_tagHttpDelivery.NewTagHandler(e, tu)

	log.Fatal(e.Start(viper.GetString("server.address"))) //nolint

	// log.Info().Msg("Started Server!")
	// Database
	// db := config.DatabaseConnection()
	// validate := validator.New()

	// db.Table("tags").AutoMigrate(&domain.Tag{})

	// // Repository
	// tagRepository := _tagRepo.NewPostgresqlTagRepository(db)

	// // Usecase
	// tagUsecase := _tagUcase.NewTagUsecase(tagRepository, validate)

	// // Delivery
	// tagDelivery := _tagHttpDelivery.NewTagHandler(tagUsecase)

	// // Router
	// routes := router.NewRouter(tagDelivery)

	// server := &http.Server{
	// 	Addr:    ":8080",
	// 	Handler: routes,
	// }

	// err := server.ListenAndServe()
	// helper.ErrorPanic(err)
}
