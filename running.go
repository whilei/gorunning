// go_running is a package to the get the file path of running program by process id.
// It:
// - FIXME: probably does not support Windows or other edge-case OS's
// - depends on existence of possible some but not all standard programs like lsof, ps, and awk
// - FIXME: depends on bash shell
// - TODO: could return more than just the file path, like time since started, owner, group...

package go_running

import (
	"os/exec"
	"regexp"
	"path/filepath"
	"errors"
	"fmt"
)

const (
	followSymlinksDefault = true
)

var (
	errInvalidArg = errors.New("invalid argument(s)")
	errRunningPathUnavailable = errors.New("path of running file unavailable")
	errBadPid = errors.New("bad pid")
	errCantCast = errors.New("cant cast value")
)

var (
	// each should accept exactly one '%d' corresponding to PID
	runningGrabbers = []string{
		`lsof -p %d | awk '$4 == "txt" { print $9 }' | head -1`,
		`ps awx | awk '$1 == %d { print $5 }'`,
	}
	shells = map[string]string{
		"bash": "-c",
	}
)

func parseArbitraryArgToBool(normal bool, arguments ...interface{}) (b bool, err error) {
	b = normal
	if len(arguments) == 0 {
		return normal, nil
	} else if len(arguments) > 1 {
		return false, errInvalidArg
	}

	if castBool, ok := arguments[0].(bool); ok {
		b = castBool
	} else {
		return false, errCantCast
	}
	return b, nil
}

// GetPath gets the path of a running file by PID.
// If argFollowSymlinkInterface is nil (or empty), the default is true.
func GetPath(pid int, argFollowSymlink ...interface{}) (runningPath string, err error) {
	if pid <= 0 {
		return "", errBadPid
	}

	followSymlink := true
	followSymlink, err = parseArbitraryArgToBool(followSymlinksDefault, argFollowSymlink...)
	if err != nil {
		return "", errInvalidArg
	}

	return getRunningFilepath(pid, followSymlink, runningGrabbers)
}

func getRunningFilepath(pid int, followSymlink bool, grabbers []string) (runningPath string, err error) {
	var cmdOutput []byte
	err = errRunningPathUnavailable

	OUTER:
	for _, cmd := range grabbers {
		for p, arg := range shells {
			cmdOutput, err = exec.Command(p, arg, fmt.Sprintf(cmd, pid)).CombinedOutput()
			if err == nil {
				break OUTER
			}
			// ignore errors because they can maybe be expected to happen
		}
	}
	if err != nil {
		return "", err
	}
	// Remove any newlines from stdout/err
	re := regexp.MustCompile(`\r?\n`)
	outputString := re.ReplaceAllString(string(cmdOutput), "")

	var fullPath string

	fullPath, err = filepath.Abs(outputString)
	if err != nil {
		return "", err
	}

	if followSymlink {
		fullPath, err = filepath.EvalSymlinks(fullPath)
		if err != nil {
			return "", err
		}
	}
	return fullPath, err
}
