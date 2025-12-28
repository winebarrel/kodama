package kodama_test

import (
	"net"
	"testing"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winebarrel/kodama"
)

func TestOptions_AfterApply_OK(t *testing.T) {
	assert := assert.New(t)

	options := &kodama.Options{}

	err := options.AfterApply()
	assert.NoError(err)
}

func TestOptions_AfterApply_WithZone_OK(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	zone := `
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

	options := &kodama.Options{
		ZoneData: zone,
	}

	err := options.AfterApply()
	assert.NoError(err)

	q := &dns.Question{Name: "mail3.example.com.", Qtype: dns.TypeA}
	rrs := options.Zone.Resolve(q)
	require.NotNil(rrs)
	require.Len(rrs, 1)
	rr := rrs[0]

	expected := &kodama.RRQ{
		RR: &dns.A{
			Hdr: dns.RR_Header{
				Name:   "mail3.example.com.",
				Rrtype: dns.TypeA,
				Class:  dns.ClassINET,
				Ttl:    3600,
			},
			A: net.ParseIP("192.0.2.5"),
		},
		Qname: "mail3.example.com.",
	}

	assert.Equal(expected, rr)
}

func TestOptions_AfterApply_WithZone_Err(t *testing.T) {
	assert := assert.New(t)

	options := &kodama.Options{
		ZoneData: "invalid",
	}

	err := options.AfterApply()
	assert.ErrorContains(err, "failed to parse zone file data: dns: ")
}
