package owasp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

// NewEditor returns an empty editor. Use Load or Import methods to
// fill with entries.
func NewEditor() *Editor {
	return &Editor{}
}

type Editor struct {
	Entries []Entry
}

// SetApplicable sets the applicable field of given entry. Returns
// error if no entry is found. The given pattern is interpreted as a
// regular expression if it's a string starts with a `^`. If it contains a `*`
// without backslash it's converted to a regular expression, otherwise
// the id is matched as is. The pattern may also be a Level.
func (me *Editor) SetApplicable(pattern interface{}, v bool) error {
	match := matcherFrom(pattern)
	var found bool
	for i, e := range me.Entries {
		if match(e) {
			me.Entries[i].Applicable = v
			found = true
		}
	}
	if !found {
		return fmt.Errorf("no entries matched by %s", pattern)
	}
	return nil
}

func matcherFrom(v interface{}) func(e Entry) bool {
	if level, ok := v.(Level); ok {
		return func(e Entry) bool {
			return e.IsLevel(level)
		}
	}
	{
		v := v.(string)

		switch {
		case len(v) > 0 && v[0] == '^':
			rx := regexp.MustCompile(v)
			return func(e Entry) bool {
				return rx.Match([]byte(e.ID))
			}
		case strings.Contains(v, "*"):
			v = strings.ReplaceAll(v, ".", `\.`)
			v = strings.ReplaceAll(v, "*", `.*`)
			v = "^" + v + "$"

			rx := regexp.MustCompile(v)
			return func(e Entry) bool {
				return rx.Match([]byte(e.ID))
			}

		default:
			// exact id match
			return func(e Entry) bool { return e.ID == v }
		}
	}
}

// ----------------------------------------

// Reset same as calling ResetVerified and ResetApplicable
func (me *Editor) Reset() {
	me.ResetVerified()
	me.ResetApplicable()
}

// ResetVerified sets Verified field to false on all entries and Manual to nil.
func (me *Editor) ResetVerified() {
	for i := range me.Entries {
		me.Entries[i].Verified = false
		me.Entries[i].Manual = nil
	}
}

// ResetApplicable sets applicable field to false on all entries.
func (me *Editor) ResetApplicable() {
	for i := range me.Entries {
		me.Entries[i].Applicable = false
	}
}

// ----------------------------------------

// SetVerified sets the given entry as verified. Returns error if
// pattern is not found or the entry is not applicable. See
// SetApplicable for pattern variations.
func (me *Editor) SetVerified(pattern interface{}, v bool, man ...Manual) error {
	match := matcherFrom(pattern)
	var found bool
	for i, e := range me.Entries {
		if match(e) {
			if !e.Applicable {
				return fmt.Errorf("%v is not applicable", e.ID)
			}
			me.Entries[i].Verified = v
			me.Entries[i].Manual = nil
			if len(man) > 0 {
				m := man[0]
				me.Entries[i].Manual = &m
			}
			found = true
		}
	}
	if !found {
		return fmt.Errorf("no entries matched by %s", pattern)
	}
	return nil
}

// ----------------------------------------

// Load entries from given json file.
func (me *Editor) Load(filename string) error {
	fh, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fh.Close()
	return me.Import(fh)
}

// SaveAs writes entries as a tidy json to the given filename.
func (me *Editor) SaveAs(filename string) error {
	fh, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fh.Close()
	_, err = me.WriteTo(fh)
	return err
}

// Import entries from json
func (me *Editor) Import(r io.Reader) error {
	return json.NewDecoder(r).Decode(&me.Entries)
}

// WriteTo exports entries as tidy json to the given writer.
func (me *Editor) WriteTo(w io.Writer) (int64, error) {
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(me.Entries)

	var tidy bytes.Buffer
	_ = json.Indent(&tidy, buf.Bytes(), "", "  ")

	return io.Copy(w, &tidy)
}

// NewReport returns a new report from the loaded entries.
func (me *Editor) NewReport(title string) *Report {
	r := &Report{
		Title:              title,
		ShortDescriptionNA: true,
		ShowNonApplicable:  true,
	}
	r.AddEntries(me.Entries...)
	return r
}
