package dmgooglemap

import (
	models "dm/models"
)

//DMGoogleMap datatype
type DMGoogleMap struct {
	*models.Datatype
}

func (googlemap DMGoogleMap) Validate() {
	googlemap.Datatype.Store()
}
