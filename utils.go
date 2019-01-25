package gateway

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"

	uuid "github.com/satori/go.uuid"
)

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
