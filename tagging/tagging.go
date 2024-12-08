package tagging

import (
	"fmt"
	"strings"

	"github.com/coreos/go-semver/semver"
)

func splitOff(input string, delim string) string {
	parts := strings.SplitN(input, delim, 2)

	if len(parts) == 2 {
		return parts[0]
	}

	return input
}

// DefaultTags returns a set of default suggested tags.
func DefaultTags(ref string, sha string) []string {
	if !strings.HasPrefix(ref, "refs/tags/") {
		return []string{sha}
	}

	v := stripTagPrefix(ref)

	version, err := semver.NewVersion(v)
	if err != nil {
		return []string{sha}
	}

	if version.PreRelease != "" || version.Metadata != "" {
		return []string{
			sha,
			"latest",
			version.String(),
		}
	}

	v = stripTagPrefix(ref)
	v = splitOff(splitOff(v, "+"), "-")
	dotParts := strings.SplitN(v, ".", 3)

	if version.Major == 0 {
		return []string{
			sha,
			"latest",
			fmt.Sprintf("%0*d.%0*d", len(dotParts[0]), version.Major, len(dotParts[1]), version.Minor),
			fmt.Sprintf("%0*d.%0*d.%0*d", len(dotParts[0]), version.Major, len(dotParts[1]), version.Minor, len(dotParts[2]), version.Patch),
		}
	}

	return []string{
		sha,
		"latest",
		fmt.Sprintf("%0*d", len(dotParts[0]), version.Major),
		fmt.Sprintf("%0*d.%0*d", len(dotParts[0]), version.Major, len(dotParts[1]), version.Minor),
		fmt.Sprintf("%0*d.%0*d.%0*d", len(dotParts[0]), version.Major, len(dotParts[1]), version.Minor, len(dotParts[2]), version.Patch),
	}
}

// UseDefaultTag to restrict latest tag for default branch.
func UseDefaultTag(ref, defaultBranch string) bool {
	if strings.HasPrefix(ref, "refs/tags/") || strings.HasPrefix(ref, "refs/pull/") {
		return true
	}

	if stripHeadPrefix(ref) == defaultBranch {
		return true
	}

	return false
}

// stripHeadPrefix just strips the ref heads prefix.
func stripHeadPrefix(ref string) string {
	return strings.TrimPrefix(ref, "refs/heads/")
}

// stripTagPrefix just strips the ref tags prefix.
func stripTagPrefix(ref string) string {
	ref = strings.TrimPrefix(ref, "refs/tags/")
	ref = strings.TrimPrefix(ref, "v")

	return ref
}
