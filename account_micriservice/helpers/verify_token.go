package helpers

func VerifyToken(tokenStr string) (bool, error) {
	_, err := ParseToken(tokenStr)

	if err != nil {
		return false, err
	}

	return true, nil
}