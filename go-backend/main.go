package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/lestrrat-go/jwx/jwk"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:4200",
			"https://localhost:4200",
		},
	}))

	jwksCertUrl := "http://localhost:8080/auth/realms/default-realm/protocol/openid-connect/certs";

	ar := jwk.NewAutoRefresh(context.Background())

	// Tell *jwk.AutoRefresh that we only want to refresh this JWKS
	// when it needs to (based on Cache-Control or Expires header from
	// the HTTP response). If the calculated minimum refresh interval is less
	// than 15 minutes, don't go refreshing any earlier than 15 minutes.
	ar.Configure(jwksCertUrl, jwk.WithMinRefreshInterval(15*time.Minute))

	// Refresh the JWKS once before getting into the main loop.
	// This allows you to check if the JWKS is available before we start
	// a long-running program
	_, err := ar.Refresh(context.Background(), jwksCertUrl)
	if err != nil {
		log.Fatalln("Failed to get JWKS")
	}

	e.GET("/protected", func(context echo.Context) error {
		return context.JSON(http.StatusOK, struct {
			Message string
		}{Message: "oke"})
	}, middleware.JWTWithConfig(middleware.JWTConfig{KeyFunc: generateGetKeyFunc(ar)}))

	if err := e.Start(":3000"); err != nil {
		fmt.Println(err)
	}
}

func generateGetKeyFunc(autoRefresher *jwk.AutoRefresh) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		// For a demonstration purpose, Google Sign-in is used.
		// https://developers.google.com/identity/sign-in/web/backend-auth
		//
		// This user-defined KeyFunc verifies tokens issued by Google Sign-In.
		//
		// Note: In this example, it downloads the keyset every time the restricted route is accessed.
		keySet, err := autoRefresher.Fetch(context.Background(), "http://localhost:8080/auth/realms/default-realm/protocol/openid-connect/certs")
		if err != nil {
			return nil, err
		}

		keyID, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("expecting JWT header to have a key ID in the kid field")
		}

		key, found := keySet.LookupKeyID(keyID)

		if !found {
			return nil, fmt.Errorf("unable to find key %q", keyID)
		}

		var pubkey interface{}
		if err := key.Raw(&pubkey); err != nil {
			return nil, fmt.Errorf("Unable to get the public key. Error: %s", err.Error())
		}

		return pubkey, nil
	}
}
