// Package recipients interacts with the SparkPost Recipient Lists API.
// https://www.sparkpost.com/api#/reference/recipient-lists
package recipient_lists

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/SparkPost/go-sparkpost/api"
)

// RecipientLists is your handle for the Recipient Lists API.
type RecipientLists struct{ api.API }

// New gets a RecipientLists object ready to use with the specified config.
func New(cfg api.Config) (*RecipientLists, error) {
	rl := &RecipientLists{}
	path := fmt.Sprintf("/api/v%d/recipient-lists", cfg.ApiVersion)
	err := rl.Init(cfg, path)
	if err != nil {
		return nil, err
	}
	return rl, nil
}

// RecipientList is the JSON structure accepted by and returned from the SparkPost Recipient Lists API.
// It's mostly metadata at this level - see Recipients for more detail.
type RecipientList struct {
	ID          string       `json:"id,omitempty"`
	Name        string       `json:"name,omitempty"`
	Description string       `json:"description,omitempty"`
	Attributes  interface{}  `json:"attributes,omitempty"`
	Recipients  *[]Recipient `json:"recipients"`

	Accepted *int `json:"total_accepted_recipients,omitempty"`
}

func (rl *RecipientList) String() string {
	n := 0
	if rl.Recipients != nil {
		n = len(*rl.Recipients)
	} else if rl.Accepted != nil {
		n = *rl.Accepted
	}
	return fmt.Sprintf("ID:\t%s\nName:\t%s\nDesc:\t%s\nCount:\t%d\n",
		rl.ID, rl.Name, rl.Description, n)
}

// Recipient represents one email (you guessed it) recipient.
type Recipient struct {
	Address          interface{} `json:"address"`
	ReturnPath       string      `json:"return_path,omitempty"`
	Tags             []string    `json:"tags,omitempty"`
	Metadata         interface{} `json:"metadata,omitempty"`
	SubstitutionData interface{} `json:"substitution_data,omitempty"`
}

// Address describes the nested object way of specifying the Recipient's email address.
// It can also be a plain string.
type Address struct {
	Email    string `json:"email"`
	Name     string `json:"name,omitempty"`
	HeaderTo string `json:"header_to,omitempty"`
}

// ParseAddress parses the various allowable Content.From values.
func ParseAddress(addr interface{}) (a Address, err error) {
	// handle the allowed types
	switch addrVal := addr.(type) {
	case string: // simple string value
		if addrVal == "" {
			err = fmt.Errorf("Recipient.Address may not be empty")
		} else {
			a.Email = addrVal
		}

	case map[string]interface{}:
		// auto-parsed nested json object
		for k, v := range addrVal {
			switch vVal := v.(type) {
			case string:
				if strings.EqualFold(k, "name") {
					a.Name = vVal
				} else if strings.EqualFold(k, "email") {
					a.Email = vVal
				} else if strings.EqualFold(k, "header_to") {
					a.HeaderTo = vVal
				}
			default:
				err = fmt.Errorf("strings are required for all Recipient.Address values")
				break
			}
		}

	case map[string]string:
		// user-provided json literal (convenience)
		for k, v := range addrVal {
			if strings.EqualFold(k, "name") {
				a.Name = v
			} else if strings.EqualFold(k, "email") {
				a.Email = v
			} else if strings.EqualFold(k, "header_to") {
				a.HeaderTo = v
			}
		}

	default:
		err = fmt.Errorf("unsupported Recipient.Address value type [%s]", reflect.TypeOf(addrVal))
	}

	return
}

// Validate runs sanity checks on a RecipientList struct.
// This should catch most errors before attempting a doomed API call.
func (rl RecipientList) Validate() error {
	// enforce required parameters
	if rl.Recipients == nil || len(*rl.Recipients) <= 0 {
		return fmt.Errorf("RecipientList requires at least one Recipient")
	}

	// enforce max lengths
	if len(rl.ID) > 64 {
		return fmt.Errorf("RecipientList id may not be longer than 64 bytes")
	} else if len(rl.Name) > 64 {
		return fmt.Errorf("RecipientList name may not be longer than 64 bytes")
	} else if len(rl.Description) > 1024 {
		return fmt.Errorf("RecipientList description may not be longer than 1024 bytes")
	}

	var err error
	for _, r := range *rl.Recipients {
		err = r.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

// Validate runs sanity checks on a Recipient struct
// This should catch most errors before attempting a doomed API call.
func (r Recipient) Validate() error {
	_, err := ParseAddress(r.Address)
	if err != nil {
		return err
	}

	// Metadata must be an object, not an array or bool etc.
	if r.Metadata != nil {
		err := api.AssertObject(r.Metadata, "metadata")
		if err != nil {
			return err
		}
	}

	// SubstitutionData must be an object, not an array or bool etc.
	if r.SubstitutionData != nil {
		err := api.AssertObject(r.SubstitutionData, "substitution_data")
		if err != nil {
			return err
		}
	}

	return nil
}

// BuildRecipient accepts a map of key/value pairs, builds, and returns a Recipient object.
// TODO: list expected keys
func (rl RecipientLists) BuildRecipient(p map[string]interface{}) (*Recipient, error) {
	R := &Recipient{}

	// Look up expected keys in the map, deleting as we find them.
	if rpathUntyped, ok := p["return_path"]; ok {
		switch rpath := rpathUntyped.(type) {
		case string:
			R.ReturnPath = rpath
			delete(p, "return_path")
		default:
			return nil, fmt.Errorf("Expected string for `return_path`, got [%s]", reflect.TypeOf(rpath))
		}
	}

	if emailUntyped, ok := p["email"]; ok {
		if R.Address == nil {
			R.Address = Address{}
		}
		switch addr := R.Address.(type) {
		case Address:
			switch email := emailUntyped.(type) {
			case string:
				addr.Email = email
				delete(p, "email")
			default:
				return nil, fmt.Errorf("BuildRecipient expected email as string, not [%s]", reflect.TypeOf(email))
			}
		default:
			return nil, fmt.Errorf("BuildRecipient expected type `Address`, got [%s].", reflect.TypeOf(addr))
		}
	}

	if nameUntyped, ok := p["name"]; ok {
		if R.Address == nil {
			R.Address = Address{}
		}
		switch addr := R.Address.(type) {
		case Address:
			switch name := nameUntyped.(type) {
			case string:
				addr.Name = name
				delete(p, "name")
			default:
				return nil, fmt.Errorf("BuildRecipient expected name as string, not [%s]", reflect.TypeOf(name))
			}
		default:
			return nil, fmt.Errorf("BuildRecipient expected type `Address`, got [%s].", reflect.TypeOf(addr))
		}
	}

	if toUntyped, ok := p["header_to"]; ok {
		if R.Address == nil {
			R.Address = Address{}
		}
		switch addr := R.Address.(type) {
		case Address:
			switch to := toUntyped.(type) {
			case string:
				addr.HeaderTo = to
				delete(p, "header_to")
			default:
				return nil, fmt.Errorf("BuildRecipient expected header_to as string, not [%s]", reflect.TypeOf(to))
			}
		default:
			return nil, fmt.Errorf("BuildRecipient expected type `Address`, got [%s].", reflect.TypeOf(addr))
		}
	}

	if tagsUntyped, ok := p["tags"]; ok {
		switch tags := tagsUntyped.(type) {
		case []interface{}: // auto-parsed tag array
			R.Tags = make([]string, len(tags))
			for idx, tagUntyped := range tags {
				switch tag := tagUntyped.(type) {
				case string:
					R.Tags[idx] = tag
				default:
					return nil, fmt.Errorf("BuildRecipient expected array of tags, got [%s].", reflect.TypeOf(tags))
				}
			}
			delete(p, "tags")

		case []string: // user-provided tag array (convenience)
			R.Tags = make([]string, len(tags))
			for idx, tag := range tags {
				R.Tags[idx] = tag
			}
			delete(p, "tags")
		default:
			return nil, fmt.Errorf("BuildRecipient expected array of tags, got [%s].", reflect.TypeOf(tags))
		}
	}

	if metaUntyped, ok := p["metadata"]; ok {
		err := api.AssertObject(metaUntyped, "metadata")
		if err != nil {
			return nil, err
		}
		R.Metadata = metaUntyped
		delete(p, "metadata")
	}

	if subUntyped, ok := p["substitution_data"]; ok {
		err := api.AssertObject(subUntyped, "substitution_data")
		if err != nil {
			return nil, err
		}
		R.SubstitutionData = subUntyped
		delete(p, "substitution_data")
	}

	// If there are any keys left, they are unsupported.
	if len(p) > 0 {
		return nil, fmt.Errorf("BuildRecipient received unsupported keys")
	}
	return R, nil
}

// BuildRecipients accepts an array of key/value pairs, builds, and returns
// an array of Recipient objects.
func (rl RecipientLists) BuildRecipients(p []interface{}) (*[]Recipient, error) {
	recipients := make([]Recipient, len(p))
	for idx, recipientUntyped := range p {
		switch recipient := recipientUntyped.(type) {
		case map[string]interface{}:
			tmp, err := rl.BuildRecipient(recipient)
			if err != nil {
				return nil, err
			}
			recipients[idx] = *tmp

		default:
			return nil, fmt.Errorf("Build received unexpected recipient format [%s]", reflect.TypeOf(recipient))
		}
	}
	return &recipients, nil
}

// Build accepts a map of key/value pairs, builds, and returns a RecipientList
// object suitable for use with Create.
// TODO: list expected keys
func (rl RecipientLists) Build(p map[string]interface{}) (*RecipientList, error) {
	RL := &RecipientList{}

	// Look up expected keys in the map, deleting as we find them.
	if idUntyped, ok := p["id"]; ok {
		switch id := idUntyped.(type) {
		case string:
			RL.ID = id
			delete(p, "id")
		default:
			return nil, fmt.Errorf("RecipientList.ID must be a string, not [%s]",
				reflect.TypeOf(id))
		}
	}
	if nameUntyped, ok := p["name"]; ok {
		switch name := nameUntyped.(type) {
		case string:
			RL.Name = name
			delete(p, "name")
		default:
			return nil, fmt.Errorf("RecipientList.Name must be a string, not [%s]",
				reflect.TypeOf(name))
		}
	}
	if descUntyped, ok := p["description"]; ok {
		switch desc := descUntyped.(type) {
		case string:
			RL.Description = desc
			delete(p, "description")
		default:
			return nil, fmt.Errorf("RecipientList.Description must be a string not [%s]",
				reflect.TypeOf(desc))
		}
	}
	if attr, ok := p["attributes"]; ok {
		RL.Attributes = attr
		delete(p, "attributes")
	}

	if recipientsUntyped, ok := p["recipients"]; ok {
		switch recipients := recipientsUntyped.(type) {
		case []interface{}:
			var err error
			RL.Recipients, err = rl.BuildRecipients(recipients)
			if err != nil {
				return nil, err
			}

		default:
			return nil, fmt.Errorf("Build received unexpected recipients format [%s]",
				reflect.TypeOf(recipients))
		}
	}

	// If there are any keys left, they are unsupported.
	if len(p) > 0 {
		return nil, fmt.Errorf("Build received unsupported keys")
	}
	return RL, nil
}

// Create accepts a populated RecipientList object, validates it,
// and performs an API call against the configured endpoint.
func (rl RecipientLists) Create(recipList *RecipientList) (id string, err error) {
	if recipList == nil {
		err = fmt.Errorf("Create called with nil RecipientList")
		return
	}

	err = recipList.Validate()
	if err != nil {
		return
	}

	jsonBytes, err := json.Marshal(recipList)
	if err != nil {
		return
	}

	url := fmt.Sprintf("%s%s", rl.Config.BaseUrl, rl.Path)
	res, err := rl.HttpPost(url, jsonBytes)
	if err != nil {
		return
	}

	if err = api.AssertJson(res); err != nil {
		return
	}

	err = rl.ParseResponse(res)
	if err != nil {
		return
	}

	if res.StatusCode == 200 {
		var ok bool
		id, ok = rl.Response.Results["id"]
		if !ok {
			err = fmt.Errorf("Unexpected response to Recipient List creation")
		}

	} else if len(rl.Response.Errors) > 0 {
		// handle common errors
		err = api.PrettyError("RecipientList", "create", res)
		if err != nil {
			return
		}

		if res.StatusCode == 400 || res.StatusCode == 422 {
			eobj := rl.Response.Errors[0]
			err = fmt.Errorf("%s: %s\n%s", eobj.Code, eobj.Message, eobj.Description)
		} else { // everything else
			err = fmt.Errorf("%d: %s", res.StatusCode, string(rl.Response.Body))
		}
	}

	return
}

func (rl RecipientLists) List() (*[]RecipientList, error) {
	url := fmt.Sprintf("%s%s", rl.Config.BaseUrl, rl.Path)
	res, err := rl.HttpGet(url)
	if err != nil {
		return nil, err
	}

	if err = api.AssertJson(res); err != nil {
		return nil, err
	}

	if res.StatusCode == 200 {
		var body []byte
		body, err = rl.ReadBody(res)
		if err != nil {
			return nil, err
		}
		rllist := map[string][]RecipientList{}
		if err = json.Unmarshal(body, &rllist); err != nil {
			return nil, err
		} else if list, ok := rllist["results"]; ok {
			return &list, nil
		}
		return nil, fmt.Errorf("Unexpected response to RecipientList list")

	} else {
		err = rl.ParseResponse(res)
		if err != nil {
			return nil, err
		}
	}

	return nil, err
}
