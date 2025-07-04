package main

import (
	"os"
	"path/filepath"
)

func recursiveUpkSearch(root string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		var fileExt = filepath.Ext(info.Name())
		if fileExt == ".upk" || fileExt == ".udk" {
			matches = append(matches, path)
		}

		return nil
	})
	return matches, err
}

func (a *App) GetMaps(paths []string) map[string]string {
	maps := map[string]string{}

	for _, path := range paths {
		new, err := recursiveUpkSearch(path)
		if err != nil {
			println("WARN: failed to search path: " + path)
			continue
		}

		for _, map_path := range new {
			maps[filepath.Base(map_path)] = map_path
		}
	}

	return maps
}
