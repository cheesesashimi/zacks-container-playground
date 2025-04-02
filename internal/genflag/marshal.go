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

// Begins the main marshaling process.
func (m marshaler) marshal(in interface{}) ([]Flag, error) {
	if in == nil {
		return nil, fmt.Errorf("cannot marshal flags from nil value")
	}

	val := reflect.ValueOf(in)
	typ := reflect.TypeOf(in)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	// If our input implements the FlagMarshaller interface, call it and return
	// the results.
	if m.isFlagMarshaller(typ) {
		return m.runFlagMarshaller(typ, val)
	}

	kind := typ.Kind()

	// If we have a struct, we need to examine its fields.
	if kind == reflect.Struct {
		return m.marshalStructFields(typ, val)
	}

	// If we have a top-level slice or array, return an error since there is
	// nothing we can do with it. Any slices or arrays must be contained within a
	// struct field.
	if kind == reflect.Slice || kind == reflect.Array {
		return m.marshalTopLevelSlice(val)
	}

	// If we have a map, convert it into keyValueFlags.
	if kind == reflect.Map {
		return m.marshalMap(nil, val)
	}

	return nil, fmt.Errorf("unknown kind %q", kind)
}

// Iterate through each struct field and retrieve any flags that
// are returned from it.
func (m marshaler) marshalStructFields(typ reflect.Type, val reflect.Value) ([]Flag, error) {
	cliflags := []Flag{}

	for i := 0; i < val.NumField(); i++ {
		f, err := m.marshalStructField(typ.Field(i), val.Field(i))
		if err != nil {
			return nil, fmt.Errorf("cannot parse struct field %q: %w", typ.Field(i).Name, err)
		}

		cliflags = append(cliflags, f...)
	}

	return cliflags, nil
}

// Determines if a given type implements the FlagMarshaller interface.
func (m marshaler) isFlagMarshaller(typ reflect.Type) bool {
	interfaceType := reflect.TypeOf((*FlagMarshaller)(nil)).Elem()
	return typ.Implements(interfaceType)
}

type stringer interface {
	String() string
}

func (m marshaler) isStringer(typ reflect.Type) bool {
	interfaceType := reflect.TypeOf((*stringer)(nil)).Elem()
	return typ.Implements(interfaceType)
}

func (m marshaler) runStringer(typ reflect.Type, val reflect.Value) (string, error) {
	if !m.isStringer(typ) {
		return "", fmt.Errorf("does not implement stringer")
	}

	methodValue := val.MethodByName("String")

	results := methodValue.Call(nil)

	if len(results) != 1 {
		return "", fmt.Errorf("expected 1 result, got %d", len(results))
	}

	if result, sliceOK := results[0].Interface().(string); sliceOK {
		return result, nil
	}

	return "", fmt.Errorf("a problem occured while calling String()")
}

// Calls the FlagMarshaller interface on a given type, if it
// implements it.
func (m marshaler) runFlagMarshaller(typ reflect.Type, val reflect.Value) ([]Flag, error) {
	if !m.isFlagMarshaller(typ) {
		return nil, fmt.Errorf("does not implement FlagMarshaller")
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

// Parses the field of a given struct.
func (m marshaler) marshalStructField(field reflect.StructField, val reflect.Value) ([]Flag, error) {
	if m.isNilPointer(val) {
		return nil, nil
	}

	val = m.getValue(val)

	kind := val.Kind()

	_, ok := field.Tag.Lookup(genFlagKeyName)
	if !ok {
		return m.marshalUntaggedStructField(field, val, kind)
	}

	return m.marshalTaggedStructField(field, val, kind)
}

// Handles untagged struct fields as well as discovering whether a given struct
// implements FlagMarshaller.
func (m marshaler) marshalUntaggedStructField(field reflect.StructField, val reflect.Value, kind reflect.Kind) ([]Flag, error) {
	// If the field's concrete tyoe implements FlagMarshaller, call
	// it directly and return the results.
	if m.isFlagMarshaller(val.Type()) {
		return m.runFlagMarshaller(val.Type(), val)
		// return m.marshal(val.Interface())
	}

	// If we found a struct, recursively handle it regardless of whether it has a
	// field tag or not. This allows us to examine nested structs which may or
	// may not have field tags.
	if kind == reflect.Struct {
		return m.marshal(val.Interface())
	}

	if (kind == reflect.Slice || kind == reflect.Array || kind == reflect.Map) && m.isFlagMarshaller(val.Type().Elem()) {
		return m.marshal(val.Interface())
	}

	return nil, nil
}

func (m marshaler) getFlagsFromFlagMarshallerSlice(val reflect.Value) ([]Flag, error) {
	flags := []Flag{}
	err := m.traverseSlice(val, func(v reflect.Value) error {
		f, err := m.marshal(v.Interface())
		if err != nil {
			return err
		}

		flags = append(flags, f...)

		return nil
	})

	return flags, err
}

// Handles the tagged struct fields.
func (m marshaler) marshalTaggedStructField(field reflect.StructField, val reflect.Value, kind reflect.Kind) ([]Flag, error) {
	tagContents, ok := field.Tag.Lookup(genFlagKeyName)
	if !ok {
		return nil, nil
	}

	gfo, err := newGenFlagOpts(field.Name, tagContents)
	if err != nil {
		return nil, err
	}

	if !field.IsExported() {
		return nil, fmt.Errorf("cannot reflect from unexported field")
	}

	if m.isStringer(val.Type()) {
		r, err := m.runStringer(val.Type(), val)
		if err != nil {
			return nil, err
		}

		return gfo.newStringFlagWithName(r)
	}

	if kind == reflect.Struct {
		return m.marshal(val.Interface())
	}

	if kind == reflect.String && val.String() != "" {
		return gfo.newStringFlagWithName(val.String())
	}

	if kind == reflect.Bool {
		return gfo.newBoolFlagWithName(val.Bool())
	}

	if kind == reflect.Map {

		return m.marshalMap(gfo, val)
	}

	if kind == reflect.Array || kind == reflect.Slice {
		return m.marshalSlice(gfo, "", val)
	}

	if kind == reflect.Interface {
		return m.marshal(val.Interface())
	}

	return nil, nil
}

// Retrieves the string slice values from a given struct field.
func (m marshaler) fetchStringSliceValues(val reflect.Value) ([]string, error) {
	items := []string{}

	for i := 0; i < val.Len(); i++ {
		sliceVal := val.Index(i)
		if m.isNilPointer(sliceVal) {
			continue
		}

		sliceVal = m.getValue(sliceVal)

		kind := sliceVal.Kind()

		if kind != reflect.String {
			return nil, fmt.Errorf("unsupported slice kind %s", kind)
		}

		items = append(items, sliceVal.String())
	}

	return items, nil
}

func (m marshaler) newStringFlag(gfo *genFlagOpts, name, value string) ([]Flag, error) {
	if gfo == nil && name == "" {
		return nil, fmt.Errorf("when genFlagOpts is not passed, a name is required")
	}

	if gfo != nil {
		if name == "" {
			return gfo.newStringFlagWithName(value)
		}

		return gfo.newStringFlag(name, value)
	}

	f, err := NewStringFlag(name, value)
	return []Flag{f}, err
}

// Fetches the string slice values from a given struct field and passes them
// into the listflag constructor either by passing them through genFlagOpts to
// get the field name and optionFuncs or by calling the constructor directly.
func (m marshaler) newListFlag(gfo *genFlagOpts, name string, items []string) ([]Flag, error) {
	if gfo != nil {
		return gfo.newListFlag(items)
	}

	if gfo == nil && name == "" {
		return nil, fmt.Errorf("when genFlagOpts is not passed, a name is required")
	}

	return NewListFlags(name, items)
}

// Passes the provided string map into the NewKeyValueFlags constructor either
// by passing them through genFlagOpts to get the field name and optionFuncs or
// by calling the constructor directly.
func (m marshaler) newKeyValueFlag(gfo *genFlagOpts, vals map[string]string) ([]Flag, error) {
	if gfo != nil {
		return gfo.newKeyValueFlag(vals)
	}

	return NewKeyValueFlags(vals)
}

// Passes the provided string map into the NewSwitchFlag constructor either by
// passing them through genFlagOpts to get the field name and optionFuncs or by
// calling the constructor directly.
func (m marshaler) newSwitchFlag(gfo *genFlagOpts, vals map[string]bool) ([]Flag, error) {
	if gfo != nil {
		return gfo.newSwitchFlag(vals)
	}

	return NewSwitchFlags(vals)
}

func (m marshaler) marshalTopLevelSlice(val reflect.Value) ([]Flag, error) {
	elemType := val.Type().Elem()
	elemKind := elemType.Kind()

	kinds := map[reflect.Kind]struct{}{
		reflect.Struct:    struct{}{},
		reflect.Interface: struct{}{},
		reflect.Pointer:   struct{}{},
	}

	if _, ok := kinds[elemKind]; !ok {
		return nil, fmt.Errorf("slices of %s not supported at top-level", elemKind)
	}

	flags := []Flag{}

	for i := 0; i < val.Len(); i++ {
		sliceVal := m.getValue(val.Index(i))
		f, err := m.marshal(sliceVal.Interface())
		if err != nil {
			return nil, err
		}

		flags = append(flags, f...)
	}

	return flags, nil
}

func (m marshaler) traverseSlice(val reflect.Value, f func(reflect.Value) error) error {
	for i := 0; i < val.Len(); i++ {
		if err := f(val.Index(i)); err != nil {
			return err
		}
	}

	return nil
}

// Converts a slice into flags.
func (m marshaler) marshalSlice(gfo *genFlagOpts, name string, val reflect.Value) ([]Flag, error) {
	flags := []Flag{}

	items := []string{}

	err := m.traverseSlice(val, func(v reflect.Value) error {
		if m.isNilPointer(v) {
			return nil
		}

		if m.isStringer(v.Type()) {
			r, err := m.runStringer(v.Type(), v)
			if err != nil {
				return err
			}

			items = append(items, r)
			return nil
		}

		v = m.getValue(v)
		kind := v.Kind()

		if kind == reflect.Struct || kind == reflect.Interface {
			if m.isStringer(v.Type()) {
				r, err := m.runStringer(v.Type(), v)
				if err != nil {
					return err
				}

				items = append(items, r)
				return nil
			}

			f, err := m.marshal(v.Interface())
			if err != nil {
				return err
			}

			flags = append(flags, f...)
			return nil
		}

		if kind == reflect.String {
			items = append(items, v.String())
			return nil
		}

		if kind == reflect.Bool {
			return fmt.Errorf("bool slices are not allowed")
		}

		if kind == reflect.Slice {
			return fmt.Errorf("slices of slices not allowed")
		}

		return fmt.Errorf("unexpected kind %s", kind)
	})

	if err != nil {
		return nil, err
	}

	f, err := m.newListFlag(gfo, name, items)
	if err != nil {
		return nil, err
	}

	return append(flags, f...), nil
}

// Determines if a given value is actually a nil pointer.
func (m marshaler) isNilPointer(val reflect.Value) bool {
	kind := val.Kind()

	if kind != reflect.Pointer && kind != reflect.Ptr {
		return false
	}

	return val.IsNil()
}

// Retrieves a given value, recursing if necessary to both dereference the
// value if it's a pointer as well as fetch the concrete type if it's an
// interface.
func (m marshaler) getValue(val reflect.Value) reflect.Value {
	kind := val.Kind()

	if kind == reflect.Interface {
		return m.getValue(reflect.ValueOf(val.Interface()))
	}

	if kind == reflect.Pointer || kind == reflect.Ptr {
		return m.getValue(val.Elem())
	}

	return val
}

// Iterates through a map and calls the supplied function once for each
// iteration until the map has been iterated or an error has been returned.
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

// Iterates through a map and calls the supplied function once for each
// iteration until the map has been iterated or an error has been returned. The
// provided function can return a list of flags which will be combined into a
// singular list of flags once this function finishes.
func (m marshaler) traverseMapAndGetFlags(val reflect.Value, f func(reflect.Value, reflect.Value) ([]Flag, error)) ([]Flag, error) {
	flags := []Flag{}

	err := m.traverseMap(val, func(k, v reflect.Value) error {
		flagsOut, err := f(k, v)

		if err != nil {
			return err
		}

		flags = append(flags, flagsOut...)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return flags, nil
}

// Marshals a given map into switchFlags or keyValueFlags based upon the
// underlying type of the map. Also handles maps with pointers to strings and
// bools as well as map[string]interface{} where all four of those combinations
// can be present.
func (m marshaler) marshalMap(gfo *genFlagOpts, val reflect.Value) ([]Flag, error) {
	boolMap := map[string]bool{}
	stringMap := map[string]string{}

	flags, err := m.traverseMapAndGetFlags(val, func(k, v reflect.Value) ([]Flag, error) {
		if k.Kind() != reflect.String {
			return nil, fmt.Errorf("map key is %s, expected string", k.Kind())
		}

		if m.isNilPointer(v) {
			return nil, nil
		}

		strVal := m.getValue(v)

		if m.isStringer(strVal.Type()) {
			r, err := m.runStringer(strVal.Type(), strVal)
			if err != nil {
				return nil, err
			}

			return m.newStringFlag(gfo, k.String(), r)
		}

		val := m.getValue(v)

		kind := val.Kind()

		if kind == reflect.Interface {
			return m.marshal(val.Interface())
		}

		if kind == reflect.Map {
			return m.marshal(val.Interface())
		}

		if kind == reflect.Struct {
			return m.marshal(val.Interface())
		}

		if kind == reflect.Slice || kind == reflect.Array {
			return m.marshalSlice(gfo, k.String(), val)
		}

		if kind == reflect.String {
			stringMap[k.String()] = val.String()
		}

		if kind == reflect.Bool {
			boolMap[k.String()] = val.Bool()
		}

		if kind != reflect.String && kind != reflect.Bool {
			return nil, fmt.Errorf("invalid map value type %q", kind)
		}

		return nil, nil
	})

	getCollectedFlags := func() ([]Flag, error) {
		// While its possible to call the map traversal function directly from here,
		// it is better that we execute it above and have this closure just return
		// its values.
		//
		// The reason is because getBoolFlags() and getKVFlags() are dependent upon
		// work done when the map traversal function executes. If that work is not
		// done before the other two closures are called, there will be no data to
		// return.
		//
		// By calling the map traversal function above and passing its return values
		// via this closure, we can guarantee that it will execute first.
		return flags, err
	}

	getBoolFlags := func() ([]Flag, error) {
		return m.newSwitchFlag(gfo, boolMap)
	}

	getKVFlags := func() ([]Flag, error) {
		return m.newKeyValueFlag(gfo, stringMap)
	}

	return combineFlags(getCollectedFlags, getBoolFlags, getKVFlags)
}
