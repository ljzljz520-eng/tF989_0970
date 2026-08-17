package web

import "embed"

//go:embed assets/* views/*
var Files embed.FS
