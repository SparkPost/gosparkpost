// Package transmissions interacts with the SparkPost Transmissions API.
// https://www.sparkpost.com/api#/reference/transmissions
package transmissions

import (
	"time"

	"github.com/SparkPost/go-sparkpost/api"
	//"github.com/SparkPost/go-sparkpost/api/recipients"
	"github.com/SparkPost/go-sparkpost/api/templates"
)

// Transmissions is your handle for the Transmissions API.
type Transmissions struct{ api.API }

func New(cfg *api.Config) (*Transmissions, error) {
	// FIXME: allow caller to set api version
	t := &Transmissions{}
	err := t.Init(cfg, "/api/v1/transmissions")
	if err != nil {
		return nil, err
	}
	return t, nil
}

// Transmission is the JSON structure accepted by and returned from the SparkPost Transmissions API.
type Transmission struct {
	ID         string      `json:"id,omitempty"`
	State      string      `json:"state,omitempty"`
	Options    Options     `json:"options,omitempty"`
	Recipients interface{} `json:"recipients"`

	Content templates.Content `json:"content"`
}

// FIXME: Open and ClickTracking default to true if not specified
// FIXME: assert Options.Sandbox is only recognized for SparkPost.com
// FIXME: assert Options.SkipSuppression is only recognized for SparkPost Elite

// Options specifies settings to apply to this Transmission.
// If not specified, and present in templates.Options, those values will be used.
type Options struct {
	StartTime       time.Time `json:"start_time,omitempty"`
	OpenTracking    bool      `json:"open_tracking,omitempty"`
	ClickTracking   bool      `json:"click_tracking,omitempty"`
	Transactional   bool      `json:"transactional,omitempty"`
	Sandbox         bool      `json:"sandbox,omitempty"`
	SkipSuppression bool      `json:"skip_suppression,omitempty"`
}
