package httpmiddleware

import (
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tamboto2000/gotoko-pos/apperror"
	commonerr "github.com/tamboto2000/gotoko-pos/common/errors"
	"github.com/tamboto2000/gotoko-pos/domain/cashiers"
	"github.com/tamboto2000/gotoko-pos/helpers/httpresponse"
	"github.com/tamboto2000/gotoko-pos/helpers/jwtparse"
	"go.uber.org/zap"
)

var ErrUnauthorized = apperror.New("Unauthorized", apperror.InvalidAuth, "")

func AuthMiddleware(cr *cashiers.CashiersRepository, cmr *cashiers.CashiersMemRepository, log *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bearer := ctx.GetHeader("Authorization")
		rgx := regexp.MustCompile(`^JWT ([a-zA-Z0-9-_]+\.[a-zA-Z0-9-_]+\.[a-zA-Z0-9-_]+)$`)
		res := rgx.FindAllStringSubmatch(bearer, 1)
		if res == nil {
			log.Warn("malformed token " + bearer)

			res := httpresponse.FromError(ErrUnauthorized)
			ctx.JSON(res.StatusCode, res)
			ctx.Abort()

			return
		}

		token := res[0][1]
		claims, err := jwtparse.ParseJWT(token)
		if err != nil {
			log.Warn("authorization failed, token " + bearer)
			res := httpresponse.FromError(ErrUnauthorized)
			ctx.JSON(res.StatusCode, res)
			ctx.Abort()

			return
		}

		cashierId, err := strconv.Atoi(claims.Subject)
		if err != nil {
			log.Error(err.Error())
			res := httpresponse.FromError(commonerr.ErrInternal)
			ctx.JSON(res.StatusCode, res)
			ctx.Abort()

			return
		}

		iat := claims.IssuedAt.Unix()

		if _, err := cmr.GetSession(ctx, cashierId, iat); err != nil {
			if err.Error() == cashiers.ErrCashierSessionNotFound.Error() {
				if _, err := cr.GetSession(ctx, cashierId, iat); err != nil {
					if err.Error() == cashiers.ErrCashierSessionNotFound.Error() {
						res := httpresponse.FromError(ErrUnauthorized)
						ctx.JSON(res.StatusCode, res)
						ctx.Abort()

						return
					} else {
						log.Error(err.Error())
						res := httpresponse.FromError(err)
						ctx.JSON(res.StatusCode, res)
						ctx.Abort()

						return
					}
				}
			} else {
				log.Error(err.Error())
				res := httpresponse.FromError(err)
				ctx.JSON(res.StatusCode, res)
				ctx.Abort()

				return
			}
		}

		ctx.Set("cashier-id", cashierId)
		ctx.Next()
	}
}
