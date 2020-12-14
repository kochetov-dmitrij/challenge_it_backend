package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-pg/pg/v10"
	db "github.com/kochetov-dmitrij/challenge_it_backend/database"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/sha3"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var database *pg.DB

const SECRET string = "secret" // TODO: change to env

// HealthCheck godoc
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router / [get]
func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": "Server is up and running",
	})
}

func Register(c echo.Context) error {
	var b map[string]interface{}
	err := json.NewDecoder(c.Request().Body).Decode(&b)
	username := fmt.Sprintf("%v", b["email"])
	name := fmt.Sprintf("%v", b["name"])
	password := fmt.Sprintf("%v", b["password"])

	hasher := sha3.Sum256([]byte(username + ":" + password))
	hash := fmt.Sprintf("%x", hasher)
	user := db.User{
		Name:        name,
		Email:       username,
		EncPassword: hash,
	}
	_, err = database.Model(&user).Insert()
	if err != nil {
		return c.String(http.StatusInternalServerError, "something goes wrong")
	}

	claims := &jwtCustomClaims{
		name,
		user.Id,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(SECRET))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

func Login(c echo.Context) error {
	var b map[string]interface{}
	err := json.NewDecoder(c.Request().Body).Decode(&b)
	username := fmt.Sprintf("%v", b["email"])
	password := fmt.Sprintf("%v", b["password"])

	user := &db.User{}
	if err := database.Model(user).Where("email = ?", username).Select(); err != nil {
		return c.String(http.StatusInternalServerError, "something goes wrong")
	}

	hasher := sha3.Sum256([]byte(username + ":" + password))
	hash := fmt.Sprintf("%x", hasher)

	if strings.Compare(hash, user.EncPassword) != 0 {
		return echo.ErrUnauthorized
	}

	// Set custom claims
	claims := &jwtCustomClaims{
		user.Name,
		user.Id,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(SECRET))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

func RejectChallenge(c echo.Context) error {
	challenge, err := strconv.Atoi(c.QueryParam("challenge"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid challenge")
	}
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	mychallenge := db.UserChallenge{
		Id: int32(challenge),
	}
	_, err = database.Model(&mychallenge).Set("status = ?", db.Rejected).Where("user_id = ?", claims.UserId).Where("id = ?", challenge).Update()
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Something goes wrong")
	}
	return c.String(http.StatusOK, "Rejected")
}

func CompleteChallenge(c echo.Context) error {
	challenge, err := strconv.Atoi(c.QueryParam("challenge"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid challenge")
	}
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	mychallenge := db.UserChallenge{
		Id: int32(challenge),
	}
	_, err = database.Model(&mychallenge).Set("status = ?", db.Completed).Where("user_id = ?", claims.UserId).Where("id = ?", challenge).Update()
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Something goes wrong")
	}
	return c.String(http.StatusOK, "Completed")
}

func TakeChallenge(c echo.Context) error {
	challenge, err := strconv.Atoi(c.QueryParam("challenge"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid challenge")
	}
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	var mychallenge db.UserChallenge
	mychallenge.UserId = claims.UserId
	mychallenge.StartDate = time.Now()
	mychallenge.Status = db.InProgress
	mychallenge.ChallengeId = int32(challenge)
	_, err = database.Model(&mychallenge).Insert()
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Something goes wrong")
	}
	return c.String(http.StatusCreated, "Taken")
}

func NewChallenge(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	var challenge db.Challenge
	challenge.AuthorId = claims.UserId
	err := json.NewDecoder(c.Request().Body).Decode(&challenge)
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "Your challenge data is invalid")
	}
	_, err = database.Model(&challenge).Insert()
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Something goes wrong")
	}
	return c.String(http.StatusCreated, "Added")
}

func CreatedChallenges(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	var challenges []db.Challenge
	err := database.Model(&challenges).Where("challenge.author_id = ?", claims.UserId).Select()
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Something goes wrong")
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"challenges": challenges})
}

func MyChallenges(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	var challenges []db.UserChallenge
	err := database.Model(&challenges).Where("user_id = ?", claims.UserId).Select()
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Something goes wrong")
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"challenges": challenges})
}

func AllChallenges(c echo.Context) error {
	var challenges []db.Challenge
	err := database.Model(&challenges).Select()
	if err != nil {
		log.Error(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Something goes wrong")
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"challenges": challenges})
}

func init() {
	database, _ = db.Connect()
}
