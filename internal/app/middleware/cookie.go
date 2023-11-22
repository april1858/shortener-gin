package middleware

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

var (
	ErrValueTooLong = errors.New("cookie value too long")
	ErrInvalidValue = errors.New("invalid cookie value")
	secretKey       = []byte("12345")
)

func (mw *MW) Cookie() gin.HandlerFunc {
	return func(c *gin.Context) {
		if signedValue, err := c.Cookie("UID"); err == nil {
			value, err := ReadSigned(signedValue)
			if err != nil {
				fmt.Println("error from ReadSigned - ", err)
			}
			c.Set("UID", value)
			c.Next()
			return
		}
		name := "UID"
		uid, signedValue, err := WriteSigned(name)
		if err != nil {
			fmt.Println("error from WriteSigned - ", err)
		}
		c.SetCookie(name, signedValue, 3600, "/", "", false, false)
		c.Set("UID", uid)
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
	_, err = mac.Write([]byte("UID"))
	if err != nil {
		return "", err
	}
	_, err = mac.Write([]byte(value))
	if err != nil {
		return "", err
	}
	expectedSignature := mac.Sum(nil)

	if !hmac.Equal([]byte(signature), expectedSignature) {
		return "", ErrInvalidValue
	}

	return string(value), nil
}

func WriteSigned(name string) (string, string, error) {
	uid, err := createCode()
	if err != nil {
		return "", "", err
	}
	value := uid
	mac := hmac.New(sha256.New, secretKey)
	_, err = mac.Write([]byte(name))
	if err != nil {
		return "", "", err
	}
	_, err = mac.Write([]byte(value))
	if err != nil {
		return "", "", err
	}
	signature := mac.Sum(nil)
	value = string(signature) + value
	value = base64.URLEncoding.EncodeToString([]byte(value))
	return uid, value, nil
}

func createCode() (string, error) {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
