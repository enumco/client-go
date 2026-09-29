package client

import (
	"encoding/base64"
	"encoding/json"
	"testing"
	"time"

	altcha "github.com/altcha-org/altcha-lib-go/v2"
)

// The payload has to satisfy the same verification enumd runs, so this issues a
// challenge the way enumd issues one and checks the answer the way enumd checks
// it. Cost is low only to keep the test quick; it does not affect the result.
func TestSolveCaptchaProducesAVerifiablePayload(t *testing.T) {
	const secret = "captcha-secret"

	expires := time.Now().Add(time.Minute)

	challenge, err := altcha.CreateChallenge(altcha.CreateChallengeOptions{
		Algorithm:           "PBKDF2/SHA-256",
		DeriveKey:           altcha.DeriveKeyPBKDF2(),
		Cost:                100,
		ExpiresAt:           &expires,
		HMACSignatureSecret: secret,
	})
	if err != nil {
		t.Fatalf("CreateChallenge: %v", err)
	}

	issued, err := json.Marshal(challenge)
	if err != nil {
		t.Fatal(err)
	}

	token, err := SolveCaptcha(string(issued))
	if err != nil {
		t.Fatalf("SolveCaptcha: %v", err)
	}

	raw, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		t.Fatalf("the token is not base64: %v", err)
	}

	var payload altcha.Payload
	if err := json.Unmarshal(raw, &payload); err != nil {
		t.Fatalf("the token is not a payload: %v", err)
	}

	result, err := altcha.VerifySolution(altcha.VerifySolutionOptions{
		Challenge:           payload.Challenge,
		Solution:            payload.Solution,
		DeriveKey:           altcha.DeriveKeyPBKDF2(),
		HMACSignatureSecret: secret,
	})
	if err != nil {
		t.Fatalf("VerifySolution: %v", err)
	}

	if result.Expired {
		t.Error("the solution reads as expired")
	}

	if !result.Verified {
		t.Error("the solution does not verify")
	}
}

func TestSolveCaptchaRejectsGarbage(t *testing.T) {
	if _, err := SolveCaptcha("not json"); err == nil {
		t.Error("SolveCaptcha accepted a challenge that is not json")
	}
}
