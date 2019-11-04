package main

// We want to replace this variable at build time with "-ldflags -X main.GitSHA=xxx", where const is not supported.
var (
	Version = "0.0.2"
	GitSHA  = ""
	Built   = ""
)
