package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

var jwt_token string
var user_id uuid.UUID

func TestMakeJWT(t *testing.T) {
	test_Userid := uuid.New()
	user_id = test_Userid
	test_jwtSecret := "test-jwt-secret"

	passCount := 0

	jwt, err := MakeJWT(test_Userid, test_jwtSecret)
	if err != nil {
		fmt.Println("=====================")
		fmt.Printf("JWT ERROR: %v\n", err)
		t.Fail()
		fmt.Println("=====================")
	} else {
		jwt_token = jwt
		passCount += 1
	}

	fmt.Printf("Make JWT : %v PASSED, %v FAILED\n", passCount, 1-passCount)
	fmt.Println("=============================")
}

func TestValidateJWT(t *testing.T) {
	test_jwtSecret := "test-jwt-secret"
	failCount := 0

	token, claims, err := ValidateJWT(jwt_token, test_jwtSecret)
	if err != nil {
		fmt.Println("=====================")
		fmt.Printf("%v\n", err)
		t.Fail()
		fmt.Println("=====================")
		failCount += 1
	}

	if !token.Valid {
		fmt.Println("=====================")
		fmt.Printf("Invalid JWT token\n")
		t.Fail()
		fmt.Println("=====================")
		failCount += 1
	}

	id, _ := uuid.Parse(claims.ID)

	if id != user_id {
		fmt.Println("=====================")
		fmt.Printf("Invalid user id\n")
		t.Fail()
		fmt.Println("=====================")
		failCount += 1
	}

	if time.Now().UTC().After(time.Unix(claims.ExpiresAt.Unix(), 0)) {
		fmt.Println("=====================")
		fmt.Printf("Expired JWT token\n")
		t.Fail()
		fmt.Println("=====================")
		failCount += 1
	}

	fmt.Printf("Validate JWT: %v PASSED, %v FAILED\n", 4-failCount, failCount)
	fmt.Println("=============================")
}
