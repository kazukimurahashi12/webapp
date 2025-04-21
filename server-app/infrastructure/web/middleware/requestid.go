package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ctxKey string

const RequestIDKey ctxKey = "requestID"
const HeaderXRequestID = "X-Request-ID"

// ヘッダからリクエストIDを取得し、なければ新規生成してContextに詰めるミドルウェア
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(HeaderXRequestID)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Contextに詰める
		ctx := context.WithValue(c.Request.Context(), RequestIDKey, requestID)
		c.Request = c.Request.WithContext(ctx)

		// レスポンスヘッダーにも追加
		c.Writer.Header().Set(HeaderXRequestID, requestID)

		c.Next()
	}
}

// ContextからrequestIDを取得
func GetRequestID(ctx context.Context) string {
	if v, ok := ctx.Value(RequestIDKey).(string); ok {
		return v
	}
	return ""
}
