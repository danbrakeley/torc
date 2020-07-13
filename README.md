# torc

TORrent Cleaner

This utility compares the contents of a torrent file with your local disk, and lists any differences found.

There is an option to delete any files found on disk that are not also part of the torrent (see `torc help cmp`).

```text
$ torc help
NAME:
   torc - compare local files to a torrent file

USAGE:
   torc [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
     compare, cmp  compares a torrent file with actual files on disk
     list, ls      list all files for a given torrent
     help, h       Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version

COPYRIGHT:
   (c) 2017-2020 Dan Brakeley
```
