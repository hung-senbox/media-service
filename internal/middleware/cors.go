package middleware

import "github.com/gofiber/fiber/v2"

func CORS() fiber.Handler {
	return func(c *fiber.Ctx) error {
		origin := c.Get("Origin")
		c.Set("Vary", "Origin")
		// Accept all domains:
		// - If Origin present: echo it and allow credentials
		// - If no Origin (curl/server-to-server): return "*", no credentials
		if origin != "" {
			c.Set("Access-Control-Allow-Origin", origin)
			c.Set("Access-Control-Allow-Credentials", "true")
		} else {
			c.Set("Access-Control-Allow-Origin", "*")
		}
		c.Set("Access-Control-Expose-Headers", "Content-Disposition")
		reqHeaders := c.Get("Access-Control-Request-Headers")
		if reqHeaders != "" {
			c.Set("Access-Control-Allow-Headers", reqHeaders)
		} else {
			c.Set("Access-Control-Allow-Headers",
				"Content-Type, Content-Length, Accept, Accept-Encoding, Accept-Language, X-CSRF-Token, Authorization, Cache-Control, X-Requested-With, X-App-Language, ngrok-skip-browser-warning")
		}
		reqMethod := c.Get("Access-Control-Request-Method")
		if reqMethod != "" {
			c.Set("Access-Control-Allow-Methods", reqMethod)
		} else {
			c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		}
		c.Set("Access-Control-Max-Age", "86400")

		if c.Method() == "OPTIONS" {
			return c.SendStatus(204)
		}

		return c.Next()
	}
}
