package image

import (
	"context"

	"github.com/digimakeras/digimaker/core/util"
	"github.com/digimakeras/digimaker/dcp/config"
)

//Get full path from image's interal path
//size: original/default/600
func ImagePath(ctx context.Context, path string, size string) string {
	imageUrl := config.DmViper(ctx).GetString("general.image_url")
	sizeStr := ""
	sizeWithSlash := ""
	if size != "" && size != "original" {
		sizeStr = size
		sizeWithSlash = sizeStr + "/"
	}
	variables := map[string]string{"path": path,
		"size":            sizeStr,
		"size_with_slash": sizeWithSlash,
	}
	result := util.ReplaceStrVar(imageUrl, variables)
	return result
}
