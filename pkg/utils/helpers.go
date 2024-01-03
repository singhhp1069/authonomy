package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/TBD54566975/ssi-sdk/credential"
	"github.com/TBD54566975/ssi-sdk/crypto/jwx"
	"github.com/lestrrat-go/jwx/v2/jws"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/pkg/errors"
)

const (
	VCJWTProperty string = "vc"
)

type HTTPResponse struct {
	Body       interface{}
	StatusCode int
}

func SendHTTPRequest(method, url string, body *bytes.Buffer) (*HTTPResponse, error) {
	var req *http.Request
	var err error

	if body != nil {
		req, err = http.NewRequest(method, url, body)
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var responseJSON interface{}
	err = json.Unmarshal(responseBody, &responseJSON)
	if err != nil {
		return nil, err
	}

	return &HTTPResponse{
		Body:       responseJSON,
		StatusCode: resp.StatusCode,
	}, nil
}

func ParseVerifiableCredentialFromJWT(token string) (jws.Headers, jwt.Token, *credential.VerifiableCredential, error) {
	parsed, err := jwt.Parse([]byte(token), jwt.WithValidate(false), jwt.WithVerify(false))
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "parsing credential token")
	}

	// get headers
	headers, err := jwx.GetJWSHeaders([]byte(token))
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "getting JWT headers")
	}

	// parse remaining JWT properties and set in the credential
	cred, err := ParseVerifiableCredentialFromToken(parsed)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "parsing credential from token")
	}

	return headers, parsed, cred, nil
}

// ParseVerifiableCredentialFromToken takes a JWT object and parses it into a VerifiableCredential
func ParseVerifiableCredentialFromToken(token jwt.Token) (*credential.VerifiableCredential, error) {
	// parse remaining JWT properties and set in the credential
	vcClaim, ok := token.Get(VCJWTProperty)
	if !ok {
		return nil, fmt.Errorf("did not find %s property in token", VCJWTProperty)
	}
	vcBytes, err := json.Marshal(vcClaim)
	if err != nil {
		return nil, errors.Wrap(err, "marshalling credential claim")
	}
	var cred credential.VerifiableCredential
	if err = json.Unmarshal(vcBytes, &cred); err != nil {
		return nil, errors.Wrap(err, "reconstructing Verifiable Credential")
	}

	jti, hasJTI := token.Get(jwt.JwtIDKey)
	jtiStr, ok := jti.(string)
	if hasJTI && ok && jtiStr != "" {
		cred.ID = jtiStr
	}

	iat, hasIAT := token.Get(jwt.IssuedAtKey)
	iatTime, ok := iat.(time.Time)
	if hasIAT && ok {
		cred.IssuanceDate = iatTime.Format(time.RFC3339)
	}

	exp, hasExp := token.Get(jwt.ExpirationKey)
	expTime, ok := exp.(time.Time)
	if hasExp && ok {
		cred.ExpirationDate = expTime.Format(time.RFC3339)
	}

	// Note: we only handle string issuer values, not objects for JWTs
	iss, hasIss := token.Get(jwt.IssuerKey)
	issStr, ok := iss.(string)
	if hasIss && ok && issStr != "" {
		cred.Issuer = issStr
	}

	sub, hasSub := token.Get(jwt.SubjectKey)
	subStr, ok := sub.(string)
	if hasSub && ok && subStr != "" {
		if cred.CredentialSubject == nil {
			cred.CredentialSubject = make(map[string]any)
		}
		cred.CredentialSubject[credential.VerifiableCredentialIDProperty] = subStr
	}

	return &cred, nil
}
