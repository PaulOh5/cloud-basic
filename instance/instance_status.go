package instance

type InstanceStatus int

const (
	RUNNING InstanceStatus = iota
	STOPPED
)

func (is InstanceStatus) String() string {
	switch is {
	case RUNNING:
		return "RUNNING"
	case STOPPED:
		return "STOPPED"
	default:
		return "UNKNOWN"
	}
}
