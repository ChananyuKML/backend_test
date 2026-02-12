package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"hole/adapters"
	_ "hole/docs" // IMPORTANT
	"hole/entities"
	"hole/use_cases"
)

const (
	host     = "postgres"   // or the Docker service name if running in another container
	port     = 5432         // default PostgreSQL port
	user     = "myuser"     // as defined in docker-compose.yml
	password = "mypassword" // as defined in docker-compose.yml
	dbname   = "auth"       // as defined in docker-compose.yml
)

// @title Hole Auth API
// @version 1.0
// @description Authentication API with JWT & Refresh Token
// @host localhost:8000
// @BasePath
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {

	app := fiber.New()
	app.Get("/swagger/*", swagger.HandlerDefault)
	godotenv.Load()

	// Configure your PostgreSQL database details here
	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	//dsn := fmt.Sprintf("host=%s port=%s user=%s "+
	// 	"password=%s dbname=%s sslmode=disable",
	// 	os.Getenv("DB_HOST"),
	// 	os.Getenv("DB_PORT"),
	// 	os.Getenv("DB_USER"),
	// 	os.Getenv("DB_PASSWORD"),
	// 	os.Getenv("DB_NAME"))

	// New logger for detailed SQL logging
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Enable color
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger, // add Logger
	})

	if err != nil {
		panic("failed to connect to database")
	}
	db.AutoMigrate(&entities.User{},
		&entities.RefreshToken{},
		&entities.Item{},
	)

	fmt.Println("Database migration completed!")

	userRepo := adapters.NewUserRepository(db)
	refreshRepo := adapters.NewRefreshTokenRepository(db)
	itemRepo := adapters.NewItemRepository(db)
	tokenSvc := adapters.NewJWTService()

	authUC := use_cases.NewAuthUseCase(
		userRepo,
		refreshRepo,
		tokenSvc,
	)

	itemUC := use_cases.NewItemUseCase(
		itemRepo,
	)

	itemHandler := adapters.NewItemHandler(itemUC)
	authHandler := adapters.NewAuthHandler(authUC)

	app.Post("/register", authHandler.Register)
	app.Post("/login", authHandler.Login)
	app.Post("/refresh", authHandler.Refresh)

	reg := app.Group("/register", adapters.Protected(tokenSvc))
	reg.Post("/sam", itemHandler.Create)
	reg.Post("/box", itemHandler.List)

	insp := app.Group("/inspect", adapters.Protected(tokenSvc))
	insp.Post("/:id", itemHandler.Create)
	app.Listen(":8000")
}
