package net

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResolveHostIP(t *testing.T) {
	_, err := ResolveHostIP("interface-not-found")
	require.Equal(t, ErrNetworkInterfaceNotFound, err)

	if os.Getenv("SKIP_TestResolveHostIP_IPLOOKUP") != "" {
		t.Skip("skip ip address lookup")
	}
	_, err = ResolveHostIP("lo0")
	require.Equal(t, ErrNoAvailableIPAddress, err)
	ip, err := ResolveHostIP("")
	require.Nil(t, err)
	require.NotNil(t, ip)
	t.Logf("HostIP: %s", ip)
}
