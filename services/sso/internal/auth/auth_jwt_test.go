package auth

import "testing"

func TestCreatingAndValidatingToken(t *testing.T) {
	id := 5
	token, err := CreateToken(id)
	if err != nil {
		t.Error(err)
	}

	claims, err := ValidateToken(token)
	if err != nil {
		t.Error(err)
	}

	if claims == nil {
		t.Fatal("expected not nill claims")
	}

	if claims.Id != id {
		t.Fatalf("expected id: %d got: %d", id, claims.Id)
	}
}
