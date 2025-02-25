# command

This package contains some highly experimental abstractions around creating command-line incantations. While there are ample libraries for parsing CLI flags, there are few for generating them in a programmatic way. This package aims to do that as well as provide a higher-leval abstraction for constructing such incantations. The general idea is to code the basics for working with each of the individual commands there are specialized structs for.

I completely acknowledge that in certain circumstances, it is much better to either use an official API client or to use the Golang standard library for certain operations. This package is for situations where one simply cannot do that or doing so is much more involved than it arguably should be.
