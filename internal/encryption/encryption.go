package encryption

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"math/big"
)

func GenerateKeyPair() (privateKey, publicKey string) {
	// 使用椭圆曲线加密算法生成密钥对
	privateKeyECDSA, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	// 将私钥转换为十六进制字符串
	privateKeyBytes := privateKeyECDSA.D.Bytes()
	privateKey = hex.EncodeToString(privateKeyBytes)

	// 将公钥转换为十六进制字符串
	publicKeyBytes := elliptic.Marshal(privateKeyECDSA.PublicKey.Curve, privateKeyECDSA.PublicKey.X, privateKeyECDSA.PublicKey.Y)
	publicKey = hex.EncodeToString(publicKeyBytes)

	return privateKey, publicKey
}

func SignMessage(privateKey string, message string) (string, error) {
	// 将私钥从十六进制字符串转换回 *ecdsa.PrivateKey
	privateKeyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return "", err
	}
	privateKeyECDSA := new(ecdsa.PrivateKey)
	privateKeyECDSA.PublicKey.Curve = elliptic.P256()
	privateKeyECDSA.D = new(big.Int).SetBytes(privateKeyBytes)
	privateKeyECDSA.PublicKey.X, privateKeyECDSA.PublicKey.Y = privateKeyECDSA.PublicKey.Curve.ScalarBaseMult(privateKeyECDSA.D.Bytes())

	// 对消息进行哈希
	hash := sha256.Sum256([]byte(message))

	// 签名
	r, s, err := ecdsa.Sign(rand.Reader, privateKeyECDSA, hash[:])
	if err != nil {
		return "", err
	}

	// 将签名转换为字节数组
	signature := append(r.Bytes(), s.Bytes()...)

	// 将签名转换为十六进制字符串
	return hex.EncodeToString(signature), nil
}

func VerifySignature(publicKey string, message string, signature string) (bool, error) {
	// 将公钥从十六进制字符串转换回 *ecdsa.PublicKey
	publicKeyBytes, err := hex.DecodeString(publicKey)
	if err != nil {
		return false, err
	}
	x, y := elliptic.Unmarshal(elliptic.P256(), publicKeyBytes)
	publicKeyECDSA := &ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}

	// 将签名从十六进制字符串转换回字节数组
	signatureBytes, err := hex.DecodeString(signature)
	if err != nil {
		return false, err
	}

	// 分离 r 和 s
	r := new(big.Int).SetBytes(signatureBytes[:len(signatureBytes)/2])
	s := new(big.Int).SetBytes(signatureBytes[len(signatureBytes)/2:])

	// 对消息进行哈希
	hash := sha256.Sum256([]byte(message))

	// 验证签名
	return ecdsa.Verify(publicKeyECDSA, hash[:], r, s), nil
}
