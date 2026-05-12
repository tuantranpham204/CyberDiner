package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tuantranpham204/CyberDiner.git/src/server/internal/store"
	"github.com/tuantranpham204/CyberDiner.git/src/server/pkg/util"
)

// Context keys populated by JWTAuth for downstream handlers.
const (
	ContextUserID = "user_id"
	ContextRole   = "role"
	ContextJTI    = "jti"
	ContextExpiry = "expiry"
)

// JWTAuth validates the Bearer token, checks the denylist, and populates context.
// All JWT concerns stop here — handlers read plain values from context.
func JWTAuth(denylist store.DenylistStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "missing token"})
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")

		userID, role, jti, expiry, err := util.ValidateAccessToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "invalid or expired token"})
			return
		}

		if denylist.Contains(jti) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "token has been revoked"})
			return
		}

		c.Set(ContextUserID, userID)
		c.Set(ContextRole, role)
		c.Set(ContextJTI, jti)
		c.Set(ContextExpiry, expiry)
		c.Next()
	}
}

// GetUserID reads the authenticated user's ID from context.
func GetUserID(c *gin.Context) uint {
	v, _ := c.Get(ContextUserID)
	id, _ := v.(uint)
	return id
}

// GetJTI reads the token's JTI from context.
func GetJTI(c *gin.Context) (string, time.Time) {
	jti, _ := c.Get(ContextJTI)
	exp, _ := c.Get(ContextExpiry)
	jtiStr, _ := jti.(string)
	expTime, _ := exp.(time.Time)
	return jtiStr, expTime
}
