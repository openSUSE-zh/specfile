package specfile

import (
	"reflect"
)

// Specfile specfile struct
type Specfile struct {
	Subpackages  []Specfile
	Tags         []Tag
	Macros       Macros
	Sections     []Section
	Dependencies []Dependency
}

// append append value to fields
func (s *Specfile) append(fld string, val interface{}) {
	sv := reflect.ValueOf(s).Elem().FieldByName(fld)
	if sv.Len() > 0 {
		sv.Set(reflect.Append(sv, reflect.ValueOf(val)))
	} else {
		nv := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(val)), 1, 1)
		nv.Index(0).Set(reflect.ValueOf(val))
		sv.Set(nv)
	}
}
