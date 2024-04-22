package auth_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/bmtrann/sesc-component/config"
	"github.com/bmtrann/sesc-component/internal/auth"
	"github.com/bmtrann/sesc-component/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type Payload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func TestAuthServiceRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	authConfig := config.LoadAuthConfig()

	database := db.InitDB(&config.DBConfig{
		URI:  "mongodb://localhost:27017",
		Name: "UserDB",
	})
	handler := auth.InitAuthHandler(authConfig, database, "users")

	payload := Payload{
		Username: "test@c123456",
		Password: "test@123!",
	}
	jsonPayload, _ := json.Marshal(payload)
	w.Write(jsonPayload)

	handler.Register(c)

	assert.Equal(t, w.Code, 200)
}
