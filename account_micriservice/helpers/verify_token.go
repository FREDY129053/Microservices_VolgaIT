package helpers

import (
	"time"
)

func VerifyToken(tokenStr string) (bool, error) {
	claims, err := ParseToken(tokenStr)

	if time.Until(time.Unix(claims.ExpiresAt, 0)) < 0 {
		return false, err
	}
	
	if err != nil {
		return false, err
	}

	return true, nil
}