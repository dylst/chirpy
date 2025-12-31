package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name          string
		password      string
		hash          string
		wantErr       bool
		matchPassword bool
	}{
		{
			name:          "Correct password",
			password:      password1,
			hash:          hash1,
			wantErr:       false,
			matchPassword: true,
		},
		{
			name:          "Incorrect password",
			password:      "wrongPassword",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Password doesn't match different hash",
			password:      password1,
			hash:          hash2,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Empty password",
			password:      "",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Invalid hash",
			password:      password1,
			hash:          "invalidhash",
			wantErr:       true,
			matchPassword: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && match != tt.matchPassword {
				t.Errorf("CheckPasswordHash() expects %v, got %v", tt.matchPassword, match)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	userIdOne := uuid.New()
	userIdTwo := uuid.New()
	secretKey := "WTfcT9HiKeuorBmrhFGXCJK0kgBiHsp1AzKCw6nDojE"
	oneSecJWT, err := MakeJWT(userIdOne, secretKey, time.Microsecond)
	if err != nil {
		t.Errorf("Failed to make oneSecJWT: %v", err)
	}
	fiveMinJWT, err := MakeJWT(userIdTwo, secretKey, 5 * time.Minute)
	if err != nil {
		t.Errorf("Failed to make fiveMinJWT: %v", err)
	}

	tests := []struct {
		name          string
		jwt           string
		secretKey     string
		wantErr       bool
		reject 		  bool
	}{
		{
			name:          "Valid JWT",
			jwt:           fiveMinJWT,
			secretKey:	   secretKey,
			wantErr:       false,
			reject: 	   false,
		},
		{
			name:          "Expired JWT",
			jwt:           oneSecJWT,
			secretKey:     secretKey,
			wantErr:       true,
			reject: 	   true,
		},
		{
			name:          "Invalid secret key",
			jwt:           fiveMinJWT,
			secretKey: 	   "random",
			wantErr:       true,
			reject: 	   true,
		},
	}

	for _, tt := range tests{
		t.Run(tt.name, func(t *testing.T) {
			valid, err := ValidateJWT(tt.jwt, tt.secretKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() Test %v, error: %v, wantErr %v", tt.name, err, tt.wantErr)
			}

			if !tt.wantErr && (valid != userIdOne && valid != userIdTwo) {
				t.Errorf("ValidateJWT() Test %v, expects %v or %v, got %v", tt.name, userIdOne, userIdTwo, valid)
			}
		})
	}
}