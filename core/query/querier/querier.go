package querier

import (
	"context"

	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/db"
)

type Querier interface {
	Fetch(ctx context.Context, contentType string, condition db.Condition) (contenttype.ContentTyper, error)

	ListWithUser(ctx context.Context, userID int, contentType string, condition db.Condition) ([]contenttype.ContentTyper, int, error)

	List(ctx context.Context, contentType string, condition db.Condition) ([]contenttype.ContentTyper, int, error)
}
