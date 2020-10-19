package mqi

// Topic struct
type topic struct {
	name string
}

// Topic interface
type Topic interface {
	Name() string
	WithName(name string) Topic
}

// NewTopic constructor
func NewTopic(name string) Topic { return topic{}.WithName(name) }

// Name getter
func (tp topic) Name() string { return tp.name }

// WithName setter
func (tp topic) WithName(name string) Topic {
	tp.name = name
	return tp
}
