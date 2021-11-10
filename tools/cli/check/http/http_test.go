package http

import (
	"context"
	"github.com/stretchr/testify/require"
	"net"
	"testing"
)

func TestCheckHTTP(t *testing.T) {
	addr, err := net.LookupIP("img.joomcdn.net")
	if err != nil {
		require.NoError(t, err)
	}
	require.True(t, len(addr) > 0)
	CheckHTTP(context.Background(), "https://img.joomcdn.net/143277b4305cfcb23573b35ba9d26448e71d8eb4_100_100.jpeg", addr[0])
	CheckHTTP(context.Background(), "https://speedtest.bozaro.ru/", nil)

	CheckHTTP(context.Background(), "https://speedtest.bozaro.ru/", addr[0])

	addr, _ = net.LookupIP("172.16.1.2")
	CheckHTTP(context.Background(), "https://speedtest1.bozaro.ru/", addr[0])
}
