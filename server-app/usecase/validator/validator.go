package validator

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// バリデーションチェック処理
func ValidationCheck(c *gin.Context, err error) error {

	// JSONデータをUser構造体にバインドしてバリデーションを実行
	if err != nil {
		// バリデーションエラーが発生した場合はエラーレスポンスを返す
		var verr validator.ValidationErrors
		if ok := errors.As(err, &verr); ok {
			var errorMsgs []string
			for _, e := range verr {
				errorMsgs = append(errorMsgs, fmt.Sprintf("%s validation failed on the %s field", e.Tag(), e.Field()))
			}
			log.Println(errorMsgs)
			return errors.New(strings.Join(errorMsgs, ", "))
			// エラーメッセージを結合したものをエラーとして返す
		}
		log.Println(err.Error())
		return err
	}
	return nil
}
