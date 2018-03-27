package frontend

import "github.com/elazarl/go-bindata-assetfs"

//go:generate go-bindata-assetfs -pkg frontend -ignore .*\.go -prefix build ./build

func AssetFS() *assetfs.AssetFS {
	afs := assetFS()
	afs.Prefix = "/"
	return afs
}
