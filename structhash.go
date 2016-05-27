package structhash

import (
	"crypto/md5"
	"crypto/sha1"
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
// at the version asked.
//
// This function uses md5 hashing function and default formatter. See also Dump()
// function.
func Hash(c interface{}, version int) (string, error) {
	return fmt.Sprintf("v%d_%x", version, Md5(c, version)), nil
}

// Dump takes a data structure and returns its byte representation. This can be
// useful if you need to use your own hashing function or formatter.
func Dump(c interface{}, version int) []byte {
	return []byte(serialize(c, version))
}

// Md5 takes a data structure and returns its md5 hash.
// This is a shorthand for md5.Sum(Dump(c, version)).
func Md5(c interface{}, version int) []byte {
	sum := md5.Sum(Dump(c, version))
	return sum[:]
}

// Sha1 takes a data structure and returns its sha1 hash.
// This is a shorthand for sha1.Sum(Dump(c, version)).
func Sha1(c interface{}, version int) []byte {
	sum := sha1.Sum(Dump(c, version))
	return sum[:]
}

type structFieldFilter func(reflect.StructField) (string, bool)

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
		if !val.IsNil() || val.Type().Elem().Kind() == reflect.Struct {
			return strValue(reflect.Indirect(val), depth, fltr)
		} else {
			return strValue(reflect.Zero(val.Type().Elem()), depth, fltr)
		}
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
			strmap[strValue(mk[i], depth+1, fltr)] =
				strValue(val.MapIndex(mk[i]), depth+1, fltr)
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
			name := field.Name
			if fltr != nil {
				if newname, ok := fltr(field); ok {
					name = newname
				} else {
					continue
				}
			}
			fval := val.Field(i)
			smap[name] = strValue(fval, depth+1, fltr)
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

func serialize(object interface{}, version int) string {
	return strValue(reflect.ValueOf(object), 0, func(f reflect.StructField) (string, bool) {
		var err error
		name := f.Name
		ver := 0
		lastver := -1
		if str := f.Tag.Get("lastversion"); str != "" {
			if lastver, err = strconv.Atoi(str); err != nil {
				return "", false
			}
		}
		if str := f.Tag.Get("version"); str != "" {
			if ver, err = strconv.Atoi(str); err != nil {
				return "", false
			}
		}
		if str := f.Tag.Get("hash"); str != "" {
			parts := strings.Split(str, ",")
			if len(parts) > 0 {
				n := strings.TrimSpace(parts[0])
				if n == "-" {
					return "", false
				} else if n != "" {
					name = n
				}
			}
			if len(parts) > 1 {
				for _, tag := range strings.Split(parts[1], " ") {
					tag = strings.TrimSpace(tag)
					if strings.HasPrefix(tag, "version(") {
						arg := strings.TrimPrefix(tag, "version(")
						arg = strings.TrimSuffix(arg, ")")
						if ver, err = strconv.Atoi(arg); err != nil {
							return "", false
						}
					} else if strings.HasPrefix(tag, "lastversion(") {
						arg := strings.TrimPrefix(tag, "lastversion(")
						arg = strings.TrimSuffix(arg, ")")
						if lastver, err = strconv.Atoi(arg); err != nil {
							return "", false
						}
					}
				}
			}
		}
		if lastver != -1 && lastver < version {
			return "", false
		}
		if ver > version {
			return "", false
		}
		return name, true
	}) + "\n"
}
