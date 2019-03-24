package field

import "dm/models"

//TextField is a field for normal text line. It implements Datatyper
type TextField struct {
	*models.Field
}
