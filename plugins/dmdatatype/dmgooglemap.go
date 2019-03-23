package dmgooglemap

import (
	"dmcaf/models/base"
)

//DMGoogleMap datatype
type DMGoogleMap struct {
	*base.Datatype
}

func (googlemap DMGoogleMap) Validate() {
	googlemap.Datatype.Save()
}
