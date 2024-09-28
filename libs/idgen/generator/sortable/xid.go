package sortable

import "github.com/rs/xid"

func NextXid() (string, error) {
	guid := xid.New()
	return guid.String(), nil
}
