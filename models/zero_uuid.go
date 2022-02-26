package models

import (
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/globalsign/mgo/bson"
	"github.com/gofrs/uuid"
)

var (
	byteGroups = []int{8, 4, 4, 4, 12}
)

type ZeroUUID uuid.UUID

func NewZeroUUIDFromstring(uidStr string) (ZeroUUID, error) {
	uid, err := uuid.FromString(uidStr)
	if err != nil {
		return ZeroUUID(uuid.Nil), nil
	}
	return ZeroUUID(uid), nil
}

func NewZeroUUIDFromUUID(uid *uuid.UUID) (ZeroUUID, error) {
	if uid == nil {
		return ZeroUUID(uuid.Nil), nil
	}

	return ZeroUUID(*uid), nil
}

func NewV4() ZeroUUID {
	uid, _ := uuid.NewV4()
	return ZeroUUID(uid)
}

func (zu ZeroUUID) IsZero() bool {
	return zu == ZeroUUID((uuid.UUID{}))
}

func (zu ZeroUUID) ToUUID() *uuid.UUID {
	if zu == ZeroUUID(uuid.Nil) {
		return nil
	}

	uid := uuid.UUID(zu)
	return &uid
}

func (zu ZeroUUID) ToBsonBinary() *bson.Binary {
	if zu == ZeroUUID(uuid.Nil) {
		return nil
	}

	uid := uuid.UUID(zu)
	return &bson.Binary{
		Kind: bson.BinaryUUID,
		Data: uid.Bytes(),
	}
}

func (zu ZeroUUID) NullUUID() NullUUID {
	var nullUID = NullUUID{}
	if zu == ZeroUUID((uuid.UUID{})) {
		nullUID.UUID = ZeroUUID((uuid.UUID{}))
		nullUID.Valid = false
		return nullUID
	}

	nullUID.UUID = zu
	nullUID.Valid = true
	return nullUID
}

func (zu ZeroUUID) Interface() interface{} {
	if zu == ZeroUUID((uuid.UUID{})) {
		return nil
	}

	return zu
}

func (zu ZeroUUID) MarshalJSON() ([]byte, error) {
	if zu == ZeroUUID((uuid.UUID{})) {
		return json.Marshal("")
	}
	return json.Marshal(uuid.UUID(zu).String())
}

func (zu ZeroUUID) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "" {
		return nil
	}
	_, err := uuid.FromString(s)
	if err != nil {
		return errors.New("invalid format uuid")
	}

	return nil
}

func (zu ZeroUUID) String() string {
	if zu == ZeroUUID((uuid.UUID{})) {
		return ""
	}
	return uuid.UUID(zu).String()
}

func (zu *ZeroUUID) Scan(src interface{}) error {
	switch src := src.(type) {
	case ZeroUUID: // support gorm convert from UUID to NullUUID
		*zu = src
		return nil

	case []byte:
		if len(src) == uuid.Size {
			return zu.UnmarshalBinary(src)
		}
		return zu.UnmarshalText(src)

	case string:
		return zu.UnmarshalText([]byte(src))
	}

	return fmt.Errorf("uuid: cannot convert %T to UUID", src)
}

func (zu ZeroUUID) Value() (driver.Value, error) {
	if zu == ZeroUUID(uuid.Nil) {
		return nil, nil
	}
	return zu.String(), nil
}

func (zu *ZeroUUID) UnmarshalBinary(data []byte) error {
	if len(data) != uuid.Size {
		return fmt.Errorf("uuid: UUID must be exactly 16 bytes long, got %d bytes", len(data))
	}
	copy(zu[:], data)

	return nil
}

func (zu ZeroUUID) MarshalBinary() ([]byte, error) {
	return uuid.UUID(zu).Bytes(), nil
}

func (zu *ZeroUUID) UnmarshalText(text []byte) error {
	switch len(text) {
	case 32:
		return zu.decodeHashLike(text)
	case 36:
		return zu.decodeCanonical(text)
	default:
		return fmt.Errorf("uuid: incorrect UUID length %d in string %q", len(text), text)
	}
}

// decodeHashLike decodes UUID strings that are using the following format:
//  "6ba7b8109dad11d180b400c04fd430c8".
func (u *ZeroUUID) decodeHashLike(t []byte) error {
	src := t[:]
	dst := u[:]

	_, err := hex.Decode(dst, src)
	return err
}

func (u *ZeroUUID) decodeCanonical(t []byte) error {
	if t[8] != '-' || t[13] != '-' || t[18] != '-' || t[23] != '-' {
		return fmt.Errorf("uuid: incorrect UUID format in string %q", t)
	}

	src := t
	dst := u[:]

	for i, byteGroup := range byteGroups {
		if i > 0 {
			src = src[1:] // skip dash
		}
		_, err := hex.Decode(dst[:byteGroup/2], src[:byteGroup])
		if err != nil {
			return err
		}
		src = src[byteGroup:]
		dst = dst[byteGroup/2:]
	}

	return nil
}
