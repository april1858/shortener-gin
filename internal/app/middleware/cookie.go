package middleware

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"

	"github.com/april1858/shortener-gin/internal/app/repository"
	"github.com/gin-gonic/gin"
)

var (
	ErrValueTooLong = errors.New("cookie value too long")
	ErrInvalidValue = errors.New("invalid cookie value")
	secretKey       = []byte("12345")
	cookie          = http.Cookie{
		Name:     "UID",
		Value:    "",
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
)

func (mw *MW) Cookie() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get cookie
		if signedValue, err := c.Cookie("UID"); err == nil {
			value, err := ReadSigned(signedValue)
			if err != nil {
				fmt.Println("err err - ", err)
			}
			repository.UID = value
			c.Next()
			return
		}

		err := WriteSigned()
		if err != nil {
			fmt.Println("err from Cookie() ", err)
		}
		c.SetCookie(cookie.Name, cookie.Value, 3600, "/", "localhost", true, true)
		c.Next()
	}
}

func ReadSigned(sValue string) (string, error) {
	signedValue, err := base64.URLEncoding.DecodeString(sValue)
	if err != nil {
		return "", ErrInvalidValue
	}
	if len(signedValue) < sha256.Size {
		return "", ErrInvalidValue
	}
	signature := signedValue[:sha256.Size]
	value := signedValue[sha256.Size:]

	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte("UID"))
	mac.Write([]byte(value))
	expectedSignature := mac.Sum(nil)

	if !hmac.Equal([]byte(signature), expectedSignature) {
		return "", ErrInvalidValue
	}

	return string(value), nil
}

func WriteSigned() error {
	cookie.Value = createCode()
	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(cookie.Name))
	mac.Write([]byte(cookie.Value))
	signature := mac.Sum(nil)
	cookie.Value = string(signature) + cookie.Value
	cookie.Value = base64.URLEncoding.EncodeToString([]byte(cookie.Value))
	return nil
}

func createCode() string {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return "error in createCode()"
	}
	repository.UID = hex.EncodeToString(b)
	fmt.Println("UID from create - ", hex.EncodeToString(b))
	return hex.EncodeToString(b)
}
