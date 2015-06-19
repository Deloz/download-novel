package sites

import (
	"reflect"
	"strings"
)

type Site interface {
	ParseNovelList()
}

// todo...
func HasSite(s Site, siteModel string) bool {
	siteModel = strings.Title(siteModel)
	sv := reflect.ValueOf(s).Elem()
	t := sv.Type()
	siteType := reflect.TypeOf((*Site)(nil)).Elem()

	for i := 0; i < sv.NumField(); i++ {
		f := t.Field(i)
		if f.Name == siteModel && f.Type.Implements(siteType) {
			return true
		}
	}

	return false
}
