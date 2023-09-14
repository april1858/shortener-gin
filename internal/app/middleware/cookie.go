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

	"github.com/gin-gonic/gin"
)

var (
	ErrValueTooLong = errors.New("cookie value too long")
	ErrInvalidValue = errors.New("invalid cookie value")
	secretKey       = []byte("12345")
	cookie          http.Cookie
)

func (mw *MW) Cookie() gin.HandlerFunc {
	return func(c *gin.Context) {
		mw.mx.Lock()
		defer mw.mx.Unlock()
		// Get cookie
		if signedValue, err := c.Cookie("UID"); err == nil {
			_, err := ReadSigned(signedValue)
			if err != nil {
				fmt.Println("error from ReadSigned - ", err)
			}
			//c.Set("UID", value)
			c.Next()
			return
		}

		_, err := WriteSigned()
		fmt.Println("Get cookie - ", cookie.Name, "+",cookie.Value)
		if err != nil {
			fmt.Println("error from WriteSigned - ", err)
		}
		cookie.Value = "UID"
		c.SetCookie(cookie.Value, cookie.Value, 3600, "/", "", false, false)
		//c.Set("UID", uid)
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

func WriteSigned() (string, error) {
	uid := createCode()
	cookie.Value = uid
	mac := hmac.New(sha256.New, secretKey)
	_, err := mac.Write([]byte(cookie.Name))
	if err != nil {
		return "", err
	}
	_, err = mac.Write([]byte(cookie.Value))
	if err != nil {
		return "", err
	}
	signature := mac.Sum(nil)
	cookie.Value = string(signature) + cookie.Value
	cookie.Value = base64.URLEncoding.EncodeToString([]byte(cookie.Value))
	return uid, nil
}

func createCode() string {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return "error in createCode()"
	}
	return hex.EncodeToString(b)
}
