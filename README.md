# autotag

Automatically tag media files based on filenames.

## Installation

Requires Golang 1.7+:

```
go get -u github.com/ghetzel/autotag
```


## Overview

`autotag` takes one or more directories, scans them for files, and uses regular expressions to parse the
filenames and apply media tags to them.  Supported tag formats are ID3v1/v2, Ogg Vorbis comments (`.ogg`,
`.flac`, `.opus`), and anything else that [TagLib](http://taglib.org) supports.


## Pattern Files

The `autotag` scanner starts in the directory specified on the command line, looking for an `autotag.list` file.
This file is a simple text format that lists, one per line, regular expressions that `autotag` will use to attempt
to parse filenames into meaningful tags.

### Format

An example file looks like this:

```
# parse album/artist-title filenames
/(?P<album>.*?)/(?P<artist>.*?) - (?<title>.*)\..{2,4}$

# parse album/track title filenames
/(?P<album>.*?)/(?P<track>\d+?) (?<title>.*)\..{2,4}$
```

### Search Paths

If the directory specified does not contain this file, that directory's parent is searched, and so
on until reaching the root (/) directory.  It is an error to not locate the `autotag.list` file in
any directories.

### Aggregation

If multiple directories contain an `autotag.list` file, they are all merged together in the order they
are found.  The first match across all pattern files will be used, so complex naming schemes can be
handled by putting ever-more-specific pattern files deeper in the file hierarchy.  For example, you
may have a general rule that applies to all files in a media library, but have specific rules that apply
to a subset of files, or even a single directory.  This search and aggregation behavior supports such
a scheme.