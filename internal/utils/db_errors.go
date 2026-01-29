package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

type Status string

const (
	StatusSuccess Status = "Success"
	StatusFailure Status = "Failure"
)

func ParsePostgresError(err error) (string, bool) {
	if err == nil {
		return "", false
	}

	errStr := err.Error()
	errorCode := ExtractPgCode(err)
	fmt.Println("errorCode :", errorCode)

	switch errorCode {
	case "23505":
		switch {
		case strings.Contains(errStr, "users_username_key"):
			return "username already exists", true
		case strings.Contains(errStr, "users_email_key"):
			return "email already exists", true
		default:
			return "duplicate value already exists", true
		}
	}

	return "", false
}

func Error(c *gin.Context, code int, message string, details interface{}) {
	c.JSON(code, gin.H{
		"status":  "Failure",
		"message": message,
		"error":   details,
	})
}

func ExtractPgCode(err error) string {
	if err == nil {
		return ""
	}

	re := regexp.MustCompile(`SQLSTATE (\d{5})`)
	matches := re.FindStringSubmatch(err.Error())
	if len(matches) == 2 {
		return matches[1] // "23505"
	}
	return ""
}

func ResponseError(c *gin.Context, code int, message string, details interface{}) {
	c.JSON(code, gin.H{
		"status":  StatusFailure,
		"message": message,
		"error":   details,
	})
}

func ResponseSuccess(c *gin.Context, code int, message string, details interface{}) {
	c.JSON(code, gin.H{
		"status":  StatusSuccess,
		"message": message,
		"data":    details,
	})
}
