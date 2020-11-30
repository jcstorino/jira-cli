package query

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type paramsErr struct {
	latest     bool
	watching   bool
	resolution bool
	issueType  bool
}

type testFlagParser struct {
	err        paramsErr
	noHistory  bool
	noWatching bool
	orderDesc  bool
	emptyType  bool
}

func (tfp testFlagParser) GetBool(name string) (bool, error) {
	if tfp.err.latest && name == "latest" {
		return false, fmt.Errorf("oops! couldn't fetch latest flag")
	}

	if tfp.err.watching && name == "watching" {
		return false, fmt.Errorf("oops! couldn't fetch watching flag")
	}

	if tfp.noHistory && name == "latest" {
		return false, nil
	}

	if tfp.noWatching && name == "watching" {
		return false, nil
	}

	if tfp.orderDesc && name == "reverse" {
		return false, nil
	}

	return true, nil
}

func (tfp testFlagParser) GetString(name string) (string, error) {
	if tfp.err.resolution && name == "resolution" {
		return "", fmt.Errorf("oops! couldn't fetch resolution flag")
	}

	if tfp.err.issueType && name == "type" {
		return "", fmt.Errorf("oops! couldn't fetch type flag")
	}

	if tfp.emptyType && name == "type" {
		return "", nil
	}

	return "test", nil
}

func TestIssueGet(t *testing.T) {
	cases := []struct {
		name       string
		initialize func() *Issue
		expected   string
	}{
		{
			name: "query with default parameters",
			initialize: func() *Issue {
				i, err := NewIssue("TEST", &testFlagParser{})
				assert.NoError(t, err)

				return i
			},
			expected: `project="TEST" AND issue IN issueHistory() AND issue IN watchedIssues() AND ` +
				`type="test" AND resolution="test" AND status="test" AND priority="test" AND reporter="test" AND assignee="test" ` +
				`ORDER BY lastViewed ASC`,
		},
		{
			name: "query without issue history parameter",
			initialize: func() *Issue {
				i, err := NewIssue("TEST", &testFlagParser{noHistory: true})
				assert.NoError(t, err)

				return i
			},
			expected: `project="TEST" AND issue IN watchedIssues() AND ` +
				`type="test" AND resolution="test" AND status="test" AND priority="test" AND reporter="test" AND assignee="test" ` +
				`ORDER BY created ASC`,
		},
		{
			name: "query only with fields filter",
			initialize: func() *Issue {
				i, err := NewIssue("TEST", &testFlagParser{noHistory: true, noWatching: true})
				assert.NoError(t, err)

				return i
			},
			expected: `project="TEST" AND ` +
				`type="test" AND resolution="test" AND status="test" AND priority="test" AND reporter="test" AND assignee="test" ` +
				`ORDER BY created ASC`,
		},
		{
			name: "query with error when fetching latest flag",
			initialize: func() *Issue {
				i, err := NewIssue("TEST", &testFlagParser{err: paramsErr{
					latest: true,
				}})
				assert.Error(t, err)

				return i
			},
			expected: "",
		},
		{
			name: "query with error when fetching watching flag",
			initialize: func() *Issue {
				i, err := NewIssue("TEST", &testFlagParser{err: paramsErr{
					watching: true,
				}})
				assert.Error(t, err)

				return i
			},
			expected: "",
		},
		{
			name: "query with error when fetching resolution flag",
			initialize: func() *Issue {
				i, err := NewIssue("TEST", &testFlagParser{err: paramsErr{
					resolution: true,
				}})
				assert.Error(t, err)

				return i
			},
			expected: "",
		},
		{
			name: "query with error when fetching type flag",
			initialize: func() *Issue {
				i, err := NewIssue("TEST", &testFlagParser{err: paramsErr{
					issueType: true,
				}})
				assert.Error(t, err)

				return i
			},
			expected: "",
		},
		{
			name: "query without issue type flag",
			initialize: func() *Issue {
				i, err := NewIssue("TEST", &testFlagParser{emptyType: true})
				assert.NoError(t, err)

				return i
			},
			expected: `project="TEST" AND issue IN issueHistory() AND issue IN watchedIssues() AND ` +
				`resolution="test" AND status="test" AND priority="test" AND reporter="test" AND assignee="test" ` +
				`ORDER BY lastViewed ASC`,
		},
		{
			name: "query with reverse set to true",
			initialize: func() *Issue {
				i, err := NewIssue("TEST", &testFlagParser{orderDesc: true})
				assert.NoError(t, err)

				return i
			},
			expected: `project="TEST" AND issue IN issueHistory() AND issue IN watchedIssues() AND ` +
				`type="test" AND resolution="test" AND status="test" AND priority="test" AND reporter="test" AND assignee="test" ` +
				`ORDER BY lastViewed DESC`,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			q := tc.initialize()

			if q != nil {
				assert.Equal(t, tc.expected, q.Get())
			}
		})
	}
}