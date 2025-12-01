package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"os"
	"strings"

	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	cognitoTypes "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var (
	cognitoClient *cognitoidentityprovider.Client
	userPoolID    string
	clientID      string
	clientSecret  string
)

func InitCognito() {
	// load .env file
	_ = godotenv.Load()

	userPoolID = os.Getenv("COGNITO_USER_POOL_ID")
	clientID = os.Getenv("COGNITO_CLIENT_ID")
	clientSecret = os.Getenv("COGNITO_CLIENT_SECRET")

	if userPoolID == "" || clientID == "" || clientSecret == "" {
		panic("missing COGNITO env vars")
	}

	// 🔥 IMPORTANT: set region explicitly, same as in main.go
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion("eu-north-1"),
	)
	if err != nil {
		panic(err)
	}

	cognitoClient = cognitoidentityprovider.NewFromConfig(cfg)
}

func computeSecretHash(username, clientID, clientSecret string) string {
	h := hmac.New(sha256.New, []byte(clientSecret))
	h.Write([]byte(username + clientID))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// ---------------------- REGISTER -----------------------

type RegisterBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}

func Register(c *gin.Context) {
	var body RegisterBody
	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	secretHash := computeSecretHash(body.Email, clientID, clientSecret)

	input := &cognitoidentityprovider.SignUpInput{
		ClientId:   aws.String(clientID),
		Username:   aws.String(body.Email),
		Password:   aws.String(body.Password),
		SecretHash: aws.String(secretHash),
		UserAttributes: []cognitoTypes.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(body.Email),
			},
			{
				Name:  aws.String("nickname"),
				Value: aws.String(body.Nickname),
			},
		},
	}

	_, err := cognitoClient.SignUp(context.Background(), input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// REMOVED: AdminConfirmSignUp call
	// SignUp automatically sends confirmation code to email
	// User will call /confirm endpoint with the code

	c.JSON(200, gin.H{"message": "registered, check your email for confirmation code"})
}

// ---------------------- LOGIN -----------------------

type LoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	var body LoginBody
	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	secretHash := computeSecretHash(body.Email, clientID, clientSecret)

	attemptAuth := func() (*cognitoidentityprovider.InitiateAuthOutput, error) {
		input := &cognitoidentityprovider.InitiateAuthInput{
			ClientId: aws.String(clientID),
			AuthFlow: cognitoTypes.AuthFlowTypeUserPasswordAuth,
			AuthParameters: map[string]string{
				"USERNAME":    body.Email,
				"PASSWORD":    body.Password,
				"SECRET_HASH": secretHash,
			},
		}
		return cognitoClient.InitiateAuth(context.Background(), input)
	}

	res, err := attemptAuth()
	if err != nil {
		if strings.Contains(err.Error(), "UserNotConfirmedException") {
			// try to confirm the user (requires IAM permission cognito-idp:AdminConfirmSignUp)
			if _, err2 := cognitoClient.AdminConfirmSignUp(context.Background(), &cognitoidentityprovider.AdminConfirmSignUpInput{
				UserPoolId: aws.String(userPoolID),
				Username:   aws.String(body.Email),
			}); err2 != nil {
				log.Printf("AdminConfirmSignUp during login failed for %s: %v", body.Email, err2)
				c.JSON(401, gin.H{"error": "user not confirmed and auto-confirm failed", "details": err2.Error()})
				return
			}
			// retry auth once
			res, err = attemptAuth()
		}
	}

	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	if res.AuthenticationResult == nil {
		c.JSON(500, gin.H{"error": "no authentication result"})
		return
	}

	c.JSON(200, gin.H{
		"id_token":      aws.ToString(res.AuthenticationResult.IdToken),
		"access_token":  aws.ToString(res.AuthenticationResult.AccessToken),
		"refresh_token": aws.ToString(res.AuthenticationResult.RefreshToken),
	})
}

// ---------------------- CONFIRMATION -----------------------

type ConfirmBody struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

// POST /confirm  -- confirm sign up with code from email
func Confirm(c *gin.Context) {
	var b ConfirmBody
	if err := c.BindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	secretHash := computeSecretHash(b.Email, clientID, clientSecret)

	input := &cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         aws.String(clientID),
		Username:         aws.String(b.Email),
		ConfirmationCode: aws.String(b.Code),
		// include SecretHash when app client has secret
		SecretHash: aws.String(secretHash),
	}

	_, err := cognitoClient.ConfirmSignUp(context.Background(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "confirmed"})
}

type ResendBody struct {
	Email string `json:"email"`
}

// POST /resend-confirm  -- resend confirmation code to email
func ResendConfirmation(c *gin.Context) {
	var b ResendBody
	if err := c.BindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	secretHash := computeSecretHash(b.Email, clientID, clientSecret)

	input := &cognitoidentityprovider.ResendConfirmationCodeInput{
		ClientId:   aws.String(clientID),
		Username:   aws.String(b.Email),
		SecretHash: aws.String(secretHash),
	}

	_, err := cognitoClient.ResendConfirmationCode(context.Background(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "code resent"})
}

// ---------------------- RESET PASSWORD -----------------------

type ChangePasswordBody struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func ChangePassword(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(401, gin.H{"error": "missing token"})
		return
	}

	token = strings.TrimPrefix(token, "Bearer ")

	var body ChangePasswordBody
	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	input := &cognitoidentityprovider.ChangePasswordInput{
		AccessToken:      aws.String(token),
		PreviousPassword: aws.String(body.OldPassword),
		ProposedPassword: aws.String(body.NewPassword),
	}

	_, err := cognitoClient.ChangePassword(context.Background(), input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})

		return
	}

	c.JSON(200, gin.H{"message": "password changed"})
}

// ---------------------- JWT MIDDLEWARE -----------------------

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
		if err != nil {
			c.JSON(401, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		if claims["token_use"] != "id" && claims["token_use"] != "access" {
			c.JSON(401, gin.H{"error": "invalid token type"})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}
