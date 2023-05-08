package sc

type protocol int

const (
	tcp protocol = 0
	udp protocol = 1
)

func (p protocol) String() string {
	switch p {
	case tcp:
		return "tcp"
	case udp:
		return "udp"
	default:
		return "tcp"
	}
}
