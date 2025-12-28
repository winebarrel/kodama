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

func TestZone_Resolve_OK(t *testing.T) {
	data := `
$ORIGIN example.com.     ; designates the start of this zone file in the namespace
$TTL 3600                ; default expiration time (in seconds) of all RRs without their own TTL value
example.com.  IN  SOA   ns.example.com. username.example.com. ( 2020091025 7200 3600 1209600 3600 )
example.com.  IN  NS    ns                    ; ns.example.com is a nameserver for example.com
example.com.  IN  NS    ns.somewhere.example. ; ns.somewhere.example is a backup nameserver for example.com
example.com.  IN  MX    10 mail.example.com.  ; mail.example.com is the mailserver for example.com
@             IN  MX    20 mail2.example.com. ; equivalent to above line, "@" represents zone origin
@             IN  MX    50 mail3              ; equivalent to above line, but using a relative host name
example.com.  IN  A     192.0.2.1             ; IPv4 address for example.com
              IN  AAAA  2001:db8:10::1        ; IPv6 address for example.com
ns            IN  A     192.0.2.2             ; IPv4 address for ns.example.com
              IN  AAAA  2001:db8:10::2        ; IPv6 address for ns.example.com
www           IN  CNAME example.com.          ; www.example.com is an alias for example.com
wwwtest       IN  CNAME www                   ; wwwtest.example.com is another alias for www.example.com
mail          IN  A     192.0.2.3             ; IPv4 address for mail.example.com
mail2         IN  A     192.0.2.4             ; IPv4 address for mail2.example.com
mail3         IN  A     192.0.2.5             ; IPv4 address for mail3.example.com
`
	zone, err := kodama.NewZone(data)
	require.NoError(t, err)

	tests := []struct {
		qname    string
		qtype    uint16
		expected []dns.RR
	}{
		{
			qname: "example.com.",
			qtype: dns.TypeSOA,
			expected: []dns.RR{
				&kodama.RRQ{
					RR: &dns.SOA{
						Hdr:     dns.RR_Header{Name: "example.com.", Rrtype: dns.TypeSOA, Class: dns.ClassINET, Ttl: 3600},
						Ns:      "ns.example.com.",
						Mbox:    "username.example.com.",
						Serial:  2020091025,
						Refresh: 7200,
						Retry:   3600,
						Expire:  1209600,
						Minttl:  3600,
					},
				},
			},
		},
		{
			qname: "example.com.",
			qtype: dns.TypeNS,
			expected: []dns.RR{
				&kodama.RRQ{
					RR: &dns.NS{
						Hdr: dns.RR_Header{Name: "example.com.", Rrtype: dns.TypeNS, Class: dns.ClassINET, Ttl: 3600},
						Ns:  "ns.example.com.",
					},
				},
				&kodama.RRQ{
					RR: &dns.NS{
						Hdr: dns.RR_Header{Name: "example.com.", Rrtype: dns.TypeNS, Class: dns.ClassINET, Ttl: 3600},
						Ns:  "ns.somewhere.example.",
					},
				},
			},
		},
		{
			qname: "example.com.",
			qtype: dns.TypeMX,
			expected: []dns.RR{
				&kodama.RRQ{
					RR: &dns.MX{
						Hdr:        dns.RR_Header{Name: "example.com.", Rrtype: dns.TypeMX, Class: dns.ClassINET, Ttl: 3600},
						Preference: 10,
						Mx:         "mail.example.com.",
					},
				},
				&kodama.RRQ{
					RR: &dns.MX{
						Hdr:        dns.RR_Header{Name: "example.com.", Rrtype: dns.TypeMX, Class: dns.ClassINET, Ttl: 3600},
						Preference: 20,
						Mx:         "mail2.example.com.",
					},
				},
				&kodama.RRQ{
					RR: &dns.MX{
						Hdr:        dns.RR_Header{Name: "example.com.", Rrtype: dns.TypeMX, Class: dns.ClassINET, Ttl: 3600},
						Preference: 50,
						Mx:         "mail3.example.com.",
					},
				},
			},
		},
		{
			qname: "example.com.",
			qtype: dns.TypeA,
			expected: []dns.RR{
				&kodama.RRQ{
					RR: &dns.A{
						Hdr: dns.RR_Header{Name: "example.com.", Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 3600},
						A:   net.ParseIP("192.0.2.1"),
					},
				},
			},
		},
		{
			qname: "example.com.",
			qtype: dns.TypeAAAA,
			expected: []dns.RR{
				&kodama.RRQ{
					RR: &dns.AAAA{
						Hdr:  dns.RR_Header{Name: "example.com.", Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: 3600},
						AAAA: net.ParseIP("2001:db8:10::1"),
					},
				},
			},
		},
		{
			qname: "ns.example.com.",
			qtype: dns.TypeA,
			expected: []dns.RR{
				&kodama.RRQ{
					RR: &dns.A{
						Hdr: dns.RR_Header{Name: "ns.example.com.", Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 3600},
						A:   net.ParseIP("192.0.2.2"),
					},
				},
			},
		},
		{
			qname: "ns.example.com.",
			qtype: dns.TypeAAAA,
			expected: []dns.RR{
				&kodama.RRQ{
					RR: &dns.AAAA{
						Hdr:  dns.RR_Header{Name: "ns.example.com.", Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: 3600},
						AAAA: net.ParseIP("2001:db8:10::2"),
					},
				},
			},
		},
		{
			qname: "www.example.com.",
			qtype: dns.TypeCNAME,
			expected: []dns.RR{
				&kodama.RRQ{
					RR: &dns.CNAME{
						Hdr:    dns.RR_Header{Name: "www.example.com.", Rrtype: dns.TypeCNAME, Class: dns.ClassINET, Ttl: 3600},
						Target: "example.com.",
					},
				},
			},
		},
		{
			qname: "wwwtest.example.com.",
			qtype: dns.TypeCNAME,
			expected: []dns.RR{
				&kodama.RRQ{
					RR: &dns.CNAME{
						Hdr:    dns.RR_Header{Name: "wwwtest.example.com.", Rrtype: dns.TypeCNAME, Class: dns.ClassINET, Ttl: 3600},
						Target: "www.example.com.",
					},
				},
			},
		},
		{
			qname: "mail.example.com.",
			qtype: dns.TypeA,
			expected: []dns.RR{
				&kodama.RRQ{
					RR: &dns.A{
						Hdr: dns.RR_Header{Name: "mail.example.com.", Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 3600},
						A:   net.ParseIP("192.0.2.3"),
					},
				},
			},
		},
		{
			qname: "mail2.example.com.",
			qtype: dns.TypeA,
			expected: []dns.RR{
				&kodama.RRQ{
					RR: &dns.A{
						Hdr: dns.RR_Header{Name: "mail2.example.com.", Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 3600},
						A:   net.ParseIP("192.0.2.4"),
					},
				},
			},
		},
		{
			qname: "mail3.example.com.",
			qtype: dns.TypeA,
			expected: []dns.RR{
				&kodama.RRQ{
					RR: &dns.A{
						Hdr: dns.RR_Header{Name: "mail3.example.com.", Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 3600},
						A:   net.ParseIP("192.0.2.5"),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		for _, name := range []string{tt.qname, strings.ToUpper(tt.qname)} {
			t.Run(name+"/"+dns.TypeToString[tt.qtype], func(t *testing.T) {
				assert := assert.New(t)
				require := require.New(t)

				q := &dns.Question{Name: name, Qtype: tt.qtype}
				rrs := zone.Resolve(q)
				require.NotNil(rrs)

				for _, rr := range tt.expected {
					rrq := rr.(*kodama.RRQ)
					rrq.Qname = name
				}

				assert.Equal(tt.expected, rrs)
			})
		}
	}
}

func TestZone_Resolve_NoAnswer(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	data := `
$ORIGIN example.com.     ; designates the start of this zone file in the namespace
$TTL 3600                ; default expiration time (in seconds) of all RRs without their own TTL value
example.com.  IN  SOA   ns.example.com. username.example.com. ( 2020091025 7200 3600 1209600 3600 )
`
	zone, err := kodama.NewZone(data)
	require.NoError(err)

	q := &dns.Question{Name: "example.com", Qtype: dns.TypeA}
	rrs := zone.Resolve(q)
	assert.Nil(rrs)
}

func TestZone_RRQ_Header(t *testing.T) {
	assert := assert.New(t)

	rrq := &kodama.RRQ{
		RR: &dns.A{
			Hdr: dns.RR_Header{Name: "mail.example.com.", Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 3600},
			A:   net.ParseIP("192.0.2.3"),
		},
		Qname: "MAIL.EXAMPLE.COM.",
	}

	h := rrq.Header()

	assert.Equal(
		&dns.RR_Header{
			Name: "MAIL.EXAMPLE.COM.", Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 3600,
		},
		h,
	)
}
