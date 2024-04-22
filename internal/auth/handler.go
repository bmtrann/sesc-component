package auth

import (
	"context"
	"crypto/sha1"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bmtrann/sesc-component/config"
	"github.com/bmtrann/sesc-component/internal/exception"
	model "github.com/bmtrann/sesc-component/internal/model/user"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthService interface {
	Login()
	Register()
}

type AuthHandler struct {
	userRepo *model.UserRepository
	config   *config.AuthConfig
}

type Payload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Token string `json:"token"`
}

type authClaims struct {
	*model.User `json:"user"`
	jwt.RegisteredClaims
}

func InitAuthHandler(config *config.AuthConfig, db *mongo.Database, collection string) *AuthHandler {
	return &AuthHandler{
		model.NewUserRepository(db, collection),
		config,
	}
}

func (handler *AuthHandler) Login(c *gin.Context) {
	payload := new(Payload)

	if err := c.BindJSON(payload); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	token, err := handler.login(c.Request.Context(), payload.Username, payload.Password)
	if err != nil {
		if err == exception.ErrUserNotFound {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, Response{Token: token})
}

func (handler *AuthHandler) login(ctx context.Context, username, password string) (string, error) {
	pwd := sha1.New()
	pwd.Write([]byte(password))
	pwd.Write([]byte(handler.config.HashSalt))
	password = fmt.Sprintf("%x", pwd.Sum(nil))

	user, err := handler.userRepo.GetUser(ctx, username, password)
	if err != nil {
		return "", exception.ErrUserNotFound
	}

	claims := authClaims{
		user,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 0, 1)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(handler.config.SignKey)
}

func (handler *AuthHandler) Register(c *gin.Context) {
	payload := new(Payload)

	if err := c.BindJSON(payload); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err := handler.register(c.Request.Context(), payload)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.Status(http.StatusCreated)
}

func (handler *AuthHandler) register(ctx context.Context, payload *Payload) error {
	pwd := sha1.New()
	pwd.Write([]byte(payload.Password))
	pwd.Write([]byte(handler.config.HashSalt))

	user := &model.MongoUser{
		Username: payload.Username,
		Password: fmt.Sprintf("%x", pwd.Sum(nil)),
	}

	return handler.userRepo.CreateUser(ctx, user)
}

func (handler *AuthHandler) ParseToken(ctx context.Context, accessToken string) (*model.User, error) {
	token, err := jwt.ParseWithClaims(accessToken, &authClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return handler.config.SignKey, nil
	})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if claims, ok := token.Claims.(*authClaims); ok && token.Valid {
		return claims.User, nil
	}

	return nil, exception.ErrInvalidAccessToken
}
