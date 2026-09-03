package assets

import _ "embed"

//go:embed containers.json
var Containers []byte

//go:embed items.json
var Items []byte

//go:embed profiles.json
var Profiles []byte
