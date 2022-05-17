package concordium

import (
	"encoding/hex"
	"reflect"
	"strings"
)

const (
	modelStructFiledTagName        = "concordium"
	modelStructFiledTagValueParam  = "model"
	modelStructFiledTagValueOption = "option"
)

type SerializeModel interface {
	SerializeModel() ([]byte, error)
}

type Model string

func (m *Model) Serialize(v any) error {
	b, err := serializeModel(v)
	if err != nil {
		return err
	}
	*m = Model(hex.EncodeToString(b))
	return nil
}

func parseStructFieldTag(field reflect.StructField) (bool, map[string]bool) {
	t, ok := field.Tag.Lookup(modelStructFiledTagName)
	if !ok {
		return false, nil
	}
	m := map[string]bool{}
	for _, v := range strings.Split(t, ",") {
		m[strings.TrimSpace(v)] = true
	}
	return true, m
}
