package dpsauth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const cookieExpiresInHours = 1
const prefix = "mymove-"

// UserIDToCookie takes the UUID of the current user and returns the cookie value.
func UserIDToCookie(userID string) (string, error) {
	expirationTime := time.Now().Add(time.Hour * time.Duration(cookieExpiresInHours)).Unix()
	value := map[string]string{
		"user_id":    userID,
		"expires_at": strconv.FormatInt(expirationTime, 10),
	}

	valueJSON, err := json.Marshal(value)
	if err != nil {
		return "", errors.Wrap(err, "Marshaling the cookie JSON")
	}

	encrypted, err := encrypt(valueJSON)
	if err != nil {
		return "", errors.Wrap(err, "Encrypting the cookie")
	}

	return prefix + encrypted, nil
}

// CookieToUserID takes a cookie value and returns the user's UUID only if it's a
// valid, unexpired cookie.
func CookieToUserID(token string) (string, error) {
	if !strings.HasPrefix(token, prefix) {
		return "", errors.New("Invalid cookie: missing prefix")
	}

	decryptedToken, err := decrypt(token[len(prefix):])
	if err != nil {
		return "", errors.Wrap(err, "Decrypting the cookie")
	}

	var values map[string]string
	err = json.Unmarshal(decryptedToken, &values)
	if err != nil {
		return "", errors.Wrap(err, "Unmarshaling the cookie JSON")
	}

	// TODO: check that the cookie is not expired

	return values["user_id"], nil
}

func encrypt(data []byte) (string, error) {
	key := getKey()
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", errors.Wrap(err, "Creating a new cipher using the key")
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", errors.Wrap(err, "NewGCM call")
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	encoded := gcm.Seal(nonce, nonce, data, nil)

	// Use base64 URL encoding since this will be passed back as an API param
	return base64.RawURLEncoding.EncodeToString(encoded), nil
}

func decrypt(data string) ([]byte, error) {
	key := getKey()
	var plaintext []byte
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return plaintext, err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return plaintext, err
	}
	dataBytes, err := base64.RawURLEncoding.DecodeString(data)
	if err != nil {
		return plaintext, err
	}
	nonceSize := gcm.NonceSize()
	nonce, cipher := dataBytes[:nonceSize], dataBytes[nonceSize:]
	plaintext, err = gcm.Open(nil, nonce, cipher, nil)
	if err != nil {
		return plaintext, err
	}
	return plaintext, nil
}

func getKey() string {
	return os.Getenv("DPS_AUTH_COOKIE_SECRET_KEY")
}
