package sync

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func jsonResponse(body string, status int) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": {"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func TestPagerDutyScheduleEmails(t *testing.T) {
	for _, lookahead := range []time.Duration{time.Second, 2400 * time.Hour} {
		t.Run(lookahead.String(), func(t *testing.T) {
			p := newPagerDutyClient("pagerduty-token")
			p.client.HTTPClient = &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, "/schedules/schedule-1/users", r.URL.Path)
				assert.Equal(t, "Token token=pagerduty-token", r.Header.Get("Authorization"))
				since, err := time.Parse(time.RFC3339, r.URL.Query().Get("since"))
				require.NoError(t, err)
				until, err := time.Parse(time.RFC3339, r.URL.Query().Get("until"))
				require.NoError(t, err)
				assert.WithinDuration(t, time.Now().UTC(), since, 2*time.Second)
				assert.WithinDuration(t, since.Add(lookahead), until, time.Second)
				return jsonResponse(`{"users":[{"email":"alice@example.com"},{"email":"bob@example.com"}]}`, http.StatusOK), nil
			})}

			emails, err := p.getEmailsForSchedule("schedule-1", lookahead)

			require.NoError(t, err)
			assert.Equal(t, []string{"alice@example.com", "bob@example.com"}, emails)
		})
	}
}

func TestPagerDutyScheduleErrors(t *testing.T) {
	for _, tc := range []struct {
		name   string
		body   string
		status int
	}{
		{"unauthorized", `{"error":{"message":"Invalid credentials","code":2006}}`, http.StatusUnauthorized},
		{"invalid JSON", `invalid`, http.StatusOK},
		{"missing users", `{}`, http.StatusOK},
	} {
		t.Run(tc.name, func(t *testing.T) {
			p := newPagerDutyClient("pagerduty-token")
			p.client.HTTPClient = &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				return jsonResponse(tc.body, tc.status), nil
			})}
			emails, err := p.getEmailsForSchedule("schedule-1", time.Second)
			assert.Error(t, err)
			assert.Nil(t, emails)
		})
	}
}

// The production Slack constructor uses http.DefaultTransport. Replacing it
// exercises the real SDK without allowing network calls. These tests are serial.
func stubSlackTransport(t *testing.T, handler roundTripFunc) {
	previous := http.DefaultTransport
	http.DefaultTransport = handler
	t.Cleanup(func() { http.DefaultTransport = previous })
}

func TestSlackUsersAndGroups(t *testing.T) {
	var calls []string
	stubSlackTransport(t, func(r *http.Request) (*http.Response, error) {
		require.NoError(t, r.ParseForm())
		assert.Equal(t, "slack.com", r.URL.Host)
		assert.Equal(t, "slack-token", r.Form.Get("token"))
		calls = append(calls, r.URL.Path)
		switch r.URL.Path {
		case "/api/usergroups.list":
			return jsonResponse(`{"ok":true,"usergroups":[{"id":"G1","name":"current-oncall-platform"}]}`, http.StatusOK), nil
		case "/api/users.list":
			if r.Form.Get("cursor") == "" {
				return jsonResponse(`{"ok":true,"members":[{"id":"U1","profile":{"email":"alice@example.com"}}],"response_metadata":{"next_cursor":"page-2"}}`, http.StatusOK), nil
			}
			assert.Equal(t, "page-2", r.Form.Get("cursor"))
			return jsonResponse(`{"ok":true,"members":[{"id":"U2","profile":{"email":"bob@example.com"}}],"response_metadata":{"next_cursor":""}}`, http.StatusOK), nil
		case "/api/usergroups.create":
			assert.Equal(t, "all-oncall-platforms", r.Form.Get("name"))
			assert.Equal(t, "all-oncall-platforms", r.Form.Get("handle"))
			return jsonResponse(`{"ok":true,"usergroup":{"id":"G2","name":"all-oncall-platforms"}}`, http.StatusOK), nil
		case "/api/usergroups.users.list":
			assert.Equal(t, "G2", r.Form.Get("usergroup"))
			return jsonResponse(`{"ok":true,"users":["U1"]}`, http.StatusOK), nil
		case "/api/usergroups.users.update":
			assert.Equal(t, "G2", r.Form.Get("usergroup"))
			assert.Equal(t, "U1,U2", r.Form.Get("users"))
			return jsonResponse(`{"ok":true,"usergroup":{"id":"G2","users":["U1","U2"]}}`, http.StatusOK), nil
		default:
			t.Errorf("unexpected Slack request: %s", r.URL.Path)
			return nil, fmt.Errorf("unexpected Slack request")
		}
	})

	s, err := newSlackClient("slack-token")
	require.NoError(t, err)
	ids, err := s.getSlackIDsFromEmails([]string{"ALICE@example.com", "bob@example.com"})
	require.NoError(t, err)
	assert.Equal(t, []string{"U1", "U2"}, ids)
	_, err = s.getSlackIDsFromEmails([]string{"missing@example.com"})
	assert.EqualError(t, err, "could not find slack user with email: missing@example.com")

	existing, err := s.createOrGetUserGroup("CURRENT-ONCALL-PLATFORM")
	require.NoError(t, err)
	assert.Equal(t, "G1", existing.ID)
	created, err := s.createOrGetUserGroup("all-oncall-platforms")
	require.NoError(t, err)
	assert.Equal(t, "G2", created.ID)
	members, err := s.Client.GetUserGroupMembers(created.ID)
	require.NoError(t, err)
	assert.Equal(t, []string{"U1"}, members)
	updated, err := s.Client.UpdateUserGroupMembers(created.ID, strings.Join(ids, ","))
	require.NoError(t, err)
	assert.Equal(t, []string{"U1", "U2"}, updated.Users)
	assert.Equal(t, []string{
		"/api/usergroups.list", "/api/users.list", "/api/users.list",
		"/api/usergroups.create", "/api/usergroups.users.list", "/api/usergroups.users.update",
	}, calls)
}

func TestSlackInitializationErrors(t *testing.T) {
	for _, endpoint := range []string{"/api/usergroups.list", "/api/users.list"} {
		t.Run(endpoint, func(t *testing.T) {
			stubSlackTransport(t, func(r *http.Request) (*http.Response, error) {
				if r.URL.Path == endpoint {
					return jsonResponse(`{"ok":false,"error":"invalid_auth"}`, http.StatusOK), nil
				}
				return jsonResponse(`{"ok":true,"usergroups":[]}`, http.StatusOK), nil
			})
			s, err := newSlackClient("slack-token")
			assert.EqualError(t, err, "invalid_auth")
			assert.Nil(t, s)
		})
	}
}
