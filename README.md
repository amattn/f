f
===

fast access to directories and scripts.

Currently tested on unix-like systems (linux, darwin)

### Install

The two primary installation options are with `go install` or by downloading from the [releases page][].  For shortcuts that you want to be `eval`'d by the parent shell (such as `cd $DIR`), you also need to add a small wrapper to your `.bashrc` or the config file of your given shell.


**1a:** Via `go install`: 

This will typically install into your `$GOBIN` directory.

    go install github.com/amattn/f

**1b:** Via [releases page][]: 

**2:** Shell wrapper

For bash, add the following lines to your bash profile:

    ff() {
      if [ ! -n "$1" ]; then
        # just print the menu
        f
      else
        # pass args to f
        eval `f --print "$@"`
      fi  
    } 


[releases page]: https://github.com/amattn/f/releases

#### TODO

- pass extra args along
- bash tab completion
