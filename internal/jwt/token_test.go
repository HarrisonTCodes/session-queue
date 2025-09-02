package jwt

import "testing"

func TestCreateAndValidateToken(t *testing.T) {
	secret := []byte("secret")
	pos := int64(10)

	token, err := CreateToken(pos, secret)
	if err != nil {
		t.Fatalf("CreateToken failed: %v", err)
	}

	validatedPos, err := ValidateToken(token, secret)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	if validatedPos != pos {
		t.Fatalf("Expected validated token position %v, got %v", pos, validatedPos)
	}
}

func TestValidateToken_Invalid(t *testing.T) {
	createSecret := []byte("wrong-secret")
	validateSecret := []byte("secret")

	token, err := CreateToken(10, createSecret)
	if err != nil {
		t.Fatalf("CreateToken failed: %v", err)
	}

	_, err = ValidateToken(token, validateSecret)
	if err == nil {
		t.Fatalf(("Expected validate token to fail, no failure"))
	}
}
