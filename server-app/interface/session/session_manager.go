package session

import "github.com/gin-gonic/gin"

type SessionManager interface {
	CreateSession(userID string) error
	GetSession(c *gin.Context) (string, error)
	DeleteSession(c *gin.Context) error
	UpdateSession(c *gin.Context, newID, oldID string) error
}
