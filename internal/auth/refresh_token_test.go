package auth

import (
	"fmt"
	"testing"
)

func TestRefreshToken(t *testing.T) {
	a := MakeRefreshToken()
	b := MakeRefreshToken()

	passCount := 0

	if len(a) != 12 {
		fmt.Println("=============================")
		fmt.Println("Invalid token length")
		t.Fail()
		fmt.Println("=============================")
	} else {
		passCount += 1
	}

	if a == b {
		fmt.Println("same string can't be generated twice")
		t.Fail()
		fmt.Println("=============================")
	} else {
		passCount += 1
	}

	fmt.Printf("Make Refresh Token : %v PASSED, %v FAILED\n", passCount, 2-passCount)
	fmt.Println("=============================")
}
