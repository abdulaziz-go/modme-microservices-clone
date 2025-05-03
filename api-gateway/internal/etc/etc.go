package etc

import (
	client "api-gateway/internal/clients"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"google.golang.org/grpc/metadata"
	"net/http"
	"strings"
	"time"
)

func AuthMiddleware(requiredRoles []string, userClient *client.UserClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			ctx.Abort()
			return
		}

		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			ctx.Abort()
			return
		}
		token := strings.TrimPrefix(authHeader, bearerPrefix)

		user, err := userClient.ValidateToken(token, requiredRoles)

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Invalid or insufficient permissions required role => %v", requiredRoles)})
			ctx.Abort()
			return
		}
		ctx.Set("user", user)
		ctx.Set("company_id", cast.ToString(user.CompanyId))
		ctx.Next()
	}
}

func NewTimoutContext(ctx context.Context) (context.Context, context.CancelFunc) {
	md := metadata.Pairs()
	for _, key := range []string{"company_id"} {
		if ctx.Value(key) != nil {
			val, ok := ctx.Value(key).(string)
			if ok {
				md.Set(key, val)
			}
		}
	}
	ctx = metadata.NewOutgoingContext(ctx, md)
	res, cancel := context.WithTimeout(ctx, time.Second*15)
	return res, cancel
}

func NewTimeoutContextManual(ctx context.Context, companyId string) (context.Context, context.CancelFunc) {
	md := metadata.Pairs()
	md.Set("company_id", companyId)
	ctx = metadata.NewOutgoingContext(ctx, md)
	return context.WithTimeout(ctx, time.Second*15)
}
