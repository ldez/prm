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
	Numbers  PRNumbers `short:"n" description:"PRs numbers."`
	All      bool      `description:"All PR."`
	NoPrompt bool      `short:"d" description:"Disable interactive prompt."`
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

// PRNumbers Slice of PR numbers
type PRNumbers []int

// Set a slice of PR numbers
func (c *PRNumbers) Set(rawValue string) error {
	values := strings.Split(rawValue, ",")
	if len(values) == 0 {
		return fmt.Errorf("bad Value format: %s", rawValue)
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

// Get a slice of PR numbers
func (c *PRNumbers) Get() interface{} { return []int(*c) }

// SetValue of a slice of PR numbers
func (c *PRNumbers) SetValue(val interface{}) {
	*c = val.(PRNumbers)
}

func (c *PRNumbers) String() string {
	var stringNumbers []string
	for _, number := range *c {
		stringNumbers = append(stringNumbers, strconv.Itoa(number))
	}

	return strings.Join(stringNumbers, ",")
}
