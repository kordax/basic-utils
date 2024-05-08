/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package umap

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"database/sql"
	"fmt"
	"hash"
	"hash/fnv"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type hashTestStruct struct {
	ID       int
	Name     string
	IsActive bool
	Ratings  []int
	Metadata map[string]string
	Profile  *profile
	Tags     sql.NullString // Using SQL types to demonstrate handling of SQL data types
}

type profile struct {
	Title   string
	Company string
}

func testHashFunction(t *testing.T, hash hash.Hash, hashName string) {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	hashSet := make(map[int64]bool)
	collisionCount := 0

	for i := 0; i < 100000; i++ {
		testValue := hashTestStruct{
			ID:       rand.Intn(1000),
			Name:     fmt.Sprintf("Name%d", rand.Intn(1000)),
			IsActive: rand.Intn(2) == 1,
			Ratings:  []int{rand.Intn(5), rand.Intn(5), rand.Intn(5)},
			Metadata: map[string]string{"role": fmt.Sprintf("role%d", rand.Intn(100)), "access": fmt.Sprintf("access%d", rand.Intn(100))},
			Profile: &profile{
				Title:   fmt.Sprintf("Title%d", rand.Intn(100)),
				Company: fmt.Sprintf("Company%d", rand.Intn(100)),
			},
			Tags: sql.NullString{String: fmt.Sprintf("tag%d", rand.Intn(100)), Valid: rand.Intn(2) == 1},
		}

		hashCalc := computeHash(hash, testValue)
		if _, exists := hashSet[hashCalc]; exists {
			collisionCount++
		} else {
			hashSet[hashCalc] = true
		}
	}

	assert.Equal(t, 0, collisionCount, fmt.Sprintf("Detected %d collisions in %s hash values", collisionCount, hashName))
}

func TestSHA256(t *testing.T) {
	testHashFunction(t, sha256.New(), "SHA-256")
}

func TestSHA1(t *testing.T) {
	testHashFunction(t, sha1.New(), "SHA-1")
}

func TestSHA512(t *testing.T) {
	testHashFunction(t, sha512.New(), "SHA-512")
}

func TestFNV128(t *testing.T) {
	testHashFunction(t, fnv.New128(), "FNV-128")
}

func TestIdenticalStructs(t *testing.T) {
	struct1 := hashTestStruct{
		ID:       1,
		Name:     "John Doe",
		IsActive: true,
		Ratings:  []int{5, 4, 3, 2, 1},
		Metadata: map[string]string{"role": "admin", "access": "full"},
		Profile: &profile{
			Title:   "Manager",
			Company: "Tech Inc",
		},
		Tags: sql.NullString{String: "important", Valid: true},
	}

	struct2 := hashTestStruct{
		ID:       1,
		Name:     "John Doe",
		IsActive: true,
		Ratings:  []int{5, 4, 3, 2, 1},
		Metadata: map[string]string{"role": "admin", "access": "full"},
		Profile: &profile{
			Title:   "Manager",
			Company: "Tech Inc",
		},
		Tags: sql.NullString{String: "important", Valid: true},
	}

	// A third struct that differs slightly
	struct3 := hashTestStruct{
		ID:       1,
		Name:     "John Doe",
		IsActive: true,
		Ratings:  []int{5, 5, 3, 2, 1}, // Slight change in ratings
		Metadata: map[string]string{"role": "admin", "access": "full"},
		Profile: &profile{
			Title:   "Manager",
			Company: "Tech Inc",
		},
		Tags: sql.NullString{String: "important", Valid: true},
	}

	for i := 0; i < 1000; i++ {
		hash1 := computeHash(sha256.New(), struct1)
		hash2 := computeHash(sha256.New(), struct2)
		hash3 := computeHash(sha256.New(), struct3)

		assert.Equal(t, hash1, hash2, "Hashes should match for identical structs")
		assert.NotEqual(t, hash1, hash3, "Hashes should not match for different structs")
	}
}
