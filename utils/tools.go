package utils

import (
	"encoding/json"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

func GenerateUUID() uuid.UUID {
	uid, _ := uuid.NewUUID()
	return uid
}

func GenerateUUIDBytes() []byte {
	uid, _ := uuid.NewUUID()
	return []byte(uid.String())
}

func UUIDStringByBytes(bytes []byte) (string, error) {
	uuid, err := uuid.FromBytes(bytes)
	return uuid.String(), err
}

/**
 *
 *
 *
**/
func PbMarshal(obj proto.Message) ([]byte, error) {
	bytes, err := proto.Marshal(obj)
	return bytes, err
}
func PbUnMarshal(bytes []byte, typeScope proto.Message) error {
	err := proto.Unmarshal(bytes, typeScope)
	return err
}

func JsonMarshal(obj interface{}) ([]byte, error) {
	return json.Marshal(obj)
}

func JsonUnMarshal(bytes []byte, obj interface{}) error {
	return json.Unmarshal(bytes, obj)
}
