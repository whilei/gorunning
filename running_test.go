package go_running

import (
	"testing"
	"os"
)

func TestGetPath(t *testing.T) {
	cases := []struct{
		pid int
		follow bool
		ok bool
	}{
		{os.Getpid(), true, true},
		{os.Getpid(), false, true},
		{0, true, false},
	}
	for _, c := range cases {
		p, err := GetPath(c.pid, c.follow)
		if (err != nil && c.ok) || (c.ok && p == "" ) {
			t.Errorf("error: %s, %v / %d %v %v", p, err, c.pid, c.follow, c.ok)
		} else if c.ok {
			t.Log(p)
		}
	}
}

func testgetRunningFilepath(t *testing.T, follow bool, grabbers []string) {
	pid := os.Getpid()
	p, err := getRunningFilepath(pid, follow, grabbers)
	if err != nil {
		t.Error("err:", err)
	} else {
		t.Log("ok:", p)
	}
}

// I'm not really sure about the symlinks. Maybe it will fail sometimes.
func TestGetRunningFilepath(t *testing.T) {
	for _, sym := range []bool{true, false} {
		for _, g := range runningGrabbers {
			gg := []string{g}
			testgetRunningFilepath(t, sym, gg)
		}
		testgetRunningFilepath(t, sym, runningGrabbers)
	}
}

func TestParseArbitraryArgToBool(t *testing.T) {
	cases := []struct{
		ok bool
		expected bool
		args []interface{}
	}{
		{true, followSymlinksDefault, nil},
		{true, true, []interface{}{true}},
		{true, false, []interface{}{false}},
		{false, false, []interface{}{true, true}},
		{false,false, []interface{}{123}},
		{false, false,[]interface{}{"nogo"}},
	}
	for _, c := range cases {
		b, err := parseArbitraryArgToBool(followSymlinksDefault, c.args...)
		if err != nil && c.ok {
			t.Error(c)
		}
		if b != c.expected {
			t.Errorf("want: %v, got: %v", c.expected, b)
		}
	}
}