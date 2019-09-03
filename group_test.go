package supervisor

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestSupervisor(t *testing.T) {
	checkErr := func(err error, msg string) bool {
		if err == nil {
			return msg == ""
		}

		return strings.Contains(err.Error(), msg)
	}

	tests := [...]struct {
		behave string
		do     func() error
		expect string
	}{
		{
			// two agents report one error
			"one-err",
			func() error {
				s, _ := WithContext(context.Background())
				s.Agent(func() error { return errors.New("test1") })
				s.Agent(func() error { return errors.New("test2") })
				return <-s.Err()
			},
			"test",
		},
		{
			// normal termination returns empty error
			"normal-term",
			func() error {
				s, _ := WithContext(context.Background())
				s.Agent(func() error { return nil })
				return <-s.Err()
			},
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.behave, func(t *testing.T) {
			if err := tt.do(); !checkErr(err, tt.expect) {
				t.Errorf("expected %s, got %s", tt.expect, err)
			}
		})
	}
}
