package dao

import "testing"

func TestAddAccess(t *testing.T) {
	AddAccess(1, "/", "127.0.0.1", "", "")
}
