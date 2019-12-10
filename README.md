# flagvar

> A tiny library with collection of [`flag.Value`][flagValue] implementations

# Overview

This is a collection of different useful [flag.Value][flagValue]
implementations. Works good with [flagutil][flagutil] library.

# Usage

Here is a simple program that reads a date from the flag and prints it as a
unix timestamp:

```go
package main

import (
	"flag"
	"time"

	"github.com/gobwas/flagvar"
)

func main() {
	flags := flag.NewFlagSet("time", flag.ExitOnError)	
	
	var t time.Time
	flags.Var(
		&flagvar.TimeValue{
			P:      &t,
			Layout: "02.01.2006",
		},
		"date", "01.01.1970",
		"time to print as a unix timestamp in form `dd.mm.yyyy`",
	)
	
	flags.Parse()

	fmt.Fprintln(os.Stdout, t.Unix())	
}
```

[flagValue]: https://golang.org/pkg/flag#FlagSet
[flagutil]: https://github.com/gobwas/flagutil
