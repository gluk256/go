package crutils

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"

	"github.com/gluk256/crypto/algo/keccak"
	"github.com/gluk256/crypto/algo/rcx"
)

const (
	AesKeySize = 32
	AesSaltSize = 12
	AesEncryptedSizeDiff = 16
	SaltSize = 48
	MinDataSize = 64
	EncryptedSizeDiff = AesEncryptedSizeDiff + SaltSize
	DefaultRollover = 1024 * 256
)

const (
	offset = 256
	begK1 = offset
	endK1 = begK1 + offset
	begK2 = endK1 + offset
	endK2 = begK2 + offset
	begRcxKey = endK2 + offset
	endRcxKey = begRcxKey + offset
	begAesKey = endRcxKey + offset
	endAesKey = begAesKey + AesKeySize
	begAesSalt = endAesKey + offset
	endAesSalt = begAesSalt + AesSaltSize
	keyHolderSize = endAesSalt
)

func calculateIterations(sz int) int {
	const Mb = 1024 * 1024
	if sz < 128 * 1024 {
		return 1023
	} else if sz < Mb {
		return 511
	} else if sz < Mb * 2 {
		return 255
	} else if sz < Mb * 4 {
		return 127
	} else if sz < Mb * 8 {
		return 63
	} else if sz < Mb * 16 {
		return 31
	} else if sz < Mb * 25 {
		return 25
	} else if sz < Mb * 40 {
		return 19
	} else if sz < Mb * 70 {
		return 15
	} else if sz < Mb * 120 {
		return 13
	} else if sz < Mb * 180 {
		return 11
	} else if sz < Mb * 300 {
		return 9
	} else if sz < Mb * 450 {
		return 7
	} else if sz < Mb * 700 {
		return 5
	} else {
		return 3
	}
}

func EncryptInplaceKeccak(key []byte, data []byte) {
	var d keccak.Keccak512
	d.Write(key)
	d.ReadXor(data)
	b := make([]byte, keccak.Rate * 4)
	d.ReadXor(b) // cleanup internal state
	AnnihilateData(b) // prevent compiler optimization
}

// don't forget to clear the data
// key is expected to be 32 bytes, salt 12 bytes
func EncryptAES(key []byte, salt []byte, data []byte) ([]byte, error) {
	if len(key) != AesKeySize {
		fmt.Errorf("wrong key size %d", len(key))
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	encrypted := aesgcm.Seal(nil, salt, data, nil)
	return encrypted, err
}

// don't forget to clear the data
func DecryptAES(key []byte, salt []byte, data []byte) ([]byte, error) {
	if len(key) != AesKeySize {
		fmt.Errorf("wrong key size %d", len(key))
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	decrypted, err := aesgcm.Open(nil, salt, data, nil)
	return decrypted, err
}

// this is the main encryption function
func Encrypt(key []byte, data []byte) ([]byte, error) {
	data, _ = addPadding(data, 0, true)
	spacing := make([]byte, len(data))
	Randomize(spacing)
	return encryptWithSpacing(key, data, spacing)
}

// data must be already padded
func encryptWithSpacing(key []byte, data []byte, spacing []byte) ([]byte, error) {
	salt, err := generateSalt()
	if err != nil {
		return nil, err
	}
	data = addSpacing(data, spacing)
	keyholder := generateKeys(key, salt)
	defer AnnihilateData(keyholder)

	EncryptInplaceKeccak(keyholder[begK1:endK1], data)
	rcx.EncryptInplaceRCX(keyholder[begRcxKey:endRcxKey], data, calculateIterations(len(data)))
	data, err = EncryptAES(keyholder[begAesKey:endAesKey], keyholder[begAesSalt:endAesSalt], data)
	if err != nil {
		return nil, err
	}
	EncryptInplaceKeccak(keyholder[begK2:endK2], data)
	data = append(data, salt...)
	return data, nil
}

func Decrypt(key []byte, data []byte) (res []byte, spacing []byte, err error) {
	if len(data) <= SaltSize {
		return nil, nil, fmt.Errorf("data size %d, less than salt size %d", len(data), SaltSize)
	}
	res = data[:len(data)-SaltSize]
	salt := data[len(data)-SaltSize:]
	keyholder := generateKeys(key, salt)
	defer AnnihilateData(keyholder)

	EncryptInplaceKeccak(keyholder[begK2:endK2], res)
	res, err = DecryptAES(keyholder[begAesKey:endAesKey], keyholder[begAesSalt:endAesSalt], res)
	if err != nil {
		return nil, nil, err
	}
	rcx.DecryptInplaceRCX(keyholder[begRcxKey:endRcxKey], res, calculateIterations(len(res)))
	EncryptInplaceKeccak(keyholder[begK1:endK1], res)

	res, spacing = splitSpacing(res)
	res, err = removePadding(res)
	return res, spacing, err
}

// steganographic content is obviously of unknown size.
// however, we know that the size of original unencrypted steg content was power_of_two;
// so, we try all possible sizes (25 iterations at most, but in reality much less).
func DecryptStegContentOfUnknownSize(key []byte, steg []byte) ([]byte, []byte, error) {
	for sz := len(steg) / 2; sz >= MinDataSize; sz /= 2 {
		trySize := sz + EncryptedSizeDiff
		content := make([]byte, trySize)
		copy(content, steg[:trySize])
		res, steg, err := Decrypt(key, content)
		if err == nil {
			return res, steg, err
		}
	}
	return nil, nil, errors.New("failed to decrypt steganographic content")
}

func generateSalt() ([]byte, error) {
	salt := make([]byte, SaltSize)
	err := StochasticRand(salt)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Stochastic rand failed: %s", err.Error()))
	}
	return salt, err
}

func generateKeys(key []byte, salt []byte) []byte {
	fullkey := make([]byte, 0, len(key) + len(salt))
	fullkey = append(fullkey, key...)
	fullkey = append(fullkey, salt...)
	fullkey = fullkey[:len(fullkey)-1]
	keyholder := keccak.Digest(fullkey, keyHolderSize)
	AnnihilateData(fullkey)
	return keyholder
}

// encrypt with steganographic content as spacing
func EncryptSteg(key []byte, data []byte, steg []byte) (res []byte, err error) {
	data, _ = addPadding(data, 0, true)
	if len(data) < len(steg) + 4 { // four bytes for padding size
		return nil, fmt.Errorf("data size is less than necessary [%d vs. %d]", len(data), len(steg)+4)
	}
	steg, err = addPadding(steg, len(data), false) // no mark: steg content must be indistinguishable from random gamma
	if err != nil {
		return nil, err
	}

	res, err = encryptWithSpacing(key, data, steg)
	AnnihilateData(data)
	AnnihilateData(steg)
	return res, err
}
