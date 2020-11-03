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
	opts *ScanOptions
}

func (self *FileMatch) String() string {
	return fmt.Sprintf("%v:\n  %v", self.Path, maputil.Join(self.Tags, `=`, "\n  "))
}

func (self *FileMatch) Apply() error {
	var writer TagWriter

	switch path.Ext(self.Path) {
	case `.flac`, `.ogg`:
		writer = new(MetaflacWriter)
	default:
		writer = new(TaglibWriter)
	}

	if err := writer.WriteFile(self.Path, self.Tags); err == nil {
		if self.opts != nil {
			if err := self.opts.writeXattr(self.Path, self.Tags); err != nil {
				return err
			}
		}

		return nil
	} else {
		return err
	}
}
