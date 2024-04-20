package server

import (
	"log"

	"github.com/bmtrann/sesc-component/config"
	"github.com/bmtrann/sesc-component/internal/auth"
	database "github.com/bmtrann/sesc-component/internal/db"
	"github.com/bmtrann/sesc-component/internal/enrolment"
	"github.com/bmtrann/sesc-component/internal/profile"
	"github.com/gin-gonic/gin"
)

func InitApp() {
	// load config
	if err := config.Init(); err != nil {
		log.Fatalf("%s", err.Error())
	}

	authConfig := config.LoadAuthConfig()
	dbConfig := config.LoadDBConfig()

	router := gin.Default()
	db := database.InitDB(dbConfig)

	// add Auth endpoints
	authHandler := auth.InitAuthHandler(authConfig, db, dbConfig.UserCollection)
	auth.AddRoutes(router, authHandler)

	// add Enrolment endpoints
	enrolmentHandler := enrolment.InitEnrolmentHandler(db, dbConfig)
	enrolment.AddRoutes(router, enrolmentHandler)

	// add Profile endpoints
	profileHandler := profile.InitProfileHandler(db, dbConfig.StudentCollection)
	profile.AddRoutes(router, profileHandler)

	router.Run()
}
