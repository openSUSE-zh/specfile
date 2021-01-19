package specfile

import (
	"errors"
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

// FindTag find a tag in specfile
func (s Specfile) FindTag(name string) (Tag, error) {
	for _, t := range s.Tags {
		if t.Name == name {
			return t, nil
		}
	}
	return Tag{}, errors.New("tag not found")
}

// FindSection find a section in specfile
func (s Specfile) FindSection(name string) (Section, error) {
	for _, section := range s.Sections {
		if section.Name == name {
			return section, nil
		}
	}
	return Section{}, errors.New("section not found")
}

// FindSubpackage find a subpackage in specfile
func (s Specfile) FindSubpackage(name string) (Specfile, error) {
	for _, spec := range s.Subpackages {
		if tag, err := spec.FindTag("Name"); err == nil {
			if tag.Value == name {
				return spec, nil
			}
		}
	}
	return Specfile{}, errors.New("specfile not found")
}

// append append value to fields
func (s *Specfile) append(field string, val interface{}) {
	sv := reflect.ValueOf(s).Elem().FieldByName(field)
	if sv.Len() > 0 {
		sv.Set(reflect.Append(sv, reflect.ValueOf(val)))
	} else {
		nv := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(val)), 1, 1)
		nv.Index(0).Set(reflect.ValueOf(val))
		sv.Set(nv)
	}
}
