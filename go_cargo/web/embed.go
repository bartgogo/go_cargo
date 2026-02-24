// Package web 嵌入前端静态资源
package web

import "embed"

//go:embed index.html app.html
var StaticFS embed.FS
