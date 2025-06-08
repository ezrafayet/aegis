package domain

import "testing"

func TestComputeNameFingerprint(t *testing.T) {
	t.Run("should compute name fingerprint", func(t *testing.T) {
		name := "John Doe"
		fingerprint, err := ComputeNameFingerprint(name)
		if err != nil {
			t.Fatal(err)
		}
		if fingerprint != "320b8e6bef45211f0f57b618925f4193" {
			t.Fatal("expected fingerprint to be 320b8e6bef45211f0f57b618925f4193", fingerprint)
		}
	})
	t.Run("prevent naming collisions", func(t *testing.T) {
		name1 := "John Doe"
		fp1, err := ComputeNameFingerprint(name1)
		if err != nil {
			t.Fatal(err)
		}
		name2 := "John doé"
		fp2, err := ComputeNameFingerprint(name2)
		if err != nil {
			t.Fatal(err)
		}
		if fp1 != fp2 {
			t.Fatal("expected fingerprints to be the same", fp1, fp2)
		}
	})
	t.Run("2 different names should have different fingerprints", func(t *testing.T) {
		name1 := "John Doe"
		fp1, err := ComputeNameFingerprint(name1)
		if err != nil {
			t.Fatal(err)
		}
		name2 := "Jane Doe"
		fp2, err := ComputeNameFingerprint(name2)
		if err != nil {
			t.Fatal(err)
		}
		if fp1 == fp2 {
			t.Fatal("expected fingerprints to be different", fp1, fp2)
		}
	})
}