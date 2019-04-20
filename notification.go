package quota

import "github.com/mdlayher/netlink"

// A Notification is a disk quota notification.
type Notification struct {
	Type        Type
	ID          int
	Warning     Warning
	DeviceMajor int
	DeviceMinor int
	CausedID    int
}

//go:generate stringer -type=Type,Warning -output=string.go

// A Type is a quota type.
type Type int

// Possible quota types.
const (
	User    Type = 0
	Group   Type = 1
	Project Type = 2
)

// A Warning is an individual event which caused a Notification to be sent.
type Warning int

// Possible Warning values. See the Linux quota documentation for details.
const (
	None           Warning = 0
	InodeHard      Warning = 1
	InodeSoftLong  Warning = 2
	InodeSoft      Warning = 3
	BlockHard      Warning = 4
	BlockSoftLong  Warning = 5
	BlockSoft      Warning = 6
	InodeHardBelow Warning = 7
	InodeSoftBelow Warning = 8
	BlockHardBelow Warning = 9
	BlockSoftBelow Warning = 10
)

// Constants taken from Linux kernel headers:
// https://github.com/torvalds/linux/blob/master/include/uapi/linux/quota.h.
//
// TODO(mdlayher): get c-for-go working or put in x/sys/unix.
const (
	cWarning = 1

	aQType    = 1
	aExcessID = 2
	aWarning  = 3
	aDevMajor = 4
	aDevMinor = 5
	aCausedID = 6
)

// parseNotification parses netlink attribute bytes into a Notification.
func parseNotification(b []byte) (*Notification, error) {
	ad, err := netlink.NewAttributeDecoder(b)
	if err != nil {
		return nil, err
	}

	var n Notification
	for ad.Next() {
		switch ad.Type() {
		case aQType:
			n.Type = Type(ad.Uint32())
		case aExcessID:
			n.ID = int(ad.Uint64())
		case aWarning:
			n.Warning = Warning(ad.Uint32())
		case aDevMajor:
			n.DeviceMajor = int(ad.Uint32())
		case aDevMinor:
			n.DeviceMinor = int(ad.Uint32())
		case aCausedID:
			n.CausedID = int(ad.Uint64())
		}
	}

	if err := ad.Err(); err != nil {
		return nil, err
	}

	return &n, nil
}
