package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

// 创建AES加密key,16,24,32位字符串的话，分别对应AES-128，AES-192，AES-256 加密方法
var aesKey = "DIS**#KKKDJJSKDI"

func addPadding(originData []byte, blockSize int) []byte {
	paddingSize := blockSize - len(originData)%blockSize
	//Repeat()函数的功能是把切片[]byte{byte(padding)}复制padding个，然后合并成新的字节切片返回
	paddedData := bytes.Repeat([]byte{byte(paddingSize)}, paddingSize)

	return append(originData, paddedData...)
}

// UnPadding aes解密后，去掉padding
func UnPadding(originData []byte) ([]byte, error) {
	length := len(originData)
	if length == 0 {
		return nil, errors.New("加密字符串错误！")
	} else {
		// 前面计算padding加入的就是paddingSize的值
		unPaddingSize := int(originData[length-1])
		return originData[:length-unPaddingSize], nil
	}
}

// AesEncrypt aes加密
func AesEncrypt(originData []byte, aesKey []byte) (encryptedData []byte, err error) {
	// 创建加密实例
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	// 获取加密块大小
	blockSize := block.BlockSize()

	// 对数据进行填充：对数据进行填充，让数据长度满足需求
	originData = addPadding(originData, blockSize)

	// 采用AES加密方法中CBC加密模式
	blockMode := cipher.NewCBCEncrypter(block, aesKey[:blockSize])
	// 执行加密
	encryptedByte := make([]byte, len(originData))
	blockMode.CryptBlocks(encryptedByte, originData)

	return encryptedByte, nil
}

// AesDeEncrypt aes解密
func AesDeEncrypt(encryptedByte []byte, aesKey []byte) ([]byte, error) {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCEncrypter(block, aesKey[:blockSize])
	// 执行解密
	originData := make([]byte, len(encryptedByte))
	blockMode.CryptBlocks(originData, encryptedByte)

	// 去除填充字符串
	originByte, err := UnPadding(originData)
	if err != nil {
		return nil, err
	}
	return originByte, nil
}

// EncodeMes base64加密，将其aes加密的[]byte在加密成string返回，cookie中需要存储string
func EncodeMes(message []byte) (string, error) {
	encryptedByte, err := AesEncrypt(message, []byte(aesKey))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encryptedByte), nil
}

func DecodeMess(encodes string) ([]byte, error) {
	decodeByte, err := base64.StdEncoding.DecodeString(encodes)
	if err != nil {
		return nil, err
	}

	// 执行aes解密
	originData, err := AesDeEncrypt(decodeByte, []byte(aesKey))
	if err != nil {
		return nil, err
	}
	return originData, nil
}
