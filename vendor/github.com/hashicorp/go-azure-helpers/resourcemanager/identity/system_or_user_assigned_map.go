package identity

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

var _ json.Marshaler = &SystemOrUserAssignedMap{}

type SystemOrUserAssignedMap struct {
	Type        Type                                   `json:"type"`
	PrincipalId string                                 `json:"principalId"`
	TenantId    string                                 `json:"tenantId"`
	IdentityIds map[string]UserAssignedIdentityDetails `json:"userAssignedIdentities"`
}

func (s *SystemOrUserAssignedMap) MarshalJSON() ([]byte, error) {
	// we use a custom marshal function here since we can only send the Type / UserAssignedIdentities field
	identityType := TypeNone
	userAssignedIdentityIds := map[string]UserAssignedIdentityDetails{}

	if s != nil {
		if s.Type == TypeSystemAssigned {
			identityType = TypeSystemAssigned
		}
		if s.Type == TypeUserAssigned {
			identityType = TypeUserAssigned
		}

		if identityType != TypeNone {
			userAssignedIdentityIds = s.IdentityIds
		}
	}

	out := map[string]interface{}{
		"type":                   string(identityType),
		"userAssignedIdentities": userAssignedIdentityIds,
	}
	return json.Marshal(out)
}

// ExpandSystemOrUserAssignedMap expands the schema input into a SystemOrUserAssignedMap struct
func ExpandSystemOrUserAssignedMap(input []interface{}) (*SystemOrUserAssignedMap, error) {
	identityType := TypeNone
	identityIds := make(map[string]UserAssignedIdentityDetails, 0)

	if len(input) > 0 {
		raw := input[0].(map[string]interface{})
		typeRaw := raw["type"].(string)
		if typeRaw == string(TypeSystemAssigned) {
			identityType = TypeSystemAssigned
		}
		if typeRaw == string(TypeUserAssigned) {
			identityType = TypeUserAssigned
		}

		identityIdsRaw := raw["identity_ids"].([]interface{})
		for _, v := range identityIdsRaw {
			identityIds[v.(string)] = UserAssignedIdentityDetails{
				// intentionally empty since the expand shouldn't send these values
			}
		}
	}

	if len(identityIds) > 0 && identityType != TypeUserAssigned {
		return nil, fmt.Errorf("`identity_ids` can only be specified when `type` is set to %q", string(TypeUserAssigned))
	}

	return &SystemOrUserAssignedMap{
		Type:        identityType,
		IdentityIds: identityIds,
	}, nil
}

// FlattenSystemOrUserAssignedMap turns a SystemOrUserAssignedMap into a []interface{}
func FlattenSystemOrUserAssignedMap(input *SystemOrUserAssignedMap) (*[]interface{}, error) {
	if input == nil || (input.Type != TypeSystemAssigned && input.Type != TypeUserAssigned) {
		return &[]interface{}{}, nil
	}

	identityIds := make([]string, 0)
	for raw := range input.IdentityIds {
		id, err := commonids.ParseUserAssignedIdentityIDInsensitively(raw)
		if err != nil {
			return nil, fmt.Errorf("parsing %q as a User Assigned Identity ID: %+v", raw, err)
		}
		identityIds = append(identityIds, id.ID())
	}

	return &[]interface{}{
		map[string]interface{}{
			"type":         string(input.Type),
			"identity_ids": identityIds,
			"principal_id": input.PrincipalId,
			"tenant_id":    input.TenantId,
		},
	}, nil
}
