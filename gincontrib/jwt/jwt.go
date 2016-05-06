package jwt

import (
	jwt_lib "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"fmt"
	"github.com/alivinco/blackflowhub/model"
)

func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := jwt_lib.ParseFromRequest(c.Request, func(token *jwt_lib.Token) (interface{}, error) {
			b := ([]byte(secret))
			return b, nil
		})
		authRequest := model.AuthRequest{}

		if err != nil {
//			c.AbortWithError(401, err)
			fmt.Println(err)
			authRequest.IsAuthenticated = false
			authRequest.Error = err
		}else{
			authRequest.IsAuthenticated = true
			authRequest.Email,_ = token.Claims["email"].(string)
			authRequest.Username,_ = token.Claims["nickname"].(string)
//			meta,ok := token.Claims["app_metadata"].(map[string]interface{})
//			fmt.Println("app_metadata = ",reflect.TypeOf(meta["bhub_role"]))
//			if ok == true {
//				c.Set("role", meta["bhub_role"])
//			}
		}
		c.Set("AuthRequest",authRequest)
	}
}

