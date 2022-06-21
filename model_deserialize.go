package concordium

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"time"
)

func DeserializeModel(b []byte, v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("nil or non-pointer given")
	}
	i, err := deserializeModelReflect(b, rv)
	if err != nil {
		return err
	}
	if i < len(b) {
		return fmt.Errorf("unexpected byte %d left", len(b)-i)
	}
	return nil
}

func deserializeModelReflect(b []byte, rv reflect.Value) (int, error) {
	if rv.CanSet() && rv.Kind() == reflect.Ptr {
		rv.Set(reflect.New(rv.Type().Elem()))
	}
	rv = reflect.Indirect(rv)
	for _, h := range []func([]byte, reflect.Value) (int, bool, error){
		deserializeModelTryCustom,
		deserializeModelTryTime,
		deserializeModelTryDuration,
	} {
		if i, ok, err := h(b, rv); ok {
			return i, err
		}
	}
	if !rv.CanSet() && rv.Kind() != reflect.Struct {
		return 0, fmt.Errorf("unsettable value")
	}
	var i int
	switch rv.Kind() {
	case reflect.Bool:
		i = 1
		if len(b) < i {
			return 0, fmt.Errorf("%q requires %d bytes", rv.Kind(), i)
		}
		rv.SetBool(b[0] == 1)
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		// Converting between int64 and uint64 doesn't change the sign bit,
		// only the way it's interpreted. So let's do that with magic below.
		var h func([]byte) uint64
		switch rv.Kind() {
		case reflect.Int64, reflect.Uint64:
			i = 8
			h = binary.LittleEndian.Uint64
		case reflect.Int32, reflect.Uint32:
			i = 4
			h = func(b []byte) uint64 { return uint64(binary.LittleEndian.Uint32(b)) }
		case reflect.Int16, reflect.Uint16:
			i = 2
			h = func(b []byte) uint64 { return uint64(binary.LittleEndian.Uint16(b)) }
		default:
			i = 1
			h = func(b []byte) uint64 { return uint64(b[0]) }
		}
		if len(b) < i {
			return 0, fmt.Errorf("%q requires %d bytes", rv.Kind(), i)
		}
		if rv.Kind() == reflect.Int8 || rv.Kind() == reflect.Int16 || rv.Kind() == reflect.Int32 || rv.Kind() == reflect.Int64 {
			rv.SetInt(int64(h(b)))
		} else {
			rv.SetUint(h(b))
		}
	case reflect.String:
		i = 4
		if len(b) < i {
			return 0, fmt.Errorf("%q requires %d bytes", rv.Kind(), i)
		}
		n := int(binary.LittleEndian.Uint32(b))
		rv.SetString(string(b[i : i+n]))
		i += n
	case reflect.Slice:
		i = 4
		if len(b) < i {
			return 0, fmt.Errorf("%q requires %d bytes", rv.Kind(), i)
		}
		n := int(binary.LittleEndian.Uint32(b))
		rs := reflect.MakeSlice(rv.Type(), 0, n)
		for j := 0; j < n; j++ {
			re := reflect.Indirect(reflect.New(reflect.Zero(rv.Type().Elem()).Type()))
			x, err := deserializeModelReflect(b[i:], re)
			if err != nil {
				return 0, err
			}
			i += x
			rs = reflect.Append(rs, re)
		}
		rv.Set(rs)
	case reflect.Array:
		i = 0
		n := rv.Cap()
		for j := 0; j < n; j++ {
			re := rv.Index(j)
			x, err := deserializeModelReflect(b[i:], re)
			if err != nil {
				return 0, err
			}
			i += x
		}
	case reflect.Map:
		i = 4
		if len(b) < i {
			return 0, fmt.Errorf("%q requires %d bytes", rv.Kind(), i)
		}
		n := int(binary.LittleEndian.Uint32(b))
		rm := reflect.MakeMap(rv.Type())
		rm.Type().Key()
		for j := 0; j < n; j++ {
			rk := reflect.Indirect(reflect.New(rm.Type().Key()))
			x, err := deserializeModelReflect(b[i:], rk)
			if err != nil {
				return 0, err
			}
			i += x
			re := reflect.Indirect(reflect.New(reflect.Zero(rv.Type().Elem()).Type()))
			x, err = deserializeModelReflect(b[i:], re)
			if err != nil {
				return 0, err
			}
			i += x
			rm.SetMapIndex(rk, re)
		}
		rv.Set(rm)
	case reflect.Struct:
		rt := rv.Type()
		for n := 0; n < rv.NumField(); n++ {
			ok, m := parseStructFieldTag(rt.Field(n))
			if !ok {
				continue
			}
			if !m[modelStructFiledTagValueParam] {
				continue
			}
			rf := rv.Field(n)
			if m[modelStructFiledTagValueOption] {
				i += 1
				if len(b) < i {
					return 0, fmt.Errorf("option field requires %d bytes", 1)
				}
				if b[i-1] == 0 {
					rf.Set(reflect.Zero(rf.Type()))
					continue
				}
			}
			x, err := deserializeModelReflect(b[i:], rf)
			if err != nil {
				return 0, err
			}
			i += x
		}
	default:
		return 0, fmt.Errorf("unexpected value type %q", rv.Kind())
	}
	return i, nil
}

func deserializeModelTryCustom(b []byte, rv reflect.Value) (int, bool, error) {
	u, ok := reflect.New(rv.Type()).Interface().(ModelDeserializer)
	if !ok {
		return 0, false, nil
	}
	i, err := u.DeserializeModel(b)
	if err != nil {
		return 0, true, fmt.Errorf("custom deserializer: %w", err)
	}
	if !rv.CanSet() {
		return 0, true, fmt.Errorf("unsettable value")
	}
	rv.Set(reflect.ValueOf(u).Elem())
	return i, true, nil
}

func deserializeModelTryTime(b []byte, rv reflect.Value) (int, bool, error) {
	_, ok := rv.Interface().(time.Time)
	if !ok {
		return 0, false, nil
	}
	if len(b) < 8 {
		return 0, true, fmt.Errorf("time.Time requires 8 bytes")
	}
	m := int64(binary.LittleEndian.Uint64(b))
	if !rv.CanSet() {
		return 0, true, fmt.Errorf("unsettable value")
	}
	rv.Set(reflect.ValueOf(time.Unix(m/1000, m%1000*1e6)))
	return 8, true, nil
}

func deserializeModelTryDuration(b []byte, rv reflect.Value) (int, bool, error) {
	_, ok := rv.Interface().(time.Duration)
	if !ok {
		return 0, false, nil
	}
	if len(b) < 8 {
		return 0, true, fmt.Errorf("time.Duration requires 8 bytes")
	}
	m := int64(binary.LittleEndian.Uint64(b))
	if !rv.CanSet() {
		return 0, true, fmt.Errorf("unsettable value")
	}
	rv.SetInt(m * 1e6)
	return 8, true, nil
}
