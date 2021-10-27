/*
   A template for adding persistent & accessible configuration systems to CLIs
                                                                                           - kendfss
*/
package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	// the following packages are only used for this particular project
	"github.com/nishanths/go-xkcd/v2"
	"time"
)

const (
	PROJECTNAME = "xkcd" // For creating a user-accessible data directory
)

var (
	//go:embed settings.json
	defaultSettings []byte
)

type (
	Settings struct {
		Keep      bool       `json:"keep"`
		Latest    xkcd.Comic `json:"latest"`
		LastSaw   xkcd.Comic `json:"lastSaw"`
		LastCheck time.Time  `json:"lastCheck"`
		Likes     []int      `json:"likes"`
	}
)

func (these Settings) Loadable() bool {
	if _, err := os.Stat(these.Path()); err != nil {
		return false
	}
	return true
}
func (these Settings) Exist() bool {
	_, err := os.Stat(these.Path())
	return !os.IsNotExist(err)
}
func (these *Settings) Load() error {
	if these.Exist() {
		data, err := ioutil.ReadFile(these.Path())
		if err != nil {
			return err
		}
		return json.Unmarshal(data, these)
	}
	if len(defaultSettings) > 0 {
		return json.Unmarshal(defaultSettings, these)
	}
	// default settings need to be created if they haven't already been
	data, err := json.MarshalIndent(these, "", "  ")
	if err != nil {
		return err
	}
	defaultSettings = data
	if err = ioutil.WriteFile("settings.json", data, fs.ModePerm); err != nil {
		return err
	}
	return these.Load() // try again
}

func (these *Settings) Restore() error {
	err := os.Remove(these.Path())
	if err != nil {
		return err
	}
	return json.Unmarshal(defaultSettings, these)
}

func (these Settings) Assure() error {
	dir := filepath.Dir(these.Path())
	return os.MkdirAll(dir, fs.ModeDir|fs.ModePerm)
}

func (these *Settings) Save() error {
	data, err := json.MarshalIndent(these, "", "  ")
	if err == nil {
		these.Assure()
		return ioutil.WriteFile(these.Path(), data, fs.ModePerm)
	}
	return fmt.Errorf("Couldn't save settings file: %w", err)
}

// This is used to generate an initial copy of the configuration file
// Do not use unless you are comfortable with a "settings.json" file
// in your current directory
func (these *Settings) Overwrite() error {
	data, err := json.MarshalIndent(these, "", "  ")
	if err == nil {
		internalError.Abortf("Couldn't save settings file: %w", err)
	}
	return ioutil.WriteFile("settings.json", data, fs.ModePerm)
}

func (these Settings) Open(dir bool) error {
	var (
		cmd  *exec.Cmd
		path string = these.Path()
	)
	if dir {
		path = filepath.Dir(path)
	}
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

func (these Settings) Path() string {
	user, err := os.UserHomeDir()
	if err != nil {
		internalError.Abort(err.Error())
	}
	return filepath.Join(user, ".local", "share", PROJECTNAME, "settings.json")
}
