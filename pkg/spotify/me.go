package spotify

import (
	"context"
	"encoding/json"
	"net/http"
)

func (c *Client) Me(ctx context.Context) (Me, error) {
	var me Me
	err := c.GetJSON(ctx, http.MethodGet, "/me", &me)
	return me, err
}

func UnmarshalMe(data []byte) (Me, error) {
	var r Me
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Me) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Me struct {
	Country         string          `json:"country"`
	DisplayName     string          `json:"display_name"`
	Email           string          `json:"email"`
	ExplicitContent ExplicitContent `json:"explicit_content"`
	ExternalUrls    ExternalUrls    `json:"external_urls"`
	Followers       Followers       `json:"followers"`
	Href            string          `json:"href"`
	ID              string          `json:"id"`
	Images          []Image         `json:"images"`
	Product         string          `json:"product"`
	Type            string          `json:"type"`
	URI             string          `json:"uri"`
}

type ExplicitContent struct {
	FilterEnabled bool `json:"filter_enabled"`
	FilterLocked  bool `json:"filter_locked"`
}

type ExternalUrls struct {
	Spotify string `json:"spotify"`
}

type Followers struct {
	Href  string `json:"href"`
	Total int64  `json:"total"`
}

type Image struct {
	URL    string `json:"url"`
	Height int64  `json:"height"`
	Width  int64  `json:"width"`
}
