# Name
**xkcd** - download XKCD comics

# Synopsis
**xkcd** [mode] [params]

# Description
**xkcd** download and view xkcd comics by index or at random. 

# Options
**config**
  -d  Open the directory of the config file, instead of the file itself
  -k  Toggle whether or not to keep comics after displaying them (default: false)
config located @ `~/.local/share/xkcd/settings.json`

**fetch**
  -f  grab the latest comic even if checked within the last 24 hours
  -id int
      the index of the comic you want to see (default 2532)

**likes**
  -c  Clean up. (removes all saved comics that aren't liked)
  -l  list the names of the comics you've liked
  -p  Add previous comic to likes
  -r  Re-download any missing likes
  -v  View your likes

**random**
  -max int
      the maximum index for a random comic (double-check your internet connection if the default is zero) (default 2532)
  -min int
      the minimum index for a random comic (default 1)

**help**
  Print this information

all comics downloaded to `~/.local/share/xkcd/comics`


# Examples

**xkcd** help [random]  
    # info on a given mode
**xkcd** next  
    # fetch the comic whose index is one before the last one you saw
**xkcd** prev  
    # fetch the comic whose index is one greater than the last comic you saw
**xkcd** random -max 22  
    # fetch a random comic whose index is smaller than 22
**xkcd** likes 12  
    # adds comic number 12 to your likes
**xkcd** fetch -f 
    # fetch the latest comic even if you checked for it in the past 24 hours



# Exit Values
**0** No Error  
**1** Internal Error  
**2** User Input Error  
**3** Warning

# License
MIT  
