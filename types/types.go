package types

// NoOption empty struct.
type NoOption struct{}

// CheckoutOptions "checkout" command options.
type CheckoutOptions struct {
	Number int `short:"n" description:"PR number."`
}

// RemoveOptions "remove" command options.
type RemoveOptions struct {
	Number int  `short:"n" description:"PR number."`
	All    bool `description:"All PR."`
}

// PushOptions "push" command options.
type PushOptions struct {
	Number int  `short:"n" description:"PR number (optional: parse the branch name)."`
	Force  bool `short:"f" description:"Force the push."`
}

// ListOptions "list" command options.
type ListOptions struct {
	All bool `description:"All PR."`
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
