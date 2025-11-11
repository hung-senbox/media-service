package middleware

import "github.com/gin-gonic/gin"

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		c.Writer.Header().Set("Vary", "Origin")
		// allow-list specific origins
		allowedOrigins := map[string]bool{
			"https://main.d3snrh9uwyhkta.amplifyapp.com": true,
		}
		if origin != "" && allowedOrigins[origin] {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		} else if origin == "" {
			// non-browser or curl; allow all without credentials
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}

		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
		reqHeaders := c.Request.Header.Get("Access-Control-Request-Headers")
		if reqHeaders != "" {
			c.Writer.Header().Set("Access-Control-Allow-Headers", reqHeaders)
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Headers",
				"Content-Type, Content-Length, Accept, Accept-Encoding, Accept-Language, X-CSRF-Token, Authorization, Cache-Control, X-Requested-With, X-App-Language, ngrok-skip-browser-warning")
		}
		reqMethod := c.Request.Header.Get("Access-Control-Request-Method")
		if reqMethod != "" {
			c.Writer.Header().Set("Access-Control-Allow-Methods", reqMethod)
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		}
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
