package tokens

import (
	"at-backend-claims/internal/pkg/roles"
	"context"
	"crypto/rsa"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

type tokensVerifier struct {
	jwksTTL    time.Duration
	jwksURI    string
	issuer     string
	mutex      sync.RWMutex
	keyID      string
	httpClient *http.Client

	signingKey        *rsa.PublicKey
	signingKeyFetched time.Time
}

func NewTokensVerifier(
	jwksTTL time.Duration,
	jwksURI string,
	issuer string,
	keyID string,
	httpTimeout time.Duration,
) *tokensVerifier {
	return &tokensVerifier{
		jwksTTL: jwksTTL,
		jwksURI: jwksURI,
		issuer:  issuer,
		keyID:   keyID,

		httpClient: &http.Client{Timeout: httpTimeout},
	}
}

func (tm *tokensVerifier) VerifyJWT(ctx context.Context, tokenString string) (uuid.UUID, roles.Role, error) {
	publicKey, err := tm.fetchSigningKey(ctx)
	if err != nil {
		return uuid.UUID{}, roles.Unknown, err
	}

	token, err := jwt.ParseString(tokenString,
		jwt.WithVerify(true),
		jwt.WithValidate(true),
		jwt.WithKey(jwa.RS256(), publicKey),
		jwt.WithIssuer(tm.issuer),
		jwt.WithRequiredClaim("role"),
		jwt.WithRequiredClaim("sub"),
	)
	if err != nil {
		return uuid.UUID{}, roles.Unknown, err
	}

	var roleString, sub string

	if err := token.Get("role", &roleString); err != nil {
		return uuid.UUID{}, roles.Unknown, err
	}

	if err := token.Get("sub", &sub); err != nil {
		return uuid.UUID{}, roles.Unknown, err
	}

	userID, err := uuid.Parse(sub)
	if err != nil {
		return uuid.UUID{}, roles.Unknown, err
	}

	return userID, roles.ToRole(roleString), nil
}

func (c *tokensVerifier) fetchSigningKey(ctx context.Context) (*rsa.PublicKey, error) {
	c.mutex.RLock()
	if c.signingKey != nil && time.Since(c.signingKeyFetched) < c.jwksTTL {
		defer c.mutex.RUnlock()
		return c.signingKey, nil
	}
	c.mutex.RUnlock()

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.signingKey != nil && time.Since(c.signingKeyFetched) < c.jwksTTL {
		return c.signingKey, nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.jwksURI, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("code: %d", resp.StatusCode)
	}

	set, err := jwk.ParseReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var publicKey *rsa.PublicKey

	for i := 0; i < set.Len(); i++ {
		key, ok := set.Key(i)
		if !ok {
			continue
		}

		if usage, ok := key.KeyUsage(); !ok || usage != "sig" {
			continue
		}

		var pk rsa.PublicKey

		if err := jwk.Export(key, &pk); err != nil {
			return nil, err
		}

		publicKey = &pk

		break
	}

	if publicKey == nil {
		return nil, fmt.Errorf("signing key not found")
	}

	c.signingKey = publicKey
	c.signingKeyFetched = time.Now()
	return publicKey, nil
}
