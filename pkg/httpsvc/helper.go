package httpsvc

// import (
// 	"errors"
// 	"time"

// 	"github.com/fahmifan/devkit/pkg/core/auth"
// 	"github.com/fahmifan/devkit/utils"
// 	"github.com/labstack/echo/v4"

// 	"github.com/dgrijalva/jwt-go"
// 	"github.com/fahmifan/devkit/config"
// 	"github.com/fahmifan/devkit/model"
// )

// // Create the JWT key used to create the signature
// var jwtKey = []byte(config.JWTKey())

// // ErrTokenInvalid error
// var ErrTokenInvalid = errors.New("token invalid")

// // Claims jwt claim
// type Claims struct {
// 	ID    string `json:"id"`
// 	Email string `json:"email"`
// 	Name  string `json:"name"`
// 	Role  string `json:"role"`
// 	jwt.StandardClaims
// }

// func createTokenExpiry() int64 {
// 	expireTime := time.Now().Add(8 * time.Hour)
// 	tokenExpiry := expireTime.UnixNano() / 1000000
// 	return tokenExpiry
// }

// func generateToken(user model.User, expiry int64) (string, error) {
// 	claims := &Claims{
// 		ID:    user.ID,
// 		Email: user.Email,
// 		Role:  user.Role.ToString(),
// 		Name:  user.Name,
// 		StandardClaims: jwt.StandardClaims{
// 			// millisecond
// 			ExpiresAt: expiry,
// 		},
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	tokenString, err := token.SignedString(jwtKey)
// 	if err != nil {
// 		return "", err
// 	}

// 	return tokenString, nil
// }

// func parseJWTToken(token string) (Claims, error) {
// 	claims := &Claims{}
// 	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
// 		return jwtKey, nil
// 	})

// 	if err != nil {
// 		if err == jwt.ErrSignatureInvalid {
// 			return *claims, err
// 		}
// 	}

// 	if tkn != nil && !tkn.Valid {
// 		return *claims, ErrTokenInvalid
// 	}

// 	return *claims, nil
// }

// func parseToken(token string) (*model.User, bool) {
// 	claims, err := parseJWTToken(token)
// 	if err != nil {
// 		return nil, false
// 	}

// 	user := &model.User{
// 		Base:  model.Base{ID: claims.ID},
// 		Email: claims.Email,
// 		Role:  auth.Role(claims.Role),
// 	}

// 	return user, true
// }

// func getUserFromCtx(c echo.Context) *model.User {
// 	res := c.Get(userInfoCtx)
// 	if val, ok := res.(model.User); ok {
// 		return &val
// 	}

// 	return nil
// }

// const userInfoCtx = "userInfoCtx"

// func setUserToCtx(c echo.Context, user *model.User) {
// 	c.Set(userInfoCtx, *user)
// }
