package querier

import (
	"context"

	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/db"
)

type Querier interface {
	Fetch(ctx context.Context, contentType string, condition db.Condition) (contenttype.ContentTyper, error)

	List(ctx context.Context, contentType string, condition db.Condition) ([]contenttype.ContentTyper, int, error)
}
