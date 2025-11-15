package middleware

import (
	"context"
	"media-service/helper"
	"media-service/internal/gateway"
	"media-service/pkg/constants"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Secured(userGw gateway.UserGateway) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authorizationHeader := c.Get("Authorization")

		// app language header
		appLanguage := helper.ParseAppLanguage(c.Get("X-App-Language"), 1)
		c.Set("X-App-Language", strconv.Itoa(int(appLanguage)))
		c.Locals(constants.AppLanguage.String(), appLanguage)
		ctx := context.WithValue(c.Context(), constants.AppLanguage, appLanguage)

		// --- Authorization ---
		if len(authorizationHeader) == 0 {
			return c.SendStatus(http.StatusForbidden)
		}

		if !strings.HasPrefix(authorizationHeader, "Bearer ") {
			return c.SendStatus(http.StatusUnauthorized)
		}

		tokenString := strings.Split(authorizationHeader, " ")[1]

		token, _, _ := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// --- UserID ---
			if userId, ok := claims[constants.UserID.String()].(string); ok {
				c.Locals(constants.UserID.String(), userId)
				ctx = context.WithValue(ctx, constants.UserID, userId)
			}

			// --- UserName ---
			if userName, ok := claims[constants.UserName.String()].(string); ok {
				c.Locals(constants.UserName.String(), userName)
				ctx = context.WithValue(ctx, constants.UserName, userName)
			}

			// --- Roles ---
			if userRoles, ok := claims[constants.UserRoles.String()].(string); ok {
				c.Locals(constants.UserRoles.String(), userRoles)
				ctx = context.WithValue(ctx, constants.UserRoles, userRoles)
			}
		}

		// --- Token ---
		c.Locals(constants.Token.String(), tokenString)
		ctx = context.WithValue(ctx, constants.Token, tokenString)

		// --- Call user-service to get current user ---
		currentUser, err := userGw.GetCurrentUser(
			context.WithValue(ctx, constants.CurrentUserKey, tokenString),
		)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized",
			})
		}

		// --- Set currentUser vào context ---
		c.Locals(string(constants.CurrentUserKey), currentUser)
		ctx = context.WithValue(ctx, constants.CurrentUserKey, currentUser)
		
		// Store the context for use in handlers
		c.SetUserContext(ctx)

		return c.Next()
	}
}

func RequireAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		rolesAny := c.Locals(constants.UserRoles.String())
		if rolesAny == nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Roles not found"})
		}

		rolesStr, ok := rolesAny.(string)
		if !ok {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid roles format"})
		}

		// ví dụ roles: "SuperAdmin, Teacher"
		roles := strings.Split(rolesStr, ",")
		isAdmin := false
		for _, role := range roles {
			if strings.TrimSpace(role) == "SuperAdmin" {
				isAdmin = true
				break
			}
		}

		if !isAdmin {
			return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "Admin access required"})
		}

		return c.Next()
	}
}
