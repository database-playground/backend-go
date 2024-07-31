package gatewayservice

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/database-playground/backend/internal/services/gateway/openapi"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
)

type AuthContextKey string

const AuthContextJwtToken = AuthContextKey("jwt_token")

var scopeMap map[string][]string = map[string][]string{
	"PostChallenges":         {"challenge"},
	"GetChallengesId":        {"challenge"},
	"GetQuestions":           {"read:question"},
	"GetQuestionsId":         {"read:question"},
	"GetQuestionsIdSolution": {"read:question", "read:solution"},
	"GetSchemasId":           {"read:schema"},

	"GetHealthz": nil,
}

// NewAuthorizationMiddleware creates a new authorization middleware that verifies the JWT token in the authorization header.
//
// logtoDomain is the domain of the Logto instance. It is used to fetch the JWKS and for verifing
// the issuers; resourceIndicator is the resource indicator of the request listener, which is a URI.
func NewAuthorizationMiddleware(ctx context.Context, logtoDomain string, resourceIndicator string, logger *slog.Logger) nethttp.StrictHTTPMiddlewareFunc {
	logtoOidc, err := url.JoinPath(logtoDomain, "oidc")
	if err != nil {
		panic(err)
	}
	logtoJwks, err := url.JoinPath(logtoOidc, "jwks")
	if err != nil {
		panic(err)
	}

	cache := jwk.NewCache(ctx)
	if err := cache.Register(logtoJwks); err != nil {
		panic(err)
	}

	return func(f nethttp.StrictHTTPHandlerFunc, operationID string) nethttp.StrictHTTPHandlerFunc {
		scopeRequired, ok := scopeMap[operationID]
		if !ok {
			slog.Error("The operation is not defined in the scope map. Denying any access.", slog.String("operationID", operationID))
			return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (response interface{}, err error) {
				w.WriteHeader(http.StatusNotFound)
				return nil, nil
			}
		}

		if len(scopeRequired) == 0 {
			return f
		}

		return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (response interface{}, err error) {
			authorizationHeader := r.Header.Get("authorization")
			if authorizationHeader == "" {
				return sendUnauthorizedError(w, "Malformed or missing authorization header")
			}

			bearerToken, ok := strings.CutPrefix(authorizationHeader, "Bearer ")
			if !ok {
				return sendUnauthorizedError(w, "Malformed or missing authorization header")
			}

			keySet, err := cache.Get(ctx, logtoJwks)
			if err != nil {
				slog.Error("failed to fetch JWKS", slog.Any("error", err))
				return sendServerError(w, "Failed to fetch JWKS")
			}

			tok, err := jwt.ParseString(bearerToken, jwt.WithKeySet(keySet), jwt.WithValidate(true), jwt.WithIssuer(logtoOidc), jwt.WithAudience(resourceIndicator), jwt.WithContext(ctx))
			if err != nil {
				slog.Error("failed to verify token", slog.Any("error", err), slog.String("endpoint", r.URL.Path))
				return sendUnauthorizedError(w, "Failed to verify token")
			}

			userScopes := parseScope(tok.PrivateClaims()["scope"])

			for _, requiredScope := range scopeRequired {
				if !slices.Contains(userScopes, requiredScope) {
					slog.Debug("insufficient scope", slog.String("required", requiredScope), slog.Any("actual", userScopes))
					return sendUnauthorizedError(w, "Insufficient scope ("+requiredScope+" is required)")
				}
			}

			ctx = context.WithValue(ctx, AuthContextJwtToken, tok)
			return f(ctx, w, r, request)
		}
	}
}

func sendUnauthorizedError(w http.ResponseWriter, message string) (any, error) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Header().Set("Content-Type", "application/json")

	marshaled, err := json.Marshal(openapi.UnauthorizedErrorJSONResponse{
		Message: message,
	})
	if err != nil {
		return nil, err
	}

	_, err = w.Write(marshaled)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func sendServerError(w http.ResponseWriter, message string) (any, error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json")

	marshaled, err := json.Marshal(openapi.Error{
		Message: message,
	})
	if err != nil {
		return nil, err
	}

	_, err = w.Write(marshaled)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// parseScope parses the "scope" Get result of [jwt.Token]
// into a slice of strings.
func parseScope(scope any) []string {
	if scope == nil {
		return nil
	}

	raw, ok := scope.(string)
	if !ok {
		return nil
	}

	return strings.Split(raw, " ")
}
