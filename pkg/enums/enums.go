package enums

type CacheKey int

const (
	_ CacheKey = iota
	JOB

)

func (i CacheKey) String() string {
	switch i {
	case JOB:
		return "hrjob"

	default:
		return ""
	}
}
