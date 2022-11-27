package middleware

import "github.com/gin-gonic/gin"

func MetaInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		if logID := c.GetHeader("LOG_ID"); logID != "" {
			c.Set("LOG_ID", logID)
		}
		c.Next()
	}
}
