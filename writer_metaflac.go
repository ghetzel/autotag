package main

import (
	"fmt"
	"os/exec"

	"github.com/ghetzel/go-stockutil/maputil"
)

type MetaflacWriter struct {
}

func (self *MetaflacWriter) WriteFile(path string, tagsmap map[string]interface{}) error {
	tagArgs := make([]string, 0)
	names := make([]string, 0)
	tags := maputil.M(tagsmap)

	if v := tags.Int(`disc`, 0); v > 0 {
		names = append(names, `DISCNUMBER`)
		tagArgs = append(tagArgs, fmt.Sprintf("DISCNUMBER=%d", v))
	}

	if v := tags.Int(`track`, 0); v > 0 {
		names = append(names, `TRACKNUMBER`)
		tagArgs = append(tagArgs, fmt.Sprintf("TRACKNUMBER=%d", v))
	}

	if v := tags.Int(`year`, 0); v > 0 {
		names = append(names, `DATE`)
		tagArgs = append(tagArgs, fmt.Sprintf("DATE=%d", v))
	}

	if v := tags.String(`album`); v != `` {
		names = append(names, `ALBUM`)
		tagArgs = append(tagArgs, fmt.Sprintf("ALBUM=%v", v))
	}

	if v := tags.String(`artist`); v != `` {
		names = append(names, `ARTIST`)
		tagArgs = append(tagArgs, fmt.Sprintf("ARTIST=%v", v))
	}

	if v := tags.String(`title`); v != `` {
		names = append(names, `TITLE`)
		tagArgs = append(tagArgs, fmt.Sprintf("TITLE=%v", v))
	}

	if len(tagArgs) > 0 {
		args := []string{}

		for _, nm := range names {
			args = append(args, fmt.Sprintf("--remove-tag=%v", nm))
		}

		for _, ta := range tagArgs {
			args = append(args, fmt.Sprintf("--set-tag=%v", ta))
		}

		args = append(args, path)

		return exec.Command(`metaflac`, args...).Run()
	} else {
		return nil
	}
}
