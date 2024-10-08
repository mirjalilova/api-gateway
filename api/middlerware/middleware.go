package middlerware

import (
	"golang.org/x/exp/slog"
	"net/http"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	token "github.com/mirjalilova/api-gateway/api/token"
)

const (
	key = "secret_key"
)

func NewAuth(enforce *casbin.Enforcer) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		allow, err := CheckPermission(ctx.FullPath(), ctx.Request, enforce)

		if err != nil {
			valid, _ := err.(jwt.ValidationError)
			if valid.Errors == jwt.ValidationErrorExpired {
				RequireRefresh(ctx)
			} else {
				RequirePermission(ctx)
			}
		} else if !allow {
			RequirePermission(ctx)
		}
		ctx.Next()
	}

}

func GetRole(r *http.Request) (string, error) {
	var (
		claims jwt.MapClaims
		err    error
	)

	jwtToken := r.Header.Get("Authorization")
	if jwtToken == "" {
		return "unauthorized", nil
	} else if strings.Contains(jwtToken, "Basic") {
		return "unauthorized", nil
	}
	claims, err = token.ExtractClaims(jwtToken, key)

	if err != nil {
		slog.Error("Error while extracting claims: ", err)
		return "unauthorized", err
	}

	return claims["role"].(string), nil
}

func GetUserIdAndUsername(r *http.Request) (string, string, error) {
	var (
		claims jwt.MapClaims
		err    error
	)

	jwtToken := r.Header.Get("Authorization")
	if jwtToken == "" {
		return "unauthorized", "", nil
	} else if strings.Contains(jwtToken, "Basic") {
		return "unauthorized", "", nil
	}
	claims, err = token.ExtractClaims(jwtToken, key)

	if err != nil {
		slog.Error("Error while extracting claims: ", err)
		return "unauthorized", "", err
	}

	return claims["userid"].(string), claims["username"].(string), nil
}

func CheckPermission(path string, r *http.Request, enforcer *casbin.Enforcer) (bool, error) {
	role, err := GetRole(r)
	if err != nil {
		slog.Error("Error while getting role from token: ", err)
		return false, err
	}
	method := r.Method
	// path := r.URL.Path

	allowed, err := enforcer.Enforce(role, path, method)
	if err != nil {
		slog.Error("Error while enforcing policy: ", err)
		return false, err
	}

	return allowed, nil
}

func InvalidToken(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		"error": "Invalid token !!!",
	})
}

// RequirePermission handles responses for insufficient permissions
func RequirePermission(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		"error": "Permission denied",
	})
}

// RequireRefresh handles responses for expired access tokens
func RequireRefresh(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"error": "Access token expired",
	})
}
