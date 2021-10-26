package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/nishanths/go-xkcd/v2"
)

var (
	comic, latest             xkcd.Comic
	max, min, id              int   // to help with choosing a random comic // first comic to look for
	err                       error // to avoid unwanted duplication of variables in lower scopes
	dir, force, view, clean   bool
	keep, list, restore, prev bool

	bgctx  context.Context = context.Background()
	client *xkcd.Client    = xkcd.NewClient()

	randomMode = flag.NewFlagSet("random", flag.ExitOnError) // get a random comic
	fetchMode  = flag.NewFlagSet("fetch", flag.ExitOnError)  // get a particular comic - defaults to random
	configMode = flag.NewFlagSet("config", flag.ExitOnError) // open settings file in json editor
	helpMode   = flag.NewFlagSet("help", flag.ExitOnError)
	likeMode   = flag.NewFlagSet("likes", flag.ExitOnError)
	modes      = map[string]*flag.FlagSet{
		randomMode.Name(): randomMode,
		fetchMode.Name():  fetchMode,
		configMode.Name(): configMode,
		helpMode.Name():   helpMode,
		likeMode.Name():   likeMode,
	}

	settings    = &Settings{}
	comicFolder = filepath.Join(filepath.Dir(settings.Path()), "comics")
)

func setLatestComic(force bool) error {
	if force || time.Since(settings.LastCheck) > 24*time.Hour {
		if latest, err = client.Latest(bgctx); err != nil {
			return err
		}
		settings.LastCheck = time.Now()
		settings.Latest = latest
		if err = settings.Save(); err != nil {
			return err
		}
	} else {
		latest = settings.Latest
	}
	max = latest.Number
	id = latest.Number
	return nil
}

func help(modeName string) {
	mode, exists := modes[modeName]
	if exists {
		fmt.Println(modeName)
		mode.PrintDefaults()
		if modeName == helpMode.Name() {
			fmt.Println("  Print this information")
		} else if modeName == configMode.Name() {
			fmt.Printf("config located @ %q\n", settings.Path())
		}
	}
}

func restoreLikes() {
	ctx, cancel := context.WithCancel(bgctx)
	defer cancel()
	for _, index := range settings.Likes {
		comic, err := client.Get(ctx, index)
		// comic, err := client.Get(context.Background(), index)
		if err != nil {
			internalError.Abortf("Restoring Likes: Comic number %d doesn't exist or couldn't be accessed. Double check your connection and try again:\n\t%s", id, err)
		}
		_, err = DownloadFile(comic.ImageURL)
		if err != nil {
			internalError.Abort(err.Error())
		}
	}
}

// func getComic(id int) error {}

func displayComic(id int) {
	comic, err := client.Get(bgctx, id)
	if err != nil {
		userInputError.Abortf("Cannot display comic number %d as it doesn't exist or couldn't be accessed. Double check your connection and try again:\n\t%s", id, err)
	}

	fmt.Printf("Comic #%d: %s\n", comic.Number, comic.Title)
	path, err := DownloadFile(comic.ImageURL)
	if err != nil {
		internalError.Abort(err.Error())
	}
	Open(path)

	settings.LastSaw = comic
	if err = settings.Save(); err != nil {
		internalError.Abortf("Couldn't save settings after displaying comic #%d: %s", id, err)
	}
}

func comicName(id int) string {
	comic, err := client.Get(context.Background(), id)
	if err != nil {
		internalError.Abortf("Couldn't get the name of comic #%d: %s\n\tDouble check your internet connection\n", id, err)
	}
	return filepath.Base(comic.ImageURL)
}

func cleanUp() error {
	faveNames := []string{}
	for _, index := range settings.Likes {
		faveNames = append(faveNames, comicName(index))
	}
	for _, name := range Files(comicFolder) {
		if !containsStr(faveNames, filepath.Base(name)) {
			err := os.Remove(name)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func userLikes() map[int]string {
	table := make(map[int]string)
	for _, index := range settings.Likes {
		table[index] = comicName(index)
	}
	return table
}

func Usage() {
	fmt.Println("Usage of", os.Args[0], "[mode]", "[params]")
	for name := range modes {
		help(name)
		fmt.Println()
	}

	fmt.Printf("all comics downloaded to %q\n", comicFolder)
}

func main() {
	min = 1
	err = settings.Load()
	if err != nil && err != io.EOF {
		internalError.Abortf("Couldn't load settings: %s", err)
	}
	keep = settings.Keep

	if err = setLatestComic(false); err != nil {
		internalError.Abortf("Couldn't get the index for the latest comic. %s\n\tDouble check your internet connection", err)
	}

	randomMode.IntVar(&max, "max", max, "the maximum index for a random comic (double-check your internet connection if the default is zero)")
	randomMode.IntVar(&min, "min", min, "the minimum index for a random comic")

	fetchMode.IntVar(&id, "id", id, "the index of the comic you want to see")
	fetchMode.BoolVar(&force, "f", force, "grab the latest comic even if checked within the last 24 hours")

	configMode.BoolVar(&keep, "k", keep, "Toggle whether or not to keep comics after displaying them (default: false)")
	configMode.BoolVar(&dir, "d", dir, "Open the directory of the config file, instead of the file itself")

	likeMode.BoolVar(&view, "v", view, "View your likes")
	likeMode.BoolVar(&clean, "c", clean, "Clean up. (removes all saved comics that aren't liked)")
	likeMode.BoolVar(&list, "l", list, "list the names of the comics you've liked")
	likeMode.BoolVar(&restore, "r", restore, "Re-download any missing likes")
	likeMode.BoolVar(&prev, "p", prev, "Add previous comic to likes")

	flag.Usage = Usage

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case likeMode.Name():
			err = likeMode.Parse(os.Args[2:])
			if err != nil {
				userInputError.Abortf("Couldn't parse args %s", err)
			}
			if prev {
				settings.Likes = append(settings.Likes, settings.LastSaw.Number)
				err = settings.Save()
				if err != nil {
					internalError.Abortf("Couldn't save settings: %s", err)
				}
				return
			} else if view {
				for _, index := range settings.Likes {
					displayComic(index)
				}
				return
			} else if list {
				for _, index := range settings.Likes {
					fmt.Printf("%d:\t%s\n", index, comicName(index))
				}
				return
			} else if clean {
				cleanUp()
				return
			} else if restore {
				restoreLikes()
				return
			}
			for _, arg := range os.Args[2:] {
				v, err := strconv.ParseInt(arg, 0, 0)
				if err != nil {
					userInputError.Abortf("Couldn't parse %q: %s", arg, err)
				}
				if !containsInt(settings.Likes, int(v)) {
					settings.Likes = append(settings.Likes, int(v))
				}
			}

			err = settings.Save()
			if err != nil {
				internalError.Abortf("Couldn't save settings: %s", err)
			}
			return
		case helpMode.Name():
			err = helpMode.Parse(os.Args[2:])
			if err != nil {
				userInputError.Abortf("Couldn't parse args %s", err)
			}
			if len(os.Args) == 2 {
				flag.Usage()
			} else {
				for _, name := range os.Args[2:] {
					help(name)
				}
			}
			return
		case configMode.Name():
			err = configMode.Parse(os.Args[2:])
			if err != nil {
				userInputError.Abortf("Couldn't parse args %s", err)
			}
			if keep {
				settings.Keep = !settings.Keep
				return
			}
			err = settings.Open(dir)
			if err != nil {
				internalError.Abortf("Couldn't open settings: %s", err)
			}
			return
		case randomMode.Name():
			err = randomMode.Parse(os.Args[2:])
			if err != nil {
				userInputError.Abortf("Couldn't parse args %s", err)
			}
			rand.Seed(time.Now().Unix())
			if max > min {
				id = min + rand.Intn(max-min)
			} else {
				id = max + rand.Intn(min-max)
			}
		case fetchMode.Name():
			err = fetchMode.Parse(os.Args[2:])
			if err != nil {
				userInputError.Abortf("Couldn't parse args %s", err)
			}
		case "next":
			id = settings.LastSaw.Number
			if id <= settings.Latest.Number-1 {
				id++
			} else {
				id = 1
			}
		case "prev":
			id = settings.LastSaw.Number
			if id >= 2 {
				id--
			} else {
				id = latest.Number
			}
		default:
			flag.Usage()
			flag.PrintDefaults()
			userInputError.Abortf("Couldn't parse args: %s", os.Args[1:])
		}
	} else {
		flag.Usage()
		flag.PrintDefaults()
		noError.Abort("")
	}

	displayComic(id)
}
