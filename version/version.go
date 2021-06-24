package version

import "strconv"

// AppVersion build config
var AppVersion string

// GitCommit build config
var GitCommit string

// ProductID
var ProductID string

// MaxUses
var SMaxUses string
var MaxUses int

// MaxMarkets
var SMaxMarkets string
var MaxMarkets int

var Variant string

func init() {
	MaxUses, _ = strconv.Atoi(SMaxUses)
	MaxMarkets, _ = strconv.Atoi(SMaxMarkets)
}
