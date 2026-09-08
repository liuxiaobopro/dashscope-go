// Copyright (c) Alibaba, Inc. and its affiliates.

package apientities

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/liuxiaobopro/dashscope-go/common"
)

// Encryption request/response AES-GCM + RSA envelope encryption.
type Encryption struct {
	PubKeyID           string
	PubKeyStr          string
	AESKeyBytes        []byte
	EncryptedAESKeyStr string
	IVBytes            []byte
	Base64IVStr        string
	Valid              bool
}

func (e *Encryption) IsValid() bool { return e != nil && e.Valid }

func (e *Encryption) GetPubKeyID() string { return e.PubKeyID }

func (e *Encryption) GetEncryptedAESKeyStr() string { return e.EncryptedAESKeyStr }

func (e *Encryption) GetBase64IVStr() string { return e.Base64IVStr }

// Initialize fetch public key and generate AES key.
func (e *Encryption) Initialize() {
	publicKeys := getPublicKeys()
	if publicKeys == nil {
		return
	}
	publicKeyStr, _ := publicKeys["public_key"].(string)
	publicKeyID, _ := publicKeys["public_key_id"].(string)
	if publicKeyStr == "" || publicKeyID == "" {
		common.Log.Error("public keys data not valid")
		return
	}
	aesKey := make([]byte, common.ENCRYPTION_AES_SECRET_KEY_BYTES)
	if _, err := io.ReadFull(rand.Reader, aesKey); err != nil {
		return
	}
	iv := make([]byte, common.ENCRYPTION_AES_IV_LENGTH)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return
	}
	encAES, err := encryptAESKeyWithRSA(aesKey, publicKeyStr)
	if err != nil {
		common.Log.Error("encrypt aes key failed: %v", err)
		return
	}
	e.PubKeyID = publicKeyID
	e.PubKeyStr = publicKeyStr
	e.AESKeyBytes = aesKey
	e.EncryptedAESKeyStr = encAES
	e.IVBytes = iv
	e.Base64IVStr = base64.StdEncoding.EncodeToString(iv)
	e.Valid = true
}

func (e *Encryption) Encrypt(dictPlaintext any) any {
	b, _ := json.Marshal(dictPlaintext)
	return encryptTextWithAES(string(b), e.AESKeyBytes, e.IVBytes)
}

func (e *Encryption) Decrypt(base64Ciphertext any) any {
	s, ok := base64Ciphertext.(string)
	if !ok {
		return base64Ciphertext
	}
	return decryptTextWithAES(s, e.AESKeyBytes, e.IVBytes)
}

func getPublicKeys() map[string]any {
	url := common.BaseHTTPAPIURL + "/public-keys/latest"
	apiKey, err := common.GetDefaultAPIKey()
	if err != nil {
		common.Log.Error("exceptional public key request: %v", err)
		return nil
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	for k, v := range common.GetSDKHeaders("utils") {
		req.Header.Set(k, v)
	}
	client := &http.Client{Timeout: time.Duration(common.DEFAULT_REQUEST_TIMEOUT_SECONDS) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		common.Log.Error("exceptional public key response: %v", err)
		return nil
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		common.Log.Error("exceptional public key response: %s", resp.Status)
		return nil
	}
	var jsonResp map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&jsonResp); err != nil {
		return nil
	}
	data, _ := jsonResp["data"].(map[string]any)
	if data == nil {
		common.Log.Error("no valid data in public key response")
		return nil
	}
	return data
}

func encryptTextWithAES(plaintext string, key, iv []byte) string {
	block, err := aes.NewCipher(key)
	if err != nil {
		return ""
	}
	gcm, err := cipher.NewGCMWithNonceSize(block, len(iv))
	if err != nil {
		return ""
	}
	ciphertext := gcm.Seal(nil, iv, []byte(plaintext), []byte{})
	return base64.StdEncoding.EncodeToString(ciphertext)
}

func decryptTextWithAES(b64 string, key, iv []byte) any {
	encrypted, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}
	gcm, err := cipher.NewGCMWithNonceSize(block, len(iv))
	if err != nil {
		return nil
	}
	plain, err := gcm.Open(nil, iv, encrypted, []byte{})
	if err != nil {
		return nil
	}
	var out any
	if err := json.Unmarshal(plain, &out); err != nil {
		return string(plain)
	}
	return out
}

func encryptAESKeyWithRSA(aesKey []byte, publicKeyStr string) (string, error) {
	der, err := base64.StdEncoding.DecodeString(publicKeyStr)
	if err != nil {
		return "", err
	}
	pub, err := x509.ParsePKIXPublicKey(der)
	if err != nil {
		return "", err
	}
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return "", common.NewInvalidParameter("not an RSA public key")
	}
	b64aes := base64.StdEncoding.EncodeToString(aesKey)
	enc, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPub, []byte(b64aes))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(enc), nil
}
