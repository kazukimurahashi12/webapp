package common

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazukimurahashi12/webapp/interface/session"
)

type CommonController struct {
	sessionManager session.SessionManager
}

func NewCommonController(sessionManager session.SessionManager) *CommonController {
	return &CommonController{
		sessionManager: sessionManager,
	}
}

// セッションからログインIDを取得するAPI
func (c *CommonController) GetLoginIdBySession(ctx *gin.Context) {
	// セッションからIDを取得
	id, err := c.sessionManager.GetSession(ctx)
	if err != nil {
		log.Printf("セッションからIDの取得に失敗しました。error: %v", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Success Get LoginId bySession :id %+v", id)
	ctx.JSON(http.StatusOK, gin.H{"id": id})
}
