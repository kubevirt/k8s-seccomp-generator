package install

var SUPPORTED_DISTROS [1]Distro = [1]Distro{Centos_stream8}

type Distro int

const (
	Centos_stream8 Distro = 0
)

// the returned string must match the tag (nithishdev/falco-loader:$tag) used in the loader dockerfile to make sure
// that we can seamless use the ENUMS for dynamically deploying the manifests
func (d Distro) String() string {
	switch d {
	case Centos_stream8:
		return "centos-stream8"
	}
	return ""
}

func DistroFromString(distro string) Distro {
	switch distro {
	case "centos-stream8":
		return Centos_stream8
	}
	return -1
}
