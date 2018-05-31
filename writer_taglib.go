package main

import (
	"github.com/ghetzel/go-stockutil/maputil"
	taglib "github.com/wtolson/go-taglib"
)

type TaglibWriter struct {
}

func (self *TaglibWriter) WriteFile(path string, tags map[string]interface{}) error {
	if file, err := taglib.Read(path); err == nil {
		tags := maputil.M(tags)
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
