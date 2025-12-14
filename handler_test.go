package kodama_test

import (
	"net"
	"testing"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winebarrel/kodama"
)

func TestResolve_OK(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	tests := []struct {
		subdomain string
		ip        string
	}{
		{
			subdomain: "127-0-0-1.example.com",
			ip:        "127.0.0.1",
		},
		{
			subdomain: "prefix-127-0-0-1.example.com",
			ip:        "127.0.0.1",
		},
		{
			subdomain: "10-0-0-0.example.com",
			ip:        "10.0.0.0",
		},
		{
			subdomain: "10-255-255-255.example.com",
			ip:        "10.255.255.255",
		},
		{
			subdomain: "172-16-0-0.example.com",
			ip:        "172.16.0.0",
		},
		{
			subdomain: "172-31-255-255.example.com",
			ip:        "172.31.255.255",
		},
		{
			subdomain: "192-168-0-0.example.com",
			ip:        "192.168.0.0",
		},
		{
			subdomain: "192-168-255-255.example.com",
			ip:        "192.168.255.255",
		},
	}

	for _, test := range tests {
		q := &dns.Question{
			Name: test.subdomain,
		}
		rr := kodama.Resolve(q)
		require.NotNil(rr)
		assert.IsType(&dns.A{}, rr)
		assert.Equal(net.ParseIP(test.ip), rr.A)
	}
}

func TestResolve_Nil(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		subdomain string
	}{
		{
			subdomain: "example.com",
		},
		{
			subdomain: "www.example.com",
		},
		{
			subdomain: "128-0-0-1.example.com",
		},
		{
			subdomain: "prefix-128-0-0-1.example.com",
		},
		{
			subdomain: "9-255-255-255.example.com",
		},
		{
			subdomain: "11-0-0-0.example.com",
		},
		{
			subdomain: "172-15-255-255.example.com",
		},
		{
			subdomain: "172-32-0-0.example.com",
		},
		{
			subdomain: "192-167-255-255.example.com",
		},
		{
			subdomain: "192-169-0-0.example.com",
		},
	}

	for _, test := range tests {
		q := &dns.Question{
			Name: test.subdomain,
		}
		rr := kodama.Resolve(q)
		assert.Nil(rr)
	}
}

type MockResponseWriter struct {
	m *dns.Msg
}

func (*MockResponseWriter) LocalAddr() net.Addr              { return nil }
func (*MockResponseWriter) RemoteAddr() net.Addr             { return nil }
func (mock *MockResponseWriter) WriteMsg(msg *dns.Msg) error { mock.m = msg; return nil }
func (*MockResponseWriter) Write([]byte) (int, error)        { return 0, nil }
func (*MockResponseWriter) Close() error                     { return nil }
func (*MockResponseWriter) TsigStatus() error                { return nil }
func (*MockResponseWriter) TsigTimersOnly(bool)              {}
func (*MockResponseWriter) Hijack()                          {}

var _ dns.ResponseWriter = &MockResponseWriter{}

func TestServeDNS(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	w := &MockResponseWriter{}
	msg := &dns.Msg{Question: []dns.Question{{Name: "127-0-0-1.example.com"}}}
	kodama.ServeDNS(w, msg)

	require.NotNil(w.m)
	require.Len(w.m.Answer, 1)

	answer := w.m.Answer[0]
	assert.Equal(answer, &dns.A{
		Hdr: dns.RR_Header{
			Name: "127-0-0-1.example.com",
		},
		A: net.ParseIP("127.0.0.1"),
	})
}
