package risk

import (
	"errors"
	"log/slog"

	"github.com/google/uuid"
)

// Going to leave this as a map for now, even if the encoding is a bit ugly
// I could either write a custom encoder to ditch the key values of the map, but
// if it was required I'd probably just use two slices. One with indeces and one
// with values, and just return values here
// Either way, I'm assuming for this, this is fine
var risks = make(map[string]Risk)

type Risk struct {
	Id          string `json:"id" validate:"isdefault"`
	State       string `json:"state" validate:"required,oneof=open closed accepted investigating"`
	Title       string `json:"title" validate:"required,min=5,max=999"`
	Description string `json:"description" validate:"max=9999`
}

// At some point could use this to support something other than v4
func (r *Risk) generateUUID() {
	r.Id = uuid.New().String()
}

func createRisk(risk Risk) error {
	// Leaving this here as I think for this case we can assume we wouldn't want to generate this UUID anywhere else
	// and I don't think I want to test any of google and validators uuid generation/validation
	risk.generateUUID()

	// I'm more likely to be struck by lightning in my basement before this gets hit, but
	// if there where other reasons to fail they'd go here
	if _, ok := risks[risk.Id]; ok {
		return errors.New("Cannot add risk, id already exists")
	}

	risks[risk.Id] = risk

	return nil
}

func getRisk(id string) *Risk {
	v, ok := risks[id]

	if !ok {
		slog.Error("Cannot get risk, doesn't exist", "id", id)
		return nil
	}

	return &v
}
