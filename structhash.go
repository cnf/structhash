package structhash

import (
	"crypto/md5"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// Version returns the version of the supplied hash as an integer
// or -1 on failure
func Version(h string) int {
	if h == "" {
		return -1
	}
	if h[0] != 'v' {
		return -1
	}
	if spos := strings.IndexRune(h[1:], '_'); spos >= 0 {
		n, e := strconv.Atoi(h[1 : spos+1])
		if e != nil {
			return -1
		}
		return n
	}
	return -1
}

// Hash takes a data structure and returns a hash string of that data structure
// at the version asked
func Hash(c interface{}, version int) (string, error) {
	serial := serialize(c, version)
	return fmt.Sprintf("v%d_%x", version, md5.Sum([]byte(serial))), nil
}

type structFieldFilter func(reflect.StructField) bool

func strValue(val reflect.Value, depth int, fltr structFieldFilter) string {
	switch val.Kind() {
	case reflect.String:
		return "\"" + val.String() + "\""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", val.Int())
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%d", val.Uint())
	case reflect.Bool:
		if val.Bool() {
			return "true"
		}
		return "false"
	case reflect.Ptr:
		return strValue(reflect.Indirect(val), depth, fltr)
	case reflect.Array, reflect.Slice:
		len := val.Len()
		ret := "["
		for i := 0; i < len; i++ {
			if i != 0 {
				ret = ret + ", "
			}
			ret = ret + strValue(val.Index(i), depth+1, fltr)
		}
		ret = ret + "]"
		return ret
	case reflect.Map:
		len := val.Len()
		mk := val.MapKeys()
		strmap := make(map[string]string, len)
		// Map/serialize all values
		for i := 0; i < len; i++ {
			strmap[strValue(mk[i], depth+1, fltr)] = strValue(val.MapIndex(mk[i]), depth+1, fltr)
		}

		// Create array to hold map keys and sort them
		skey := make([]string, 0, len)
		for k := range strmap {
			skey = append(skey, k)
		}
		sort.Strings(skey)

		ret := "["
		for i := 0; i < len; i++ {
			if i != 0 {
				ret = ret + ", "
			}
			ret = ret + skey[i] + ": " + strmap[skey[i]]
		}
		ret = ret + "]"
		return ret
	case reflect.Struct:

		vtype := val.Type()
		flen := vtype.NumField()
		smap := make(map[string]string)
		// Get all fields
		for i := 0; i < flen; i++ {
			field := vtype.Field(i)
			if (fltr != nil) && (!fltr(field)) {
				continue
			}
			fval := val.Field(i)
			smap[field.Name] = strValue(fval, depth+1, fltr)
		}
		// Get the keys and sort them
		skey := make([]string, len(smap))
		c := 0
		for k := range smap {
			skey[c] = k
			c++
		}
		sort.Strings(skey)

		flen = len(skey)
		ret := "{\n"
		for i := 0; i < flen; i++ {
			if i != 0 {
				ret = ret + ", \n"
			}
			ret = ret + strings.Repeat("  ", depth+1) + "\"" + skey[i] + "\": "
			ret = ret + smap[skey[i]]
		}
		return ret + "\n" + strings.Repeat("  ", depth) + "}"
	default:
		return val.String()
	}
}

// func serializeFilter(object interface{}, fltr structFieldFilter) string {
// 	return strValue(reflect.ValueOf(object), 0, fltr) + "\n"
// }

func serialize(object interface{}, version int) string {
	return strValue(reflect.ValueOf(object), 0, func(f reflect.StructField) bool {
		var err error
		ver := 0
		if lastver, err := strconv.Atoi(f.Tag.Get("lastversion")); err == nil {
			if lastver < version {
				return false
			}
		}
		if ver, err = strconv.Atoi(f.Tag.Get("version")); err != nil {
			return false
		}
		if ver <= version {
			return true
		}
		return false
	}) + "\n"
}
