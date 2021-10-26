xkcd-cli
---

A command-line interface to the example in the outer repository.  
**xkcd** allows the user to download and view xkcd comics by index (`fetch` mode) or at random, as well as keeping track of favourite comics (`like` mode).  
  
You may also find more robust documentation by running `man xkcd`, after installing.  

## Install
```shell
git clone https://github.com/nishanths/go-xkcd
cd go-xkcd/cli
make
```

## Uninstall
```shell
make remove
```
- This _will not_ delete downloaded comics or settings under `~/.local/share/xkcd`  

## Usage
```shell
xkcd [mode] [params]
```
modes:
- `config`  
- `fetch`  
- `likes`  
- `random`  
- `help`  

## Examples
```shell
xkcd help [mode]  # info on a given mode
xkcd next  # fetch the comic whose index is one before the last one you saw
xkcd prev  # fetch the comic whose index is one greater than the last comic you saw
xkcd fetch -f  # fetch the latest comic even if you checked for it in the past 24 hours
xkcd random  # fetch a random
xkcd random -max 22  # fetch a random comic whose index is smaller than 22
xkcd likes 12  # add comic number 12 to your likes
```

## Exit Values
- `0` No Error  
- `1` Internal Error  
- `2` User Input Error  
- `3` Warning  

