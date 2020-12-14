package main

import (
	"github.com/dgrijalva/jwt-go"
	db "github.com/kochetov-dmitrij/challenge_it_backend/database"
	_ "github.com/kochetov-dmitrij/challenge_it_backend/docs/echo_server"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"os"
)

type jwtCustomClaims struct {
	Name   string `json:"name"`
	UserId int32  `json:"uid"`
	jwt.StandardClaims
}

// @title Echo Swagger Example API
// @version 1.0
// @description This is a sample server server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /
// @schemes http
func main() {
	// Init db
	_ = db.InitDB()

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes
	e.GET("/", HealthCheck)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Auth
	e.POST("/login", Login)
	e.POST("/register", Register)

	r := e.Group("/challenge")

	// Configure middleware with the custom claims type
	config := middleware.JWTConfig{
		Claims:     &jwtCustomClaims{},
		SigningKey: []byte("secret"),
	}

	r.Use(middleware.JWTWithConfig(config))

	r.POST("/new", NewChallenge)
	r.GET("/created", CreatedChallenges)
	r.GET("/take", TakeChallenge)
	r.GET("/my", MyChallenges)
	r.GET("/all", AllChallenges)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
