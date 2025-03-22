package genflag

import (
	"fmt"
	"reflect"
)

// Marshals the given interface into a list of validated flags,
// halting on any errors.
func MarshalFlags(in interface{}) ([]Flag, error) {
	m := marshaler{}

	flags, err := m.marshal(in)
	if err != nil {
		return nil, err
	}

	if err := m.validateFlagList(flags); err != nil {
		return nil, err
	}

	return flags, nil
}

type marshaler struct{}

// Ensures that flag names and values are unique according to the following rules:
//
// 1. Names must be unique across all different types of flags
// unless they are listflags. Listflags will be compared to the
// rest of the flags in order to determine whether a collision
// occurs.
// 2. Values can be different amongst different flags.
func (m marshaler) validateFlagList(flags []Flag) error {
	seenCombined := map[string]struct{}{}
	seenBooleanFlags := map[string]struct{}{}
	seenStringFlags := map[string]struct{}{}

	for _, flag := range flags {
		name := flag.Name()
		val := ""

		isBoolFlag := false

		switch c := flag.(type) {
		case stringFlag:
			val = c.value
		case boolFlag:
			val = fmt.Sprintf("%v", c.value)
			isBoolFlag = true
		case listFlag:
			val = c.value
		case keyValueFlag:
			val = c.value
		case switchFlag:
			val = fmt.Sprintf("%v", c.value)
			isBoolFlag = true
		default:
			val = flag.Value()
		}

		_, seenBoolOK := seenBooleanFlags[name]
		_, seenStringOK := seenStringFlags[name]

		if isBoolFlag {
			if seenStringOK {
				return fmt.Errorf("flag name collision")
			}

			if !seenBoolOK {
				seenBooleanFlags[name] = struct{}{}
			}
		} else {
			if seenBoolOK {
				return fmt.Errorf("flag name collision")
			}
		}

		if !seenStringOK {
			seenStringFlags[name] = struct{}{}
		}

		key := fmt.Sprintf("%s/%s", name, val)

		_, ok := seenCombined[key]
		if !ok {
			seenCombined[key] = struct{}{}
		} else {
			return fmt.Errorf("flags cannot have the same name / value, found: %q", key)
		}
	}

	return nil
}

// Dereferences any pointer values that were given.
func (m marshaler) handlePointers(in interface{}) (reflect.Type, reflect.Value) {
	val := reflect.ValueOf(in)
	typ := reflect.TypeOf(in)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	return typ, val
}

// Begins the main marshaling process.
func (m marshaler) marshal(in interface{}) ([]Flag, error) {
	if in == nil {
		return nil, fmt.Errorf("cannot marshal flags from nil value")
	}

	typ, val := m.handlePointers(in)

	if m.isMarshalFlags(typ) {
		return m.runMarshalFlags(typ, val)
	}

	kind := typ.Kind()

	// If we have a struct, we need to examine its fields.
	if kind == reflect.Struct {
		return m.marshalStructFields(typ, val)
	}

	// If we have a top-level slice or array, return an error since
	// there is nothing we can do with it.
	if kind == reflect.Slice || kind == reflect.Array {
		return nil, fmt.Errorf("top-level slices / arrays not supported")
	}

	// If we have a map, convert it into keyValueFlags.
	if kind == reflect.Map {
		return m.getFlagMap(nil, val)
	}

	return nil, fmt.Errorf("unknown kind %q", kind)
}

// Iterate through each struct field and retrieve any flags that
// are returned from it.
func (m marshaler) marshalStructFields(typ reflect.Type, val reflect.Value) ([]Flag, error) {
	cliflags := []Flag{}

	for i := 0; i < val.NumField(); i++ {
		f, err := m.parseField(typ.Field(i), val.Field(i))
		if err != nil {
			return nil, fmt.Errorf("cannot parse struct field %q: %w", typ.Field(i).Name, err)
		}

		cliflags = append(cliflags, f...)
	}

	return cliflags, nil
}

// Determines if a given type implements the FlagMarshaler interface.
func (m marshaler) isMarshalFlags(typ reflect.Type) bool {
	interfaceType := reflect.TypeOf((*FlagMarshaler)(nil)).Elem()
	return typ.Implements(interfaceType)
}

// Calls the FlagMarshaler interface on a given type, if it
// implements it.
func (m marshaler) runMarshalFlags(typ reflect.Type, val reflect.Value) ([]Flag, error) {
	if !m.isMarshalFlags(typ) {
		return nil, fmt.Errorf("does not implement FlagMarshaler")
	}

	methodValue := val.MethodByName("MarshalFlags")

	results := methodValue.Call(nil)

	if len(results) != 2 {
		return nil, fmt.Errorf("expected 2 results, got %d", len(results))
	}

	if resultSlice, sliceOK := results[0].Interface().([]Flag); sliceOK {
		resultErr, errOK := results[1].Interface().(error)
		if errOK || results[1].IsNil() {
			return resultSlice, resultErr
		}

		if errOK && !results[1].IsNil() {
			return resultSlice, resultErr
		}

		return nil, fmt.Errorf("expected an error value, got: %q", results[1].Interface())
	}

	return nil, fmt.Errorf("expected a slice of Flags, got: %q", results[1].Interface())
}

func (m marshaler) isStructTagged(typ reflect.Type, val reflect.Value) bool {
	// Dereference any pointers.
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			// If the pointer is nil, there is nothing to do.
			return false
		}
	}

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		_, ok := field.Tag.Lookup(genFlagKeyName)
		if ok {
			return true
		}
	}

	return false
}

// Parses the field of a given struct.
func (m marshaler) parseField(field reflect.StructField, val reflect.Value) ([]Flag, error) {
	// Dereference any pointers.
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			// If the pointer is nil, there is nothing to do.
			return nil, nil
		}

		val = val.Elem()
	}

	// If the field's concrete tyoe implements FlagMarshaler, call
	// it directly and return the results.
	if m.isMarshalFlags(val.Type()) {
		return m.runMarshalFlags(val.Type(), val)
	}

	kind := val.Kind()

	// If we found a struct, recursively handle it regardless of
	// whether the field is tagged or not.
	// TODO: Determine if this is the desired behavior or not.
	if kind == reflect.Struct {
		return m.marshal(val.Interface())
	}

	tagContents, ok := field.Tag.Lookup(genFlagKeyName)
	if !ok {
		return nil, nil
	}

	if !field.IsExported() {
		return nil, fmt.Errorf("cannot reflect from unexported field")
	}

	gfo, err := newGenFlagOpts(field.Name, tagContents)
	if err != nil {
		return nil, err
	}

	// If we found a struct, recursively handle it.
	// if kind == reflect.Struct {
	// 	return m.marshal(val.Interface())
	// }

	if kind == reflect.String && val.String() != "" {
		return toPlural(gfo.newStringFlagWithName(val.String()))
	}

	if kind == reflect.Bool {
		return toPlural(gfo.newBoolFlagWithName(val.Bool()))
	}

	if kind == reflect.Map {
		return m.getFlagMap(gfo, val)
	}

	if kind == reflect.Array || kind == reflect.Slice {
		if reflect.TypeOf(val.Interface()).Elem().Kind() == reflect.Struct {
			return nil, fmt.Errorf("listed nested structs not supported")
		}
		return m.getFlagSlice(gfo, val)
	}

	if kind == reflect.Interface {
		return m.parseField(field, reflect.ValueOf(val.Interface()))
	}

	return nil, nil
}

// Traverses a slice or array and calls the given function on each value found within.
func (m marshaler) traverseSliceOrArray(val reflect.Value, f func(reflect.Value) error) error {
	for i := 0; i < val.Len(); i++ {
		if err := f(val.Index(i)); err != nil {
			return err
		}
	}

	return nil
}

func (m marshaler) getFlagSlice(gfo *genFlagOpts, val reflect.Value) ([]Flag, error) {
	items := []string{}

	flags := []Flag{}

	err := m.traverseSliceOrArray(val, func(sliceVal reflect.Value) error {
		kind := sliceVal.Kind()

		if kind == reflect.Pointer {
			sliceVal = sliceVal.Elem()
		}

		if kind == reflect.Struct {
			return fmt.Errorf("list of structs not supported")
		}

		s, b, err := m.getMapOrSliceElement(sliceVal)
		if err != nil {
			return err
		}

		if b != nil {
			return fmt.Errorf("unexpected bool")
		}

		items = append(items, *s)

		return nil
	})

	if err != nil {
		return nil, err
	}

	f, err := gfo.newListFlag(items)
	if err != nil {
		return nil, err
	}

	return append(flags, f...), nil
}

func (m marshaler) getMapOrSliceElement(val reflect.Value) (*string, *bool, error) {
	switch val.Kind() {
	case reflect.String:
		s := val.String()
		return &s, nil, nil
	case reflect.Bool:
		b := val.Bool()
		return nil, &b, nil
	case reflect.Interface:
		return m.getMapOrSliceElement(reflect.ValueOf(val.Interface()))
	case reflect.Pointer:
		return m.getMapOrSliceElement(val.Elem())
	default:
		return nil, nil, fmt.Errorf("invalid kind %s", val.Kind())
	}
}

func (m marshaler) traverseMap(val reflect.Value, f func(reflect.Value, reflect.Value) error) error {
	iter := val.MapRange()

	for iter.Next() {
		k := iter.Key()
		v := iter.Value()

		if k.Kind() != reflect.String {
			return fmt.Errorf("map key is %s, expected string", k.Kind())
		}

		if err := f(k, v); err != nil {
			return err
		}
	}

	return nil
}

func (m marshaler) getFlagsFromStringMap(gfo *genFlagOpts, items map[string]string) ([]Flag, error) {
	if gfo != nil {
		return gfo.newKeyValueFlag(items)
	}

	return NewKeyValueFlags(items)
}

func (m marshaler) getFlagsFromBoolMap(gfo *genFlagOpts, items map[string]bool) ([]Flag, error) {
	if gfo != nil {
		return gfo.newSwitchFlag(items)
	}

	return NewSwitchFlags(items)
}

func (m marshaler) getFlagMap(gfo *genFlagOpts, val reflect.Value) ([]Flag, error) {
	boolMap := map[string]bool{}
	stringMap := map[string]string{}

	out := []Flag{}

	err := m.traverseMap(val, func(k, v reflect.Value) error {
		s, b, err := m.getMapOrSliceElement(v)
		if err != nil {
			return err
		}

		if s != nil {
			stringMap[k.String()] = *s
		}

		if b != nil {
			boolMap[k.String()] = *b
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(stringMap) != 0 {
		f, err := m.getFlagsFromStringMap(gfo, stringMap)
		if err != nil {
			return nil, err
		}

		out = append(out, f...)
	}

	if len(boolMap) != 0 {
		f, err := m.getFlagsFromBoolMap(gfo, boolMap)
		if err != nil {
			return nil, err
		}

		out = append(out, f...)
	}

	return out, nil
}

func stringMapToSlice(in map[string]struct{}) []string {
	out := []string{}

	for key := range in {
		out = append(out, key)
	}

	return out
}
