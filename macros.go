package gosparkpost

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

// Macro enables user-defined functions to run on Template.Content before sending to the SparkPost API.
// This enables e.g. external content, locale-specific date formatting, and so on.
type Macro struct {
	Name string
	Func func(string) string
}

const (
	StaticToken = iota
	MacroToken
)

// TokenType differentiates between static content and macros that require extra processing.
type TokenType int

// ContentToken represents a piece of content, with one of the types defined above.
type ContentToken struct {
	Type TokenType
	Text string
}

var wordChars = regexp.MustCompile(`^\w+$`)

// RegisterMacro associates a Macro with a Client.
// The provided Macro.Func is wrapped so that its argument is the string between the Macro name and closing delimiter.
// Splitting of arguments can be done in the Macro Func if desired.
// Everything between the Macro Name and the closing delimiter will be passed as the single string argument.
func (c *Client) RegisterMacro(m *Macro) error {
	if m == nil {
		return errors.New(`can't add nil Macro`)
	} else if !wordChars.MatchString(m.Name) {
		return errors.New(`Macro names must only contain \w characters`)
	} else if m.Func == nil {
		return errors.New(`Macro must have non-nil Func field`)
	}

	if c.macros == nil {
		c.macros = map[string]Macro{}
	}
	c.macros[m.Name] = *m
	return nil
}

// Apply substitutes top-level string values from the Recipient's SubstitutionData and Metadata
// (in that order) for placeholders in the provided string.
func (r *Recipient) Apply(in string) (string, error) {
	tokens, err := Tokenize(in)
	if err != nil {
		return "", err
	}
	chunks := make([]string, len(tokens))

	addr, err := ParseAddress(r.Address)
	if err != nil {
		return "", errors.Wrap(err, "parsing recipient address")
	}

	var sub, meta map[string]interface{}
	var ok bool
	if r.SubstitutionData != nil {
		if sub, ok = r.SubstitutionData.(map[string]interface{}); !ok {
			return "", errors.Errorf("unexpected substitution data type for recipient %s", addr.Email)
		}
	}
	if r.Metadata != nil {
		if meta, ok = r.Metadata.(map[string]interface{}); !ok {
			return "", errors.Errorf("unexpected metadata type for recipient %s", addr.Email)
		}
	}

	for idx, token := range tokens {
		switch token.Type {
		case StaticToken:
			chunks[idx] = token.Text

		case MacroToken:
			key := strings.TrimSpace(strings.Trim(token.Text, "{}"))
			for _, subst := range []map[string]interface{}{sub, meta} {
				if ival, ok := subst[key]; ok {
					switch val := ival.(type) {
					case string:
						chunks[idx] = val
					default:
						chunks[idx] = token.Text
					}
					break
				}
			}
		}
	}

	if len(chunks) == 1 {
		return chunks[0], nil
	}
	return strings.Join(chunks, ""), nil
}

// ApplyMacros runs all Macros registered with the Client against the provided string, returning the result.
// If a Recipient is provided, substitution is performed on the macro parameter before the macro runs.
// Any placeholders not handled by a macro are left intact.
func (c *Client) ApplyMacros(in string, r *Recipient) (string, error) {
	if c.macros == nil {
		// if no macros are defined, this is a no-op
		return in, nil
	}

	tokens, err := Tokenize(in)
	if err != nil {
		return "", err
	}
	chunks := make([]string, len(tokens))

	for idx, token := range tokens {
		switch token.Type {
		case StaticToken:
			chunks[idx] = token.Text

		case MacroToken:
			body := strings.TrimSpace(strings.Trim(token.Text, "{}"))
			// split off macro name
			atoms := strings.SplitN(body, " ", 2)
			if m, ok := c.macros[atoms[0]]; ok {
				var params string
				if len(atoms) == 2 {
					params = atoms[1]
				} else {
					params = ""
				}
				if r != nil {
					params, err = r.Apply(params)
					if err != nil {
						return "", err
					}
				}
				chunks[idx] = m.Func(params)
			} else {
				// no client macro matches this block, pass it through
				chunks[idx] = token.Text
			}
		}
	}

	if len(chunks) == 1 {
		return chunks[0], nil
	}
	return strings.Join(chunks, ""), nil
}

func Tokenize(str string) (out []ContentToken, err error) {
	strlen := len(str)
	for {
		open := strings.Index(str, "{{")
		if open >= 0 && open < strlen {
			if open > 0 {
				// we have a macro, make a token with the static text leading up to it
				out = append(out, ContentToken{Text: str[:open]})
				str = str[open:]
				strlen -= open
			} else {
				// Do nothing if macro starts at index 0,
				// otherwise we end up with blank StaticTokens
			}
		} else {
			break
		}

		// advance to the end of the macro
		curlies := 0
		var last int
		for last = 0; last < strlen; last++ {
			switch str[last] {
			case '{':
				curlies++
			case '}':
				curlies--
			}
			if curlies == 0 {
				last++
				break
			}
		}

		if curlies != 0 {
			return nil, errors.Errorf("mismatched curly braces near %q", str)
		}

		out = append(out, ContentToken{
			Type: MacroToken,
			Text: str[:last],
		})

		str = str[last:]
		strlen -= last
	}
	if strlen > 0 {
		out = append(out, ContentToken{Text: str})
	}
	return out, nil
}
