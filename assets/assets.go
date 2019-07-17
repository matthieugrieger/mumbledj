package assets

import "github.com/gobuffalo/packr/v2"

// Assets is global variable for packr.NewBox, which return Box with files from specified folder
var Assets = packr.New("assets", "assets")
