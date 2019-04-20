package quota

import (
	"fmt"
	"time"

	"github.com/mdlayher/genetlink"
)

// Generic netlink parameters.
const (
	familyName = "VFS_DQUOT"
	groupName  = "events"
)

// A Client provides access to Linux kobject userspace events. Clients are safe
// for concurrent use.
type Client struct {
	c *genetlink.Conn
}

// New creates a new Client.
func New() (*Client, error) {
	c, err := genetlink.Dial(nil)
	if err != nil {
		return nil, err
	}

	// The genetlink connection must be closed if an error occurs to avoid
	// leaking files.

	group, err := getGroup(c)
	if err != nil {
		_ = c.Close()
		return nil, err
	}

	// genltest does not currently support joining groups, so we have to return
	// the group back to this function to join it for the real connection.
	if err := c.JoinGroup(group); err != nil {
		_ = c.Close()
		return nil, err
	}

	return &Client{c: c}, nil
}

// getGroup is the entry point for tests.
func getGroup(c *genetlink.Conn) (uint32, error) {
	f, err := c.GetFamily(familyName)
	if err != nil {
		return 0, err
	}

	// Determine the ID of the events multicast group.
	var id uint32
	for _, g := range f.Groups {
		if g.Name == groupName {
			id = g.ID
			break
		}
	}
	if id == 0 {
		return 0, fmt.Errorf("quota: could not find %q multicast group", groupName)
	}

	return id, nil
}

// Close releases resources used by a Client.
func (c *Client) Close() error {
	return c.c.Close()
}

// Receive waits until a quota netlink notification is triggered, and then
// returns the Notification.
func (c *Client) Receive() (*Notification, error) {
	msgs, _, err := c.c.Receive()
	if err != nil {
		return nil, err
	}

	if l := len(msgs); l != 1 {
		return nil, fmt.Errorf("quota: expected 1 generic netlink message, but received %d", l)
	}
	if cmd := msgs[0].Header.Command; cmd != cWarning {
		return nil, fmt.Errorf("quota: unexpected generic netlink command: %d", cmd)
	}

	return parseNotification(msgs[0].Data)
}

// SetDeadline sets the read deadline associated with the connection.
func (c *Client) SetDeadline(t time.Time) error {
	return c.c.SetDeadline(t)
}
