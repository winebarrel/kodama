package kodama_test

import (
	"testing"
	"time"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winebarrel/kodama"
)

func TestServerStart_TCP(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	options := &kodama.Options{Domain: "example.com", Addr: "127.0.0.1:0"}
	svr := kodama.NewServer(options)
	go svr.Start(t.Context())
	dnsSvr := svr.TCP

	for dnsSvr.Listener == nil {
		time.Sleep(100 * time.Millisecond)
	}

	msg := &dns.Msg{}
	msg.SetQuestion(dns.Fqdn("127-0-0-1.example.com"), dns.TypeA)
	client := &dns.Client{Net: "tcp", Timeout: 5 * time.Second}
	resp, _, err := client.Exchange(msg, dnsSvr.Listener.Addr().String())

	require.NoError(err)
	require.NotNil(resp)
	require.Len(resp.Answer, 1)

	answer := resp.Answer[0]
	assert.Equal("127-0-0-1.example.com.	0	IN	A	127.0.0.1", answer.String())
}

func TestServerStart_UDP(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	options := &kodama.Options{Domain: "example.com", Addr: "127.0.0.1:0"}
	svr := kodama.NewServer(options)
	go svr.Start(t.Context())
	dnsSvr := svr.UDP

	for dnsSvr.PacketConn == nil {
		time.Sleep(100 * time.Millisecond)
	}

	msg := &dns.Msg{}
	msg.SetQuestion(dns.Fqdn("127-0-0-1.example.com"), dns.TypeA)
	client := &dns.Client{Net: "udp", Timeout: 5 * time.Second}
	resp, _, err := client.Exchange(msg, dnsSvr.PacketConn.LocalAddr().String())

	require.NoError(err)
	require.NotNil(resp)
	require.Len(resp.Answer, 1)

	answer := resp.Answer[0]
	assert.Equal("127-0-0-1.example.com.	0	IN	A	127.0.0.1", answer.String())
}
