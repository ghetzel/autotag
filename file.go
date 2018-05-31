package main

import (
	"fmt"
	"path"

	"github.com/ghetzel/go-stockutil/maputil"
)

type FileMatch struct {
	Path string
	Rule *Rule
	Tags map[string]interface{}
}

func (self *FileMatch) String() string {
	return fmt.Sprintf("%v:\n  %v", self.Path, maputil.Join(self.Tags, `=`, "\n  "))
}

func (self *FileMatch) Apply() error {
	var writer TagWriter

	switch path.Ext(self.Path) {
	case `.flac`, `.ogg`:
		writer = &MetaflacWriter{}
	default:
		writer = &TaglibWriter{}
	}

	return writer.WriteFile(self.Path, self.Tags)
}
