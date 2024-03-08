package spotify

import (
	"context"
	"encoding/json"
	"net/http"
)

func (c *Client) AvailableDevices(ctx context.Context) (AvailableDevices, error) {
	var availableDevices AvailableDevices
	err := c.GetJSON(ctx, http.MethodGet, "/me/player/devices", &availableDevices)
	return availableDevices, err
}

func UnmarshalAvailableDevices(data []byte) (AvailableDevices, error) {
	var r AvailableDevices
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *AvailableDevices) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type AvailableDevices struct {
	Devices []Device `json:"devices"`
}

type Device struct {
	ID               string `json:"id"`
	IsActive         bool   `json:"is_active"`
	IsPrivateSession bool   `json:"is_private_session"`
	IsRestricted     bool   `json:"is_restricted"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	VolumePercent    int64  `json:"volume_percent"`
	SupportsVolume   bool   `json:"supports_volume"`
}
