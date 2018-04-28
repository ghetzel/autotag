package main

import (
	"fmt"

	"github.com/ghetzel/go-stockutil/maputil"
	taglib "github.com/wtolson/go-taglib"
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
	if file, err := taglib.Read(self.Path); err == nil {
		tags := maputil.M(self.Tags)
		didAnything := false

		if v := tags.Int(`track`, 0); v > 0 {
			file.SetTrack(int(v))
			didAnything = true
		}

		if v := tags.Int(`year`, 0); v > 0 {
			file.SetYear(int(v))
			didAnything = true
		}

		if v := tags.String(`album`); v != `` {
			file.SetAlbum(v)
			didAnything = true
		}

		if v := tags.String(`artist`); v != `` {
			file.SetArtist(v)
			didAnything = true
		}

		if v := tags.String(`title`); v != `` {
			file.SetTitle(v)
			didAnything = true
		}

		if didAnything {
			return file.Save()
		} else {
			return nil
		}
	} else {
		return err
	}
}
