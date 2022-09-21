package util

import (
	"encoding/json"
	"errors"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func FromMessage(from protoreflect.ProtoMessage, to interface{}) error {
	marshaller := protojson.MarshalOptions{
		UseProtoNames: true,
	}
	data, err := marshaller.Marshal(from)
	if err != nil {
		return errors.New("marshaller.Marshal")
	}
	err = json.Unmarshal(data, to)
	if err != nil {
		return errors.New("json.Unmarshal")
	}
	return nil
}

func ToMessage(from interface{}, to protoreflect.ProtoMessage) error {
	data, err := json.Marshal(from)
	if err != nil {
		return errors.New("marshaller.Marshal")
	}
	unmarshaller := protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}
	err = unmarshaller.Unmarshal(data, to)
	if err != nil {
		return errors.New("json.Unmarshal")
	}
	return nil
}

func Keys[K comparable, V any](m map[K]V) []K {
    keys := make([]K, 0, len(m))
    for k := range m {
        keys = append(keys, k)
    }
    return keys
}