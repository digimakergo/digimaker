//Author xc, Created on 2019-04-03 21:00
//{COPYRIGHTS}

package util

//DefaultSettingsType and DefaultSettings set all the default,
// including which config folder it is.
type DefaultSettingsType struct {
	ConfigFile   string
	ConfigFolder string
}

var DefaultSettings = DefaultSettingsType{"site", "/Users/xc/go/caf-prototype/src/dm/configs"}
