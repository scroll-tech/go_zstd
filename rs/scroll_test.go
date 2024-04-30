package rs

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestCompressScrollBatchBytes(t *testing.T) {
	var tests []struct {
		filename     string
		rawSize      int
		comprSize    int
		expectedHash common.Hash
	}

	data, err := os.ReadFile("../testdata/input.txt")
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		var batch string
		var rawSize, comprSize int
		var comprKeccakHash string

		n, err := fmt.Sscanf(line, "%s raw_size= %d, compr_size= %d, compr_keccak_hash=%s", &batch, &rawSize, &comprSize, &comprKeccakHash)
		if err != nil || n != 4 {
			t.Fatalf("failed to parse line: %s, error: %v", line, err)
		}

		batch = strings.TrimSuffix(batch, ",")
		//fmt.Println("test", fmt.Sprintf("%s.hex", batch), rawSize, comprSize, common.HexToHash(comprKeccakHash))

		tests = append(tests, struct {
			filename     string
			rawSize      int
			comprSize    int
			expectedHash common.Hash
		}{
			filename:     fmt.Sprintf("../testdata/%s.hex", batch),
			rawSize:      rawSize,
			comprSize:    comprSize,
			expectedHash: common.HexToHash(comprKeccakHash),
		})
	}

	for _, test := range tests {
		hexData, err := os.ReadFile(test.filename)
		if err != nil {
			t.Fatalf("failed to read file %s: %v", test.filename, err)
		}

		batchBytes, err := hex.DecodeString(strings.TrimSpace(string(hexData)))
		if err != nil {
			t.Fatalf("failed to decode hex data from file %s: %v", test.filename, err)
		}

		if len(batchBytes) != test.rawSize {
			t.Errorf("raw size mismatch for file %s: expected %d, got %d", test.filename, test.rawSize, len(batchBytes))
		}

		compressedBytes, err := CompressScrollBatchBytes(batchBytes)
		if err != nil {
			t.Errorf("CompressScrollBatchBytes failed for file %s: %v", test.filename, err)
		}

		if len(compressedBytes) != test.comprSize {
			fmt.Println(hex.EncodeToString(compressedBytes))
			t.Errorf("compressed size mismatch for file %s: expected %d, got %d", test.filename, test.comprSize, len(compressedBytes))
		}

		hash := crypto.Keccak256Hash(compressedBytes)
		if hash != test.expectedHash {
			t.Errorf("hash mismatch for file %s: expected %v, got %v", test.filename, test.expectedHash, hash)
		}
	}
}
