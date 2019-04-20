# quota [![builds.sr.ht status](https://builds.sr.ht/~mdlayher/quota.svg)](https://builds.sr.ht/~mdlayher/quota?) [![GoDoc](https://godoc.org/github.com/mdlayher/quota?status.svg)](https://godoc.org/github.com/mdlayher/quota) [![Go Report Card](https://goreportcard.com/badge/github.com/mdlayher/quota)](https://goreportcard.com/report/github.com/mdlayher/quota)

Package `quota` provides access to Linux quota netlink notifications.

Quota notifications occur when a user or group's disk quota is exceeded, or when
disk usage falls below a given quota.

For more information on quotas, please see
<https://www.kernel.org/doc/Documentation/filesystems/quota.txt>.

MIT Licensed.
