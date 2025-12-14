package kodama_test

import (
	"cmp"
	"net"
	"slices"
	"testing"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winebarrel/kodama"
)

func TestResolveA_OK(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

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
		q := &dns.Question{
			Name: tt.subdomain,
		}
		rr := handler.ResolveA(q)
		require.NotNil(rr)
		assert.IsType(&dns.A{}, rr)
		assert.Equal(net.ParseIP(tt.ip), rr.(*dns.A).A)
	}
}

func TestResolveA_Nil(t *testing.T) {
	assert := assert.New(t)

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
		q := &dns.Question{
			Name:  tt.subdomain,
			Qtype: dns.TypeA,
		}
		rr := handler.ResolveA(q)
		assert.Nil(rr)
	}
}

func TestResolveA_NS(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	handler := &kodama.Handler{Options: &kodama.Options{
		Domain: "example.com",
		NS:     map[string]string{"ns.example.com": "203.0.113.0"},
	}}

	q := &dns.Question{
		Name:   "ns.example.com.",
		Qtype:  dns.TypeA,
		Qclass: dns.ClassINET,
	}
	rr := handler.ResolveA(q)
	require.NotNil(rr)
	assert.IsType(&dns.A{}, rr)
	assert.Equal(
		&dns.A{
			Hdr: dns.RR_Header{Name: "ns.example.com.", Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 300},
			A:   net.ParseIP("203.0.113.0"),
		},
		rr.(*dns.A),
	)
}

func TestResolveA_NoNS(t *testing.T) {
	assert := assert.New(t)

	handler := &kodama.Handler{Options: &kodama.Options{
		Domain: "example.com",
	}}

	q := &dns.Question{
		Name:   "ns.example.com.",
		Qtype:  dns.TypeA,
		Qclass: dns.ClassINET,
	}
	rr := handler.ResolveA(q)
	assert.Nil(rr)
}

func TestResolveA_NoMatchNS(t *testing.T) {
	assert := assert.New(t)

	handler := &kodama.Handler{Options: &kodama.Options{
		Domain: "example.com",
		NS:     map[string]string{"ns1.example.com": "203.0.113.0"},
	}}

	q := &dns.Question{
		Name:   "ns.example.com.",
		Qtype:  dns.TypeA,
		Qclass: dns.ClassINET,
	}
	rr := handler.ResolveA(q)
	assert.Nil(rr)
}

func TestResolveNS_OK(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	handler := &kodama.Handler{Options: &kodama.Options{
		Domain: "example.com",
		NS: map[string]string{
			"ns1.example.com": "203.0.113.0",
			"ns2.example.com": "203.0.113.1",
		},
	}}

	q := &dns.Question{
		Name: "example.com.",
	}
	rr := handler.ResolveNS(q)

	require.NotNil(rr)
	slices.SortFunc(rr, func(i, j dns.RR) int { return cmp.Compare(i.(*dns.NS).Ns, j.(*dns.NS).Ns) })
	assert.Equal([]dns.RR{
		&dns.NS{Hdr: dns.RR_Header{Name: "example.com.", Ttl: 300}, Ns: "ns1.example.com."},
		&dns.NS{Hdr: dns.RR_Header{Name: "example.com.", Ttl: 300}, Ns: "ns2.example.com."},
	}, rr)
}

func TestResolveNS_Nil(t *testing.T) {
	assert := assert.New(t)

	handler := &kodama.Handler{Options: &kodama.Options{
		Domain: "example.com",
		NS: map[string]string{
			"ns1.example.com": "203.0.113.0",
			"ns2.example.com": "203.0.113.1",
		},
	}}

	q := &dns.Question{
		Name: "example.co.",
	}
	rr := handler.ResolveNS(q)
	assert.Nil(rr)
}

func TestResolveNS_NoConf(t *testing.T) {
	assert := assert.New(t)

	handler := &kodama.Handler{Options: &kodama.Options{
		Domain: "example.com",
	}}

	q := &dns.Question{
		Name: "example.co.",
	}
	rr := handler.ResolveNS(q)
	assert.Nil(rr)
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

	handler := &kodama.Handler{Options: &kodama.Options{
		Domain: "example.com",
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
			Ttl:    0,
		},
		A: net.ParseIP("127.0.0.1"),
	})
}

func TestServeDNS_NoAnswer(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	handler := &kodama.Handler{Options: &kodama.Options{
		Domain: "example.com",
	}}

	w := &MockResponseWriter{}
	q := dns.Question{Name: "127.0.0.1.example.com", Qtype: dns.TypeA, Qclass: dns.ClassINET}
	msg := &dns.Msg{Question: []dns.Question{q}}
	handler.ServeDNS(w, msg)

	require.NotNil(w.m)
	assert.Len(w.m.Answer, 0)
}

func TestServeDNS_CNAME(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	handler := &kodama.Handler{Options: &kodama.Options{
		Domain: "example.com",
	}}

	w := &MockResponseWriter{}
	q := dns.Question{Name: "127-0-0-1.example.com", Qtype: dns.TypeCNAME, Qclass: dns.ClassINET}
	msg := &dns.Msg{Question: []dns.Question{q}}
	handler.ServeDNS(w, msg)

	require.NotNil(w.m)
	assert.Len(w.m.Answer, 0)
}
