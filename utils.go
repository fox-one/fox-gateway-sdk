package gateway

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"os"

	uuid "github.com/satori/go.uuid"
)

var (
	// production
	DefaultPublicKey = `-----BEGIN RSA PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0/SMrN1Ki50xAD0mjIjA
NroYtZ+dtFh9i2gT8ANy9ObQplKJQedM5VDviEqnNiyNgQj6byso3EnykgG7JbpQ
qwt7XgAwO+uE01EdGi46G59DzvobBfwchmV9q9caHE0od95XukCq7vQzlpL/IS2+
BWaG6RjYeqcE7mxdmcVIzQ6ifcY4tfcAnEXqVz5kAcKM+GbLVDOhdeb3LPpkydNT
Li+q8vY1PrnnWDlGnJORosBuRS5IXab7QxojKFx1lrq4EvnKGeyB6m3+h14Ixlcv
/QO5p7RR4lI9hs11Ecatritck25xQQ+YO4n0gYAvScxV0t0nQGBjmsN11Nm4Hl1x
kwIDAQAB
-----END RSA PUBLIC KEY-----`

	pkEnvKey = "FOX_GATEWAY_PUBLIC_KEY"
)

func getPublicKey() (*rsa.PublicKey, error) {
	key := os.Getenv(pkEnvKey)
	if key == "" {
		key = DefaultPublicKey
	}

	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return nil, errors.New("decode pem failed")
	}

	pkey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaKey, ok := pkey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("key is not public key")
	}

	return rsaKey, nil
}

func rsaEncrypt(data []byte) (string, error) {
	pub, err := getPublicKey()
	if err != nil {
		return "", err
	}

	hash := sha256.New()
	random := rand.Reader

	encryptedData, err := rsa.EncryptOAEP(hash, random, pub, data, nil)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(encryptedData), nil
}

func signRequest(method, uri, body string) string {
	h := sha256.New()
	h.Write([]byte(method + uri + body))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func newNonce() string {
	return uuid.Must(uuid.NewV4()).String()
}

func MD5(str string) []byte {
	h := md5.New()
	h.Write([]byte(str))
	return h.Sum(nil)
}

func expendKey(key []byte) []byte {
	for len(key) < 16 {
		key = append(key, key...)
	}
	return key[:16]
}

/**
 *  PKCS7补码
 *  这里可以参考下http://blog.studygolang.com/167.html
 */
func PKCS7Padding(data []byte) []byte {
	blockSize := 16
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)

}

/**
 *  去除PKCS7的补码
 */
func UnPKCS7Padding(data []byte) []byte {
	length := len(data)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(data[length-1])
	if length <= unpadding {
		return nil
	}
	return data[:(length - unpadding)]
}

func Encrypt(data []byte, key, iv []byte) (string, error) {
	key = expendKey(key)
	iv = expendKey(iv)

	ckey, err := aes.NewCipher(key)
	if nil != err {
		return "", err
	}

	encrypter := cipher.NewCBCEncrypter(ckey, iv)

	// PKCS7补码
	str := PKCS7Padding([]byte(data))
	out := make([]byte, len(str))

	encrypter.CryptBlocks(out, str)

	return base64.StdEncoding.EncodeToString(out), nil
}

func Decrypt(base64Str string, key, iv []byte) ([]byte, error) {
	key = expendKey(key)
	iv = expendKey(iv)

	ckey, err := aes.NewCipher(key)
	if nil != err {
		return nil, err
	}

	decrypter := cipher.NewCBCDecrypter(ckey, iv)

	base64In, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, err
	}

	out := make([]byte, len(base64In))
	decrypter.CryptBlocks(out, base64In)

	// 去除PKCS7补码
	out = UnPKCS7Padding(out)
	if out == nil {
		return nil, nil
	}

	return out, nil
}
