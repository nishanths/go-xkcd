package main

import (
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func Listdir(pth string) []string {
	files, err := os.ReadDir(pth)
	if err != nil {
		internalError.Abort(err.Error())
	}
	rack := make([]string, 0)
	for _, file := range files {
		rack = append(rack, file.Name())
	}
	return rack
}
func Files(root string) []string {
	paths := make([]string, 0)
	for _, name := range Listdir(root) {
		pth := filepath.Join(root, name)
		stat, err := os.Lstat(pth)
		if err != nil {
			internalError.Abort(err.Error())
		}
		switch mode := stat.Mode(); {
		case mode.IsRegular():
			paths = append(paths, pth)
		case mode.IsDir():
			go Merge(&paths, Files(pth))
		}
	}
	return paths
}

func Open(path string) error {
	var (
		cmd *exec.Cmd
	)
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", path)
	case "linux":
		cmd = exec.Command("xdg-open", path)
	case "windows":
		cmd = exec.Command("start", path)
	default:
		internalError.Abort("Do not know how to open files on this Operating System")
	}
	return cmd.Run()
}

func DownloadFile(url string) (string, error) {
	_, err := os.Stat(comicFolder)
	if os.IsNotExist(err) {
		err := os.MkdirAll(comicFolder, fs.ModeDir|fs.ModePerm)
		if err != nil {
			return "", err
		}
	}
	path := filepath.Join(comicFolder, filepath.Base(url))
	// avoid downloading the same comic twice
	_, err = os.Stat(path)
	if !os.IsNotExist(err) {
		return path, nil
	}
	out, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	return path, nil
}
