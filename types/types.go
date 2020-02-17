package types

// CheckoutOptions "checkout" command options.
type CheckoutOptions struct {
	Number int
}

// RemoveOptions "remove" command options.
type RemoveOptions struct {
	Numbers []int
	All     bool
}

// PushOptions "push" command options.
type PushOptions struct {
	Number int // deprecated
	Force  bool
}

// PullOptions "pull" command options.
type PullOptions struct {
	Force bool
}

// ListOptions "list" command options.
type ListOptions struct {
	All       bool
	SkipEmpty bool
}

// CloneOptions "clone" command options.
type CloneOptions struct {
	NoFork        bool
	Repo          string
	UserAsRootDir bool
	Organization  string
}
