package main

type TagWriter interface {
	WriteFile(path string, tags map[string]interface{}) error
}
