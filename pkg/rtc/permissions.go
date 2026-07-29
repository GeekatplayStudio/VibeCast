// Package rtc manages room state, participant connections, media track publication, and permissions.
package rtc

// Permissions defines fine-grained access control rights for room participants.
type Permissions struct {
	CanPublish     bool `json:"can_publish"`
	CanSubscribe   bool `json:"can_subscribe"`
	CanPublishData bool `json:"can_publish_data"`
	IsAdmin        bool `json:"is_admin"`
	IsHidden       bool `json:"is_hidden"`
}

// DefaultPermissions provides baseline rights for connected room participants.
func DefaultPermissions() Permissions {
	return Permissions{
		CanPublish:     true,
		CanSubscribe:   true,
		CanPublishData: true,
		IsAdmin:        false,
		IsHidden:       false,
	}
}

// AdminPermissions grants full room administration rights.
func AdminPermissions() Permissions {
	return Permissions{
		CanPublish:     true,
		CanSubscribe:   true,
		CanPublishData: true,
		IsAdmin:        true,
		IsHidden:       false,
	}
}
