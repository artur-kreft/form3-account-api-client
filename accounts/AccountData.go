package accounts

import (
	"github.com/satori/go.uuid"
)

type AccountData struct {
	Attributes     *AccountAttributes `json:"attributes,omitempty"`
	ID             uuid.UUID          `json:"id,omitempty"`
	OrganisationID uuid.UUID          `json:"organisation_id,omitempty"`
	Type           string             `json:"type,omitempty"`
	Version        *uint64            `json:"version,omitempty"`
}
