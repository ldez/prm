package types

import (
	"fmt"
	"strconv"
	"strings"
)

// NoOption empty struct.
type NoOption struct{}

// CheckoutOptions "checkout" command options.
type CheckoutOptions struct {
	Number int `short:"n" description:"PR number."`
}

// RemoveOptions "remove" command options.
type RemoveOptions struct {
	Numbers PRNumbers `short:"n" description:"PRs numbers."`
	All     bool      `description:"All PR."`
}

// PushOptions "push" command options.
type PushOptions struct {
	Number int  `short:"n" description:"PR number (optional: parse the branch name)."`
	Force  bool `short:"f" description:"Force the push."`
}

// PullOptions "pull" command options.
type PullOptions struct {
	Force bool `short:"f" description:"Force the push."`
}

// ListOptions "list" command options.
type ListOptions struct {
	All bool `description:"All PR."`
}

type PRNumbers []int

func (c *PRNumbers) Set(rawValue string) error {
	values := strings.Split(rawValue, ",")
	if len(values) == 0 {
		return fmt.Errorf("Bad Value format: %s", rawValue)
	}
	for _, value := range values {
		number, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return err
		}
		*c = append(*c, int(number))
	}
	return nil
}

func (c *PRNumbers) Get() interface{} { return []int(*c) }

func (c *PRNumbers) String() string {

	stringNumbers := []string{}
	for _, number := range *c {
		stringNumbers = append(stringNumbers, strconv.FormatInt(int64(number), 10))
	}

	return strings.Join(stringNumbers, ",")
}

func (c *PRNumbers) SetValue(val interface{}) {
	*c = PRNumbers(val.(PRNumbers))
}

// Repository Git repository model.
type Repository struct {
	Owner string
	Name  string
}

// Remote Git remote model.
type Remote struct {
	Name string
	URL  string
}

// ByRemoteName sort remote by name.
type ByRemoteName []Remote

func (r ByRemoteName) Len() int           { return len(r) }
func (r ByRemoteName) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r ByRemoteName) Less(i, j int) bool { return r[i].Name < r[j].Name }
