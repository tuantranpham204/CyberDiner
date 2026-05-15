package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tuantranpham204/CyberDiner.git/src/server/internal/model/payload/response"
	"github.com/tuantranpham204/CyberDiner.git/src/server/internal/store"
	"github.com/tuantranpham204/CyberDiner.git/src/server/pkg/util"
)

// Keys used to stash auth claims in gin.Context.
const (
	CtxUserID    = "auth.user_id"
	CtxUsername  = "auth.username"
	CtxRole      = "auth.role"
	CtxTokenID   = "auth.jti"
	CtxExpiresAt = "auth.exp"
)

// Auth verifies a Bearer JWT and rejects revoked tokens against the denylist.
// Successful requests have their claims stashed in gin.Context.
func Auth(jwt *util.JWTManager, denylist store.TokenDenylist) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.NewError(
				http.StatusUnauthorized,
				"Authorization header is missing or malformed.",
				nil,
			))
			return
		}
		tokenStr := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))

		claims, err := jwt.Verify(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.NewError(
				http.StatusUnauthorized,
				"Token is invalid or expired.",
				nil,
			))
			return
		}

		denied, err := denylist.Contains(c.Request.Context(), claims.ID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, response.NewError(
				http.StatusInternalServerError,
				"Failed to verify session.",
				nil,
			))
			return
		}
		if denied {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.NewError(
				http.StatusUnauthorized,
				"Session has been revoked. Please sign in again.",
				nil,
			))
			return
		}

		c.Set(CtxUserID, claims.UserID)
		c.Set(CtxUsername, claims.Username)
		c.Set(CtxRole, claims.Role)
		c.Set(CtxTokenID, claims.ID)
		if claims.ExpiresAt != nil {
			c.Set(CtxExpiresAt, claims.ExpiresAt.Time)
		}
		c.Next()
	}
}

// UserID returns the authenticated user's numeric ID parsed from the JWT
// subject. Returns (0, false) if the context is unauthenticated or the claim
// is malformed.
func UserID(c *gin.Context) (int64, bool) {
	v, ok := c.Get(CtxUserID)
	if !ok {
		return 0, false
	}
	s, ok := v.(string)
	if !ok {
		return 0, false
	}
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, false
	}
	return id, true
}

// TokenID returns the jti stored by the Auth middleware.
func TokenID(c *gin.Context) string {
	v, _ := c.Get(CtxTokenID)
	s, _ := v.(string)
	return s
}

// TokenExpiresAt returns the expiration time stored by the Auth middleware.
func TokenExpiresAt(c *gin.Context) time.Time {
	v, _ := c.Get(CtxExpiresAt)
	t, _ := v.(time.Time)
	return t
}
