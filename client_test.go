package quota

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mdlayher/genetlink"
	"github.com/mdlayher/genetlink/genltest"
	"github.com/mdlayher/netlink"
	"github.com/mdlayher/netlink/nlenc"
	"github.com/mdlayher/netlink/nltest"
)

func TestClientReceive(t *testing.T) {
	tests := []struct {
		name string
		msgs []genetlink.Message
		n    *Notification
	}{
		{
			name: "ok",
			msgs: []genetlink.Message{{
				Header: genetlink.Header{
					Command: cWarning,
				},
				Data: nltest.MustMarshalAttributes([]netlink.Attribute{
					{
						Type: aQType,
						Data: nlenc.Uint32Bytes(uint32(User)),
					},
					{
						Type: aExcessID,
						Data: nlenc.Uint64Bytes(1),
					},
					{
						Type: aWarning,
						Data: nlenc.Uint32Bytes(uint32(InodeHard)),
					},
					{
						Type: aDevMajor,
						Data: nlenc.Uint32Bytes(10),
					},
					{
						Type: aDevMinor,
						Data: nlenc.Uint32Bytes(20),
					},
					{
						Type: aCausedID,
						Data: nlenc.Uint64Bytes(1),
					},
				}),
			}},
			n: &Notification{
				Type:        User,
				ID:          1,
				Warning:     InodeHard,
				DeviceMajor: 10,
				DeviceMinor: 20,
				CausedID:    1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := testClient(t, func(_ genetlink.Message, _ netlink.Message) ([]genetlink.Message, error) {
				return tt.msgs, nil
			})
			defer c.Close()

			n, err := c.Receive()
			if err != nil {
				t.Fatalf("failed to receive: %v", err)
			}

			if diff := cmp.Diff(tt.n, n); diff != "" {
				t.Fatalf("unexpected Notification (-want +got):\n%s", diff)
			}
		})
	}
}

const (
	familyID = 20
	groupID  = 21
)

func testClient(t *testing.T, fn genltest.Func) *Client {
	t.Helper()

	family := genetlink.Family{
		ID:      familyID,
		Version: 1,
		Name:    familyName,
		Groups: []genetlink.MulticastGroup{{
			ID:   groupID,
			Name: groupName,
		}},
	}

	conn := genltest.Dial(genltest.ServeFamily(family, fn))

	group, err := getGroup(conn)
	if err != nil {
		t.Fatalf("failed to get multicast group: %v", err)
	}

	// Expect the same group we entered to be returned, so the real Client is
	// able to join it.
	if diff := cmp.Diff(groupID, int(group)); diff != "" {
		t.Fatalf("unexpected multicast group ID (-want +got):\n%s", diff)
	}

	return &Client{c: conn}
}
