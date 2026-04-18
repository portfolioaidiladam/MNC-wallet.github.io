package middleware

import (
	"strings"

	"github.com/aidiladam/mnc-wallet/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ContextKeyUserID adalah key Gin context untuk UUID user yang ter-autentikasi.
const ContextKeyUserID = "userID"

// JWTAuth memvalidasi header "Authorization: Bearer <token>". Kalau valid,
// set user UUID di context; kalau tidak, abort dengan 401.
func JWTAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" {
			util.Error(c, 401, "missing authorization header")
			c.Abort()
			return
		}
		parts := strings.SplitN(h, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
			util.Error(c, 401, "invalid authorization header")
			c.Abort()
			return
		}
		claims, err := util.ParseAccessToken(jwtSecret, parts[1])
		if err != nil {
			util.Error(c, 401, "invalid or expired token")
			c.Abort()
			return
		}
		c.Set(ContextKeyUserID, claims.UserID)
		c.Next()
	}
}

// UserIDFromContext mengambil UUID user yang sudah di-set JWTAuth.
// Return zero UUID + false kalau context tidak punya (mis. dipanggil di
// endpoint non-authenticated).
func UserIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	v, ok := c.Get(ContextKeyUserID)
	if !ok {
		return uuid.Nil, false
	}
	id, ok := v.(uuid.UUID)
	return id, ok
}