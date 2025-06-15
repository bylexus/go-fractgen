package lib

import "io/fs"

type AppContext struct {
	EmbeddedPresets []byte
	WebrootFS       fs.FS
}
