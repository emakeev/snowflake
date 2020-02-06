// package snowflake provides API for managing & using snowflake UUIDs
package snowflake

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"sync/atomic"
	"time"
	"unsafe"
)

const (
	DefaultSnowflakeFile = "/etc/snowflake"
	ValidationTimePeriod = time.Minute * 5
	UUIDHexLen           = 36
)

// RFC 4122 UUID
type UUID [16]byte

type cache struct {
	uuid           UUID
	filePath       string
	validationTime time.Time
}

// Get returns UUID from a file specified by the first argument or default location if no argument given
func Get(args ...string) (*UUID, error) {
	fpath := sfFilePath(args)
	cached := (*cache)(atomic.LoadPointer(&current))
	if cached == nil || cached.filePath != fpath || cached.validationTime.Before(time.Now()) {
		return Read(fpath) // first time read OR cache refresh
	}
	res := cached.uuid
	return &res, nil
}

// Read loads UUID from a file specified by the first argument or default location, updates cache & returns parsed UUID
func Read(args ...string) (*UUID, error) {
	fpath := sfFilePath(args)
	f, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	data := make([]byte, UUIDHexLen+8)
	n, err := f.Read(data)
	if err != nil {
		return nil, err
	}
	if n < UUIDHexLen {
		return nil, fmt.Errorf("File '%s' is too small", fpath)
	}
	data = data[:n]
	uuid, err := Decode(string(data))
	if err == nil {
		nc := &cache{uuid: *uuid, filePath: fpath, validationTime: time.Now().Add(ValidationTimePeriod)}
		atomic.StorePointer(&current, unsafe.Pointer(nc))
	}
	return uuid, err
}

// Make returns snowflake ID if snowflake file exists. Otherwise, it creates one.
func Make(args ...string) (*UUID, error) {
	fpath := sfFilePath(args)
	u, err := Get(fpath)
	if err == nil {
		return u, nil
	}
	perr, ok := err.(*os.PathError)
	if ok && perr != nil && perr.Op == "open" {
		if err = WriteNew(fpath); err == nil {
			return Get(fpath)
		}
	}
	return nil, err
}

// WriteNew generates a new UUID, encodes it and writes it into fname file
// WriteNew will overwrite UUID if the file already exist
func WriteNew(fname string) error {
	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	_, err = f.WriteString(Gen().Encode())
	f.Close()
	if err == nil {
		InvalidateCache()
	}
	return err
}

// InvalidateCache resets internal cache, next Get will trigger a new Read
func InvalidateCache() {
	atomic.StorePointer(&current, nil) // invalidate cache
}

// Encode returns string representation of the UUID
func (u *UUID) Encode() string {
	if u == nil {
		u = &zeroUUID
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}

// String interface for UUID, see Encode()
func (u *UUID) String() string {
	return u.Encode()
}

// Decode decodes given UUID String into a new UUID and returns it
func Decode(uuidStr string) (*UUID, error) {
	uuidStr = strings.ToLower(strings.TrimSpace(string(uuidStr)))
	if !reMatch.MatchString(uuidStr) {
		return nil, fmt.Errorf("Decode error: '%s' is not a valid UUID", uuidStr)
	}
	var (
		res UUID
		i   int
		err error
	)
	dst := res[:]
	for i = range decoder {
		d := &decoder[i]
		if i, err = hex.Decode(dst[d.bLeft:d.bRight], []byte(uuidStr[d.sLeft:d.sRight])); err != nil {
			return &res, fmt.Errorf(
				"Error while decoding '%s' of '%s': %v", uuidStr[d.sLeft:d.sRight], uuidStr, err)
		}
	}
	return &res, nil
}

// Gen generates a new (random) UUID
func Gen() *UUID {
	u := new(UUID)
	rand.Read(u[:])
	u[8] = (u[8] | 0x40) & 0x7F
	u[6] = (u[6] & 0x0F) | 0x40
	return u
}

func sfFilePath(args []string) string {
	if len(args) == 0 {
		return DefaultSnowflakeFile
	} else {
		return args[0]
	}
}

var (
	current  unsafe.Pointer
	reMatch  *regexp.Regexp
	zeroUUID = UUID{}
	decoder  = [5]struct{ bLeft, bRight, sLeft, sRight int }{
		{0, 4, 0, 8},
		{4, 6, 9, 13},
		{6, 8, 14, 18},
		{8, 10, 19, 23},
		{10, 16, 24, 36},
	}
)

func init() {
	reMatch = regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`)
	rand.Seed(time.Now().UnixNano())
}
