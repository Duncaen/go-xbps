package repodata

import (
	"archive/tar"
	"bytes"
	"io"
	"strings"
	"sync"
	"reflect"
	"path"

	"howett.net/plist"

	"github.com/klauspost/compress/zstd"
)

type Format int

const (
	FormatUnknown Format = iota
	FormatTar
)

var zstdFrameMagic = []byte{0x28, 0xb5, 0x2f, 0xfd}

type Decoder struct {
	reader io.ReadSeeker
}

type invalidError struct {
	format string
	err    error
}

func (e invalidError) Error() string {
	if e.err != nil {
		return strings.Join([]string{"repodata: invalid ", e.format, ": ", e.err.Error()}, "")
	}
	return strings.Join([]string{"repodata: invalid ", e.format}, "")
}

type field struct {
	name string
	omitEmpty bool
	index []int
	typ reflect.Type
}

type structFields struct {
	list []field
	nameIndex map[string]int
}

var fieldCache sync.Map // map[reflect.Type]structFields

func typeFields(t reflect.Type) structFields {
	current := []field{}
	next := []field{{typ: t}}

	var count, nextCount map[reflect.Type]int

	visited := map[reflect.Type]struct{}{}

	var fields []field

	for len(next) > 0 {
		current, next = next, current[:0]
		count, nextCount = nextCount, map[reflect.Type]int{}
		for _, f := range current {
			if _, ok := visited[f.typ]; ok {
				continue
			}
			visited[f.typ] = struct{}{}

			for i := 0; i < f.typ.NumField(); i++ {
				sf := f.typ.Field(i)
				if sf.Anonymous {
					t := sf.Type
					if t.Kind() == reflect.Pointer {
						t = t.Elem()
					}
					if !sf.IsExported() && t.Kind() != reflect.Struct {
						// Ignore embedded fields of unexported non-struct types.
						continue
					}
					// Do not ignore embedded fields of unexported struct types
					// since they may have exported fields.
				} else if !sf.IsExported() {
					// Ignore unexported non-embedded fields.
					continue
				}
				tag := sf.Tag.Get("repodata")
				if tag == "-" {
					// skip field
					continue
				}
				name, opts, _ := strings.Cut(tag, ",")

				index := make([]int, len(f.index)+1)
				copy(index, f.index)
				index[len(f.index)] = i

				ft := sf.Type
				if ft.Name() == "" && ft.Kind() == reflect.Pointer {
					ft = ft.Elem()
				}
				if name != "" {
					// only support tagged fields
					field := field{
						name: name,
						index: index,
					}
					for _, opt := range strings.Split(opts, ",") {
						switch opt {
						case "omitempty":
							field.omitEmpty = true
						}
					}
					fields = append(fields, field)
					if count[f.typ] > 1 {
						fields = append(fields, fields[len(fields)-1])
					}
					continue
				}

				// Record new anonymous struct to explore in next round.
				nextCount[ft]++
				if nextCount[ft] == 1 {
					next = append(next, field{name: ft.Name(), index: index, typ: ft})
				}
			}
		}
	}
	// XXX: handle conflicts
	nameIndex := make(map[string]int, len(fields))
	for i, field := range fields {
		nameIndex[field.name] = i
	}
	return structFields{fields, nameIndex}
}

func cachedTypeFields(t reflect.Type) structFields {
	if f, ok := fieldCache.Load(t); ok {
		return f.(structFields)
	}
	f, _ := fieldCache.LoadOrStore(t, typeFields(t))
	return f.(structFields)
}

func (d *Decoder) unmarshalTar(tr *tar.Reader, val reflect.Value) error {
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			val.Set(reflect.New(val.Type().Elem()))
		}
		val = val.Elem()
	}
	if val.Kind() == reflect.Interface && val.NumMethod() == 0 {
	}
	switch val.Kind() {
	case reflect.Interface:
	case reflect.Struct:
		fields := cachedTypeFields(val.Type())
		for {
			hdr, err := tr.Next()
			if err != nil {
				if err == io.EOF {
					break
				}
				return invalidError{format: "tar", err: err}
			}
			switch path.Ext(hdr.Name) {
			case ".plist":
				var f *field
				if i, ok := fields.nameIndex[hdr.Name]; ok {
					f = &fields.list[i]
				}
				if f == nil {
					// panic("missing")
					continue
				}
				if f != nil {
					subv := val
					for _, i := range f.index {
						if subv.Kind() == reflect.Pointer {
							if subv.IsNil() {
								if !subv.CanSet() {
									subv = reflect.Value{}
									break
								}
								subv.Set(reflect.New(subv.Type().Elem()))
							}
							subv = subv.Elem()
						}
						subv = subv.Field(i)
					}
					switch subv.Kind() {
					case reflect.Struct:
						subv = subv.Addr()
					case reflect.Map:
						subv.Set(reflect.MakeMap(subv.Type()))
					}
					buf, err := io.ReadAll(tr)
					if err != nil {
						return err
					}
					rd := bytes.NewReader(buf)
					pd := plist.NewDecoder(rd)
					err = pd.Decode(subv.Interface())
					if err != nil {
						return err
					}
				}
			}
		}
	case reflect.Map:
		if val.IsNil() {
			val.Set(reflect.MakeMap(val.Type()))
		}
		for {
			hdr, err := tr.Next()
			if err != nil {
				if err == io.EOF {
					break
				}
				return invalidError{format: "tar", err: err}
			}
			switch path.Ext(hdr.Name) {
			case "plist":
			}
		}
	}
	return nil
}

func (d *Decoder) Decode(v interface{}) error {
	magic := make([]byte, len(zstdFrameMagic))
	d.reader.Read(magic)
	d.reader.Seek(0, 0)
	rd := d.reader.(io.Reader)
	if bytes.Equal(magic, zstdFrameMagic) {
		zrd, err := zstd.NewReader(rd)
		if err != nil {
			return err
		}
		defer zrd.Close()
		rd = zrd
	}
	tr := tar.NewReader(rd)
	err := d.unmarshalTar(tr, reflect.ValueOf(v))
	return err
}

func NewDecoder(r io.ReadSeeker) *Decoder {
	return &Decoder{reader: r}
}
