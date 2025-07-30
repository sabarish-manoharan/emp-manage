package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sabarish-manoharan/emp-management/controllers"
)

func AuthMiddleware(next http.Handler)http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
 
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader,"Bearer "){
			http.Error(w,"Missing or invalid Authorization header",http.StatusUnauthorized);
			return ;
		}

		tokenString := strings.TrimPrefix(authHeader,"Bearer ");

		token,err := jwt.Parse(tokenString,func(t *jwt.Token) (interface{}, error) {
			if _,ok := t.Method.(*jwt.SigningMethodHMAC);!ok{
				return  nil,http.ErrAbortHandler;
			}
			return  controllers.JwtSecret,nil;
		});
		if err!=nil || !token.Valid{
          http.Error(w,"Invalid Token",http.StatusUnauthorized);
		  return 
		}
		next.ServeHTTP(w,r);
	});
}