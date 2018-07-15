package asset

import "github.com/gobuffalo/packr"

// Assets is global variable for packr.NewBox, which return Box with files from specified folder
var Assets = packr.NewBox(".")
