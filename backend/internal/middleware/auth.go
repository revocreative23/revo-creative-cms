package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const (
	SessionKeyUserID = "user_id"
	SessionKeyEmail  = "email"
	SessionKeyRole   = "role"
)

// RequireAuth memblok request yang tidak punya session valid.
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		s := sessions.Default(c)
		uid := s.Get(SessionKeyUserID)
		if uid == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.Set("user_id", uid)
		c.Set("email", s.Get(SessionKeyEmail))
		c.Set("role", s.Get(SessionKeyRole))
		c.Next()
	}
}
