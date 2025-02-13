package login

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kazukimurahashi12/webapp/model/redis"
)

// func TestGetLogin(t *testing.T) {
// 	type args struct {
// 		c *gin.Context
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want string
// 	}{
// 		//正常系
// 		c *gin.Context
// 		{"テスト1", args{c}, ""},
// 		//異常系
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			GetLogin(tt.args.c)
// 		})
// 	}
// }

func TestGetLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// ダミーのgin.Contextを作成
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set("UserId", "dummyUserId") // セッション情報をセット

	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name           string
		args           args
		expectedStatus int
	}{
		// 正常系: セッションからユーザー情報を取得できる場合
		{"正常系: ユーザー情報取得成功", args{c}, http.StatusOK},

		// 正常系: セッションからユーザー情報が取得できない場合
		// {"正常系: セッション情報なし", args{gin.CreateTestContext(httptest.NewRecorder())}, http.StatusUnauthorized},

		// 異常系: ユーザー情報取得時にエラーが発生する場合
		{"異常系: ユーザー情報取得エラー", args{c}, http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			redis := redis.NewRedisSessionStore()
			// c, _ := gin.CreateTestContext(response)
			GetLogin(tt.args.c, redis)
			if response.Code != tt.expectedStatus {
				t.Errorf("HTTP Statusコードが期待値と一致しません。期待値: %d, 実際: %d", tt.expectedStatus, response.Code)
			}
		})
	}
}
