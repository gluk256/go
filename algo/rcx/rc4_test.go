package rcx

import (
	"bytes"
	"encoding/hex"
	mrand "math/rand"
	"testing"
	"time"

	"github.com/gluk256/crypto/algo/primitives"
)

func TestKeyStream(t *testing.T) {
	testSingleKeyStream(t, "Key", "EB9F7781B734CA72A719")
	testSingleKeyStream(t, "Wiki", "6044DB6D41B7")
	testSingleKeyStream(t, "Secret", "04D46B053CA87B59")

	testSingleEncrypt(t, "Key", "Plaintext", "BBF316E8D940AF0AD3")
	testSingleEncrypt(t, "Wiki", "pedia", "1021BF0420")
	testSingleEncrypt(t, "Secret", "Attack at dawn", "45A01F645FC35B383552544B9BF5")

	var k []byte
	for i := 0; i < 32; i++ {
		k = append(k, byte(i+1))
	}

	testSingleKeyStream(t, string(k), "eaa6bd25880bf93d3f5d1e4ca2611d91cfa45c9f7e714b54bdfa80027cb14380")
	testSingleKeyStream(t, string(k[:24]), "0595e57fe5f0bb3c706edac8a4b2db11dfde31344a1af769c74f070aee9e2326")
}

func testSingleKeyStream(t *testing.T, key string, expected string) {
	data := make([]byte, len(expected)/2)
	exp := make([]byte, len(expected)/2)
	hex.Decode(exp, []byte(expected))

	var r RC4
	r.InitKey([]byte(key))
	r.j = 0
	r.XorInplace(data)
	if !bytes.Equal(data, exp) {
		t.Fatalf("wrong keystream, key: %s", key)
	}
}

func testSingleEncrypt(t *testing.T, key string, data string, expected string) {
	d := []byte(data)
	exp := make([]byte, len(expected)/2)
	hex.Decode(exp, []byte(expected))

	var r RC4
	r.InitKey([]byte(key))
	r.j = 0
	r.XorInplace(d)
	if !bytes.Equal(d, exp) {
		t.Fatalf("encryption failed, key: %s", key)
	}
}

func generateRandomBytes(t *testing.T, align bool) []byte {
	sz := mrand.Intn(256) + 256
	if align {
		for sz%4 != 0 {
			sz++
		}
	}
	b := make([]byte, sz)
	_, err := mrand.Read(b)
	if err != nil {
		t.Fatal("failed to generate random bytes")
	}
	return b
}

func TestEncryptionRC4(t *testing.T) {
	seed := time.Now().Unix()
	mrand.Seed(seed)

	for i := 0; i < 64; i++ {
		key := generateRandomBytes(t, false)
		x := generateRandomBytes(t, false)
		y := make([]byte, len(x))
		copy(y, x)

		var re, rd RC4
		re.InitKey(key)
		rd.InitKey(key)

		re.XorInplace(y)
		if bytes.Equal(x, y) {
			t.Fatalf("failed encrypt, round %d with seed %d", i, seed)
		}
		ok := primitives.IsDeepNotEqual(x, y, len(x))
		if !ok {
			t.Fatalf("failed encrypt deep check, round %d with seed %d", i, seed)
		}

		rd.XorInplace(y)
		if !bytes.Equal(x, y) {
			t.Fatalf("failed decrypt, round %d with seed %d", i, seed)
		}
	}
}

func TestEncryptionMix(t *testing.T) {
	seed := time.Now().Unix()
	mrand.Seed(seed)

	key := generateRandomBytes(t, false)
	orig := generateRandomBytes(t, false)
	buf := make([]byte, len(orig))
	encrypted := make([]byte, len(orig))
	copy(buf, orig)

	var re RC4
	re.InitKey(key)
	dummy := make([]byte, 16*256*256)
	re.XorInplace(dummy) // roll forward

	re.XorInplace(buf)
	copy(encrypted, buf)
	if bytes.Equal(orig, buf) {
		t.Fatalf("failed encrypt, with seed %d", seed)
	}
	if !primitives.IsDeepNotEqual(orig, buf, len(orig)) {
		t.Fatalf("failed encrypt deep check, with seed %d", seed)
	}

	DecryptInplaceRcx(key, buf, 0)
	if !bytes.Equal(orig, buf) {
		t.Fatalf("failed decrypt, with seed %d", seed)
	}

	EncryptInplaceRcx(key, buf, 0)
	if !bytes.Equal(encrypted, buf) {
		t.Fatalf("failed decrypt, with seed %d", seed)
	}
}

func TestConversion(t *testing.T) {
	a := byte(0xad)
	b := byte(0xde)
	v := Bytes2uint(a, b)
	if v != 0xDEAD {
		t.Fatalf("Bytes2uint failed, val = %x", v)
	}
	y, z := Uint2bytes(v)
	if y != a || z != b {
		t.Fatalf("Uint2bytes failed [%x, %x, %x, %x]", a, b, y, z)
	}
	v = Bytes2uint(0xa1, 0x0f)
	if v != 0x0FA1 {
		t.Fatalf("Bytes2uint failed second run, val = %x", v)
	}
	y, z = Uint2bytes(v)
	if y != 0xa1 || z != 0x0f {
		t.Fatalf("Uint2bytes failed second run [%x, %x, %x, %x]", a, b, y, z)
	}
}

func TestSingleRunRCX(t *testing.T) {
	seed := time.Now().Unix()
	mrand.Seed(seed)

	for i := 0; i < 4; i++ {
		key := generateRandomBytes(t, false)
		x := generateRandomBytes(t, true)
		y := make([]byte, len(x))
		copy(y, x)

		var cipher RCX
		cipher.InitKey(key)

		cipher.encryptSingleRun(y)
		if bytes.Equal(x, y) {
			t.Fatalf("failed encrypt, round %d with seed %d", i, seed)
		}
		ok := primitives.IsDeepNotEqual(x, y, len(x))
		if !ok {
			t.Fatalf("failed encrypt deep check, round %d with seed %d", i, seed)
		}

		cipher.encryptSingleRun(y)
		if !bytes.Equal(x, y) {
			t.Fatalf("failed decrypt, round %d with seed %d", i, seed)
		}
	}
}

// encrypt array containing zeros
func TestSingleRunRcxZero(t *testing.T) {
	seed := time.Now().Unix()
	mrand.Seed(seed)
	const sz = 1024 * 64

	for i := 0; i < 3; i++ {
		key := generateRandomBytes(t, false)
		x := make([]byte, sz)
		zero := make([]byte, sz)

		var cipher RCX
		cipher.InitKey(key)
		cipher.encryptSingleRun(x)
		ok := primitives.IsDeepNotEqual(x, zero, sz)
		if !ok {
			t.Fatalf("failed encrypt deep check, round %d with seed %d", i, seed)
		}

		cipher.EncryptCascade(x, 255)
		ok = primitives.IsDeepNotEqual(x, zero, sz)
		if !ok {
			t.Fatalf("failed encrypt deep check, round %d with seed %d", i, seed)
		}
	}
}

func TestCascade(t *testing.T) {
	seed := time.Now().Unix()
	mrand.Seed(seed)

	for i := 0; i < 3; i++ {
		key := generateRandomBytes(t, false)
		x := generateRandomBytes(t, true)
		y := make([]byte, len(x))
		copy(y, x)

		var c RCX
		c.InitKey(key)

		c.EncryptCascade(y, 511)
		if bytes.Equal(x, y) {
			t.Fatalf("failed encrypt, round %d with seed %d", i, seed)
		}
		ok := primitives.IsDeepNotEqual(x, y, len(x))
		if !ok {
			t.Fatalf("failed encrypt deep check, round %d with seed %d", i, seed)
		}

		c.DecryptCascade(y, 511)
		if !bytes.Equal(x, y) {
			t.Fatalf("failed decrypt, round %d with seed %d", i, seed)
		}
	}
}

func TestAvalancheRcx(t *testing.T) {
	seed := time.Now().Unix()
	mrand.Seed(seed)

	for i := 0; i < 3; i++ {
		key := generateRandomBytes(t, false)
		x := generateRandomBytes(t, true)
		y := make([]byte, len(x))
		z := make([]byte, len(x))
		copy(y, x)
		copy(z, x)

		var cipher RCX
		cipher.InitKey(key)

		x[0]-- // change at least one bit, which is supposed to cause an avalanche effect
		cycles := len(x)
		cipher.EncryptCascade(x, cycles)
		cipher.EncryptCascade(y, cycles)
		cipher.EncryptCascade(z, cycles)

		if !bytes.Equal(y, z) {
			t.Fatalf("failed to encrypt, round %d with seed %d", i, seed)
		}

		ok := primitives.IsDeepNotEqual(x, y, len(x))
		if !ok {
			t.Fatalf("failed deep check, round %d with seed %d and len=%d", i, seed, len(x))
		}
	}
}

func TestAvalancheRC4(t *testing.T) {
	seed := time.Now().Unix()
	mrand.Seed(seed)

	for i := 0; i < 128; i++ {
		key := generateRandomBytes(t, false)
		var a, b RC4
		a.InitKey(key)
		b.InitKey(key)
		b.s[0], b.s[1] = b.s[1], b.s[0]

		// usually it takes no more than 3 iterations, but we allow 8
		var done bool
		for j := 0; j < 8 && !done; j++ {
			y := make([]byte, 256)
			x := make([]byte, 256)
			a.XorInplace(x)
			b.XorInplace(y)
			done = primitives.IsDeepNotEqual(x, y, len(x))
		}
		if !done {
			t.Fatalf("failed with seed %d, iter %d", seed, i)
		}
	}
}

func TestEncryptionRCX(t *testing.T) {
	seed := time.Now().Unix()
	mrand.Seed(seed)

	for i := 0; i < 3; i++ {
		key := generateRandomBytes(t, false)
		x := generateRandomBytes(t, false)
		y := make([]byte, len(x))
		copy(y, x)

		EncryptInplaceRcx(key, y, 511)
		if bytes.Equal(x, y) {
			t.Fatalf("failed encrypt, round %d with seed %d", i, seed)
		}
		ok := primitives.IsDeepNotEqual(x, y, len(x))
		if !ok {
			t.Fatalf("failed encrypt deep check, round %d with seed %d", i, seed)
		}

		DecryptInplaceRcx(key, y, 511)
		if !bytes.Equal(x, y) {
			t.Fatalf("failed decrypt, round %d with seed %d\n%x\n%x", i, seed, x, y)
		}
	}
}

// tests the ability to generate consistent gamma
func TestConsistencyRC4(t *testing.T) {
	seed := time.Now().Unix()
	mrand.Seed(seed)
	b := make([]byte, 1024)
	for i := 0; i < 16; i++ {
		key := generateRandomBytes(t, false)
		sz := Bytes2uint(key[0], key[1])
		x := make([]byte, sz)
		y := make([]byte, sz)

		var r1, r2 RC4
		r1.InitKey(key)
		r2.InitKey(key)

		for j := 0; j < 170; j++ {
			r1.XorInplace(b[:33])
		}
		for j := 0; j < 330; j++ {
			r2.XorInplace(b[:17])
		}

		r1.XorInplace(x)
		r2.XorInplace(y)
		if !bytes.Equal(x, y) {
			t.Fatalf("failed to generate consistent gamma, round %d with seed %d", i, seed)
		}
	}
}

func TestCleanupRcx(t *testing.T) {
	seed := time.Now().Unix()
	mrand.Seed(seed)

	for i := 0; i < 3; i++ {
		key := generateRandomBytes(t, false)
		var x RCX
		x.InitKey(key)
		prev := make([]byte, 256)
		copy(prev, x.rc4.s[:])
		x.cleanup()
		if !primitives.IsDeepNotEqual(prev, x.rc4.s[:], 256) {
			t.Fatalf("cleanup failed, round %d with seed %d", i, seed)
		}
	}
}
