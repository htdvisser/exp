package envcrypto

import (
	"testing"
)

func TestNew(t *testing.T) {
	t.Parallel()

	m, err := New()
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}
	if m == nil {
		t.Fatal("New() returned nil map")
	}

	if _, ok := m["ENVCRYPTO_PUBLIC_KEY"]; !ok {
		t.Error("New() did not set ENVCRYPTO_PUBLIC_KEY")
	}
	if _, ok := m["ENVCRYPTO_PRIVATE_KEY"]; !ok {
		t.Error("New() did not set ENVCRYPTO_PRIVATE_KEY")
	}
}

func TestOpen(t *testing.T) {
	t.Parallel()

	m, _ := New()

	box, err := Open(m)
	if err != nil {
		t.Fatalf("Open() returned error: %v", err)
	}
	if box == nil {
		t.Fatal("Open() returned nil box")
	}

	if box.publicKey == nil {
		t.Error("Open() did not set box.publicKey")
	}
	if box.privateKey == nil {
		t.Error("Open() did not set box.privateKey")
	}
}

func TestOpen_withoutKeys(t *testing.T) {
	t.Parallel()

	_, err := Open(MapSource{})
	if err == nil {
		t.Fatal("Open() did not return error")
	}
}

func TestOpen_withInvalidPublicKey(t *testing.T) {
	t.Parallel()

	m, _ := New()
	m["ENVCRYPTO_PUBLIC_KEY"] = "invalid"

	_, err := Open(m)
	if err == nil {
		t.Fatal("Open() did not return error")
	}
}

func TestOpen_withInvalidPrivateKey(t *testing.T) {
	t.Parallel()

	m, _ := New()
	m["ENVCRYPTO_PRIVATE_KEY"] = "invalid"

	_, err := Open(m)
	if err == nil {
		t.Fatal("Open() did not return error")
	}
}

func TestBox_Encrypt(t *testing.T) {
	t.Parallel()

	m, _ := New()
	box, _ := Open(m)

	encryptedValue, err := box.Encrypt("foo")
	if err != nil {
		t.Fatalf("Encrypt() returned error: %v", err)
	}
	if encryptedValue == "" {
		t.Fatal("Encrypt() returned empty string")
	}
	if encryptedValue == "foo" {
		t.Error("Encrypt() returned unencrypted value")
	}
}

func TestBox_Decrypt(t *testing.T) {
	t.Parallel()

	m, _ := New()
	box, _ := Open(m)

	encryptedValue, _ := box.Encrypt("foo")

	decryptedValue, err := box.Decrypt(encryptedValue)
	if err != nil {
		t.Fatalf("Decrypt() returned error: %v", err)
	}
	if decryptedValue == "" {
		t.Fatal("Decrypt() returned empty string")
	}
	if decryptedValue != "foo" {
		t.Error("Decrypt() returned wrong value")
	}
}

func TestBox_Get(t *testing.T) {
	t.Parallel()

	m, _ := New()
	box, _ := Open(m)

	m["FOO"], _ = box.Encrypt("foo")

	decryptedValue, err := box.Get("FOO")
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}
	if decryptedValue == "" {
		t.Fatal("Get() returned empty string")
	}
	if decryptedValue != "foo" {
		t.Error("Get() returned wrong value")
	}
}

func TestBox_Get_noValue(t *testing.T) {
	t.Parallel()

	m, _ := New()
	box, _ := Open(m)

	_, err := box.Get("FOO")
	if err == nil {
		t.Fatal("Get() did not return error")
	}
}

func TestBox_Get_unencryptedValue(t *testing.T) {
	t.Parallel()

	m, _ := New()
	box, _ := Open(m)

	m["FOO"] = "foo"

	decryptedValue, err := box.Get("FOO")
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}
	if decryptedValue == "" {
		t.Fatal("Get() returned empty string")
	}
	if decryptedValue != "foo" {
		t.Error("Get() returned wrong value")
	}
}

func TestBox_All(t *testing.T) {
	t.Parallel()

	m, _ := New()
	box, _ := Open(m)

	m["FOO"], _ = box.Encrypt("foo")
	m["BAR"], _ = box.Encrypt("bar")

	all, err := box.All()
	if err != nil {
		t.Fatalf("All() returned error: %v", err)
	}
	if all == nil {
		t.Fatal("All() returned nil map")
	}
	if len(all) != 2 {
		t.Fatalf("All() returned map with %d values, want 2", len(all))
	}
	if all["FOO"] != "foo" {
		t.Errorf("All() returned wrong value for FOO: %q, want %q", all["FOO"], "foo")
	}
	if all["BAR"] != "bar" {
		t.Errorf("All() returned wrong value for BAR: %q, want %q", all["BAR"], "bar")
	}
}
