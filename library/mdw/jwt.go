package mdw

import "github.com/griffin702/service/jwt-iris"

/**
 * jwt中间件
 */
func JwtHandler() *jwt.Middleware {
	return jwt.New(jwt.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("secret123123"), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
}
