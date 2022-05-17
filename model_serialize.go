package concordium

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"time"
)

func serializeModel(v any) ([]byte, error) {
	return serializeModelReflect(reflect.Indirect(reflect.ValueOf(v)))
}

func serializeModelReflect(rv reflect.Value) ([]byte, error) {
	for _, h := range []func(reflect.Value) ([]byte, bool, error){
		serializeModelTryCustom,
		serializeModelTryTime,
		serializeModelTryDuration,
	} {
		if b, ok, err := h(rv); ok {
			return b, err
		}
	}
	var b []byte
	switch rv.Kind() {
	case reflect.Interface:
		return serializeModelReflect(reflect.Indirect(reflect.ValueOf(rv.Interface())))
	case reflect.Ptr:
		return serializeModelReflect(reflect.Indirect(rv))
	case reflect.Invalid:
		b = []byte{}
	case reflect.Bool:
		if rv.Bool() {
			b = []byte{1}
		} else {
			b = []byte{0}
		}
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		// Converting between int64 and uint64 doesn't change the sign bit,
		// only the way it's interpreted. So let's do that with magic below.
		var i int
		var h func([]byte, uint64)
		switch rv.Kind() {
		case reflect.Int64, reflect.Uint64:
			i = 8
			h = binary.LittleEndian.PutUint64
		case reflect.Int32, reflect.Uint32:
			i = 4
			h = func(b []byte, u uint64) { binary.LittleEndian.PutUint32(b, uint32(u)) }
		case reflect.Int16, reflect.Uint16:
			i = 2
			h = func(b []byte, u uint64) { binary.LittleEndian.PutUint16(b, uint16(u)) }
		default:
			i = 1
			h = func(b []byte, u uint64) { b[0] = uint8(u) }
		}
		b = make([]byte, i)
		if rv.Kind() == reflect.Int8 || rv.Kind() == reflect.Int16 || rv.Kind() == reflect.Int32 || rv.Kind() == reflect.Int64 {
			h(b, uint64(rv.Int()))
		} else {
			h(b, rv.Uint())
		}
	case reflect.String:
		s := rv.String()
		c := rv.Len()
		b = make([]byte, 4+c)
		binary.LittleEndian.PutUint32(b, uint32(c))
		copy(b[4:], s)
	case reflect.Slice:
		i := 4
		c := rv.Len()
		bs := make([][]byte, c)
		for j := 0; j < c; j++ {
			var err error
			bs[j], err = serializeModelReflect(rv.Index(j))
			if err != nil {
				return nil, err
			}
			i += len(bs[j])
		}
		b = make([]byte, i)
		binary.LittleEndian.PutUint32(b, uint32(c))
		z := 4
		for _, x := range bs {
			copy(b[z:], x)
			z += len(x)
		}
	case reflect.Map:
		i := 4
		c := rv.Len()
		bs := make([][]byte, c*2)
		var j int
		for _, rk := range rv.MapKeys() {
			for _, v := range []reflect.Value{rk, rv.MapIndex(rk)} {
				var err error
				bs[j], err = serializeModelReflect(v)
				if err != nil {
					return nil, err
				}
				i += len(bs[j])
				j++
			}
		}
		b = make([]byte, i)
		binary.LittleEndian.PutUint32(b, uint32(c))
		z := 4
		for _, x := range bs {
			z += copy(b[z:], x)
		}
	case reflect.Struct:
		var i int
		var bs [][]byte
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
			be, err := serializeModelReflect(rf)
			if err != nil {
				return nil, err
			}
			if m[modelStructFiledTagValueOption] {
				if len(be) == 0 {
					bs = append(bs, []byte{0})
				} else {
					bs = append(bs, []byte{1})
				}
				i++
			}
			bs = append(bs, be)
			i += len(be)
		}
		b = make([]byte, i)
		var z int
		for _, x := range bs {
			z += copy(b[z:], x)
		}
	default:
		return nil, fmt.Errorf("unexpected value type %q", rv.Kind())
	}
	return b, nil
}

func serializeModelTryCustom(rv reflect.Value) ([]byte, bool, error) {
	if !rv.IsValid() {
		return nil, false, nil
	}
	// this is required for process non-pointer value
	if rv.Kind() != reflect.Ptr {
		rvn := reflect.New(rv.Type())
		rvi := reflect.Indirect(rvn)
		if rvi.CanSet() {
			rvi.Set(rv)
			rv = rvn
		}
	}
	u, ok := rv.Interface().(SerializeModel)
	if !ok {
		return nil, false, nil
	}
	// avoid 'invalid memory address or nil pointer dereference' error
	if rv.IsNil() {
		return nil, true, nil
	}
	b, err := u.SerializeModel()
	if err != nil {
		return nil, true, fmt.Errorf("custom serializer: %w", err)
	}
	return b, true, nil
}

func serializeModelTryTime(rv reflect.Value) ([]byte, bool, error) {
	if !rv.IsValid() {
		return nil, false, nil
	}
	t, ok := rv.Interface().(time.Time)
	if !ok {
		return nil, false, nil
	}
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(t.UnixNano()/1e6))
	return b, true, nil
}

func serializeModelTryDuration(rv reflect.Value) ([]byte, bool, error) {
	if !rv.IsValid() {
		return nil, false, nil
	}
	d, ok := rv.Interface().(time.Duration)
	if !ok {
		return nil, false, nil
	}
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(d/1e6))
	return b, true, nil
}
