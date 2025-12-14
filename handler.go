package kodama

import (
	"net"
	"regexp"
	"strings"

	"github.com/miekg/dns"
)

var (
	rIP = regexp.MustCompile(`\b(\d{1,3})-(\d{1,3})-(\d{1,3})-(\d{1,3})$`)
)

func ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := &dns.Msg{}
	msg.SetReply(r)
	msg.Authoritative = true

	for _, q := range r.Question {
		rr := Resolve(&q)

		if rr != nil {
			msg.Answer = append(msg.Answer, rr)
		}
	}

	w.WriteMsg(msg) //nolint:errcheck
}

func Resolve(q *dns.Question) *dns.A {
	subdomain := strings.Split(q.Name, ".")[0]
	m := rIP.FindStringSubmatch(subdomain)

	if m == nil {
		return nil
	}

	ip := net.ParseIP(strings.Join(m[1:], "."))

	if ip == nil || (!ip.IsLoopback() && !ip.IsPrivate()) {
		return nil
	}

	h := dns.RR_Header{Name: q.Name, Rrtype: q.Qtype, Class: q.Qclass, Ttl: 0}
	return &dns.A{Hdr: h, A: ip}
}
