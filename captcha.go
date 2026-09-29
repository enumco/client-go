package client

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	altcha "github.com/altcha-org/altcha-lib-go/v2"
)

// SolveCaptcha answers the challenge GetCaptchaChallenge returned and encodes it
// the way RegisterRequest.captcha_token expects. It is proof of work, so it costs
// the caller a few seconds of CPU and needs no browser.
func SolveCaptcha(challenge string) (string, error) {
	var parsed altcha.Challenge
	if err := json.Unmarshal([]byte(challenge), &parsed); err != nil {
		return "", fmt.Errorf("the captcha challenge is not valid json: %w", err)
	}

	solution, err := altcha.SolveChallenge(altcha.SolveChallengeOptions{
		Challenge: parsed,
		DeriveKey: altcha.DeriveKeyPBKDF2(),
	})
	if err != nil {
		return "", err
	}

	if solution == nil {
		return "", errors.New("the captcha challenge has no solution")
	}

	payload, err := json.Marshal(altcha.Payload{Challenge: parsed, Solution: *solution})
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(payload), nil
}
