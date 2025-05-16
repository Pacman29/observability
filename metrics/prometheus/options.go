package prometheus

type Option func(*options)

type options struct {
	namespace           string
	subsystem           string
	labelsPoolCapSave   int
	labelsPoolCapCreate int
	metricsSize         int
}

func newOptions() *options {
	return &options{
		namespace:           "",
		subsystem:           "",
		labelsPoolCapSave:   10,
		labelsPoolCapCreate: 10,
		metricsSize:         10,
	}
}

func WithNamespace(n string) Option {
	return func(o *options) {
		o.namespace = n
	}
}

func WithSubsystem(s string) Option {
	return func(o *options) {
		o.subsystem = s
	}
}

func WithLabelsPoolCapCreate(s int) Option {
	return func(o *options) {
		o.labelsPoolCapCreate = s
	}
}

func WithLabelsPoolCapSave(s int) Option {
	return func(o *options) {
		o.labelsPoolCapSave = s
	}
}

func WithMetricsSize(s int) Option {
	return func(o *options) {
		o.metricsSize = s
	}
}
