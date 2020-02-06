package snowflake

import (
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
)

var (
	uuidStr1 = "7232c1c3-f6d1-4aec-bedd-c7e4c10dc8d3"
	uuid1    = UUID{0x72, 0x32, 0xc1, 0xc3, 0xf6, 0xd1, 0x4a, 0xec, 0xbe, 0xdd, 0xc7, 0xe4, 0xc1, 0x0d, 0xc8, 0xd3}
	uuidStr2 = "ee2b1891-ccd3-4a23-9246-4ce40d20e740"
	uuid2    = UUID{0xee, 0x2b, 0x18, 0x91, 0xcc, 0xd3, 0x4a, 0x23, 0x92, 0x46, 0x4c, 0xe4, 0x0d, 0x20, 0xe7, 0x40}
)

func TestEncoding(t *testing.T) {
	encodingTest(t, uuidStr1, &uuid1)
	encodingTest(t, uuidStr2, &uuid2)
}

func TestSnowflakeFile(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "snoflake_test")
	if err != nil {
		t.Fatal(err)
	}
	fname := tmpfile.Name()
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}
	u, err := Get(fname)
	if err == nil {
		t.Fatal("Expected error")
	}
	u, err = Make(fname)
	if err == nil {
		t.Fatal("Expected error")
	}
	err = WriteNew(fname)
	if err != nil {
		t.Fatal("Unexpected error", err)
	}
	u, err = Get(fname)
	if err != nil {
		t.Fatal("Unexpected error", err)
	}
	_, err = Decode(u.Encode())
	if err != nil {
		t.Fatalf("Failed to decode: %v", err)
	}
	InvalidateCache()
	ug, err := Get(fname)
	if err != nil {
		t.Fatal("Unexpected error", err)
	}
	if !reflect.DeepEqual(*ug, *u) {
		t.Fatalf("Decode error; Expected: %v, got: %v", *ug, *u)
	}
	os.Remove(fname)
	InvalidateCache()
	u, err = Get(fname)
	if err == nil {
		t.Fatal("Expected error")
	}
	u1, err := Make(fname)
	if err != nil {
		t.Fatal("Unexpected error", err)
	}
	if reflect.DeepEqual(*u1, *ug) {
		t.Fatal("Should create new UUID")
	}
}

func encodingTest(t *testing.T, str string, expected *UUID) {
	u, err := Decode(str)
	if err != nil || u == nil {
		t.Fatalf("Failed to decode %s: %v", str, err)
	}
	if !reflect.DeepEqual(*u, *expected) {
		t.Fatalf("Decode error; Expected: %v, got: %v", *expected, *u)
	}
	if e := u.Encode(); e != str {
		t.Fatalf("Encode error; Expected: %v, got: %v", str, e)
	}
}
