package kodama_test

import (
	"net"
	"strings"
	"testing"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winebarrel/kodama"
)

func TestResolveDynamic_OK(t *testing.T) {
	handler := &kodama.Handler{Options: &kodama.Options{
		Domain: "example.com",
	}}

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

	for _, tt := range tests {
		for _, name := range []string{tt.subdomain, strings.ToUpper(tt.subdomain)} {
			t.Run(name+"->"+tt.ip, func(t *testing.T) {
				assert := assert.New(t)
				require := require.New(t)

				q := &dns.Question{Name: name, Qtype: dns.TypeA}
				rrs := handler.ResolveDynamic(q)

				require.NotNil(rrs)
				require.Len(rrs, 1)
				rr := rrs[0]
				assert.IsType(&dns.A{}, rr)
				assert.Equal(net.ParseIP(tt.ip), rr.(*dns.A).A)
			})
		}
	}
}

func TestTestResolveDynamic_Nil(t *testing.T) {
	handler := &kodama.Handler{Options: &kodama.Options{
		Domain: "example.com",
	}}

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

	for _, tt := range tests {
		t.Run(tt.subdomain+"->nil", func(t *testing.T) {
			assert := assert.New(t)

			q := &dns.Question{Name: tt.subdomain, Qtype: dns.TypeA}
			rr := handler.ResolveDynamic(q)

			assert.Nil(rr)
		})
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

func TestServeDNS_Dynamic(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	handler := &kodama.Handler{Options: &kodama.Options{
		Domain: "example.com",
		TTL:    600,
		Zone:   &kodama.Zone{},
	}}

	w := &MockResponseWriter{}
	q := dns.Question{Name: "127-0-0-1.example.com", Qtype: dns.TypeA, Qclass: dns.ClassINET}
	msg := &dns.Msg{Question: []dns.Question{q}}
	handler.ServeDNS(w, msg)

	require.NotNil(w.m)
	require.Len(w.m.Answer, 1)

	answer := w.m.Answer[0]
	assert.Equal(answer, &dns.A{
		Hdr: dns.RR_Header{
			Name:   "127-0-0-1.example.com",
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
			Ttl:    600,
		},
		A: net.ParseIP("127.0.0.1"),
	})
}
