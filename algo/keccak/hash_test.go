package keccak

import (
	"testing"
	"bytes"
	"encoding/hex"
)

var input = []string{
	"",
	"zxcv",
	" The supply of government exceeds the demand by huge margin. ",
	"Ineptocracy - a system of government where the least capable to lead are elected by the least capable to produce, and where the members of society least likely to sustain themselves or succeed, are rewarded with goods and services paid for by the confiscated wealth of a diminishing number of producers.",
	"Of all tyrannies, a tyranny sincerely exercised for the good of its victims may be the most oppressive. It would be better to live under robber barons than under omnipotent moral busybodies. The robber baron's cruelty may sometimes sleep, his cupidity may at some point be satiated; but those who torment us for our own good will torment us without end for they do so with the approval of their own conscience.",
}

var expected = []string{
	"0eab42de4c3ceb9235fc91acffe746b29c29a8c366b7c60e4e67c466f36a4304c00fa9caf9d87976ba469bcbe06713b435f091ef2769fb160cdab33d3670680e",
	"e7ff11ec516b91fa7feb61ebf19c89487b78bf49a077824efd08e392549522817b24ca65776341dddd91ce499951bca6bf267d00cf1e81629441e32eb70e2111",
	"1e737292b0b2d00227eb29b851ffd92c00908e44a51ef866fe7a934421b54191bafb86f4b46adf4252d2e6c5f3c7d04954045fdcea04b7d1e5057e94a2e7b1f6",
	"3d0f32465f276ca3b9b78832e2deaff4ae4c75d09db9ab98bcef68869d65c7e7f0c19807d791d5ffe4f4dcc7700397a590325bf8bf78b3b7bef7af64574e572c",
	"0b4726f2c9b79347d1f2340ee2ba35a6d9711dd84d6bcde7907135f0c57f4cedb3205ccb2b436b81510f199e996c3b3601ec2a92456718165c62a43e09ab5c11",
}

const sz = 64

func TestHash(t *testing.T) {
	exp := make([]byte, sz)
	for i, text := range input {
		hash := Digest([]byte(text), nil, sz)
		hex.Decode(exp, []byte(expected[i]))
		if bytes.Compare(hash, exp) != 0 {
			t.Fatalf("failed test number %d, result: \n[%x]", i, hash)
		}
	}

	res := make([]byte, sz)
	var k Keccak512
	for i := 0; i < len(input); i++ {
		k.Reset()
		k.Write([]byte(input[i]))
		k.Read(res, sz)
		hex.Decode(exp, []byte(expected[i]))
		if bytes.Compare(res, exp) != 0 {
			t.Fatalf("failed advanced test number %d, result: \n[%x]", i, res)
		}
	}
}

func BenchmarkKeccak512(b *testing.B) {
	buf := make([]byte, sz)
	var k Keccak512
	k.Write([]byte(input[3]))
	for i := 0; i < b.N; i++ {
		k.Read(buf, sz)
	}
}