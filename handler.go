package kodama

import (
	"log"
	"net"
	"regexp"
	"strings"

	"github.com/miekg/dns"
)

var (
	rIP = regexp.MustCompile(`\b(\d{1,3})-(\d{1,3})-(\d{1,3})-(\d{1,3})$`)
)

type Handler struct {
	*Options
}

func (h *Handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := &dns.Msg{}
	msg.SetReply(r)
	msg.Authoritative = true

	for _, q := range r.Question {
		log.Println(">", q.String())

		var rr dns.RR
		switch q.Qtype {
		case dns.TypeNS:
			if h.NS != "" && dns.CanonicalName(h.Domain) == q.Name {
				rr = &dns.NS{
					Hdr: dns.RR_Header{Name: q.Name, Rrtype: q.Qtype, Class: q.Qclass, Ttl: 300},
					Ns:  dns.CanonicalName(h.NS),
				}
			}
		case dns.TypeA:
			rr = h.Resolve(&q)
		}

		if rr != nil {
			log.Println("<", rr.String())
			msg.Answer = append(msg.Answer, rr)
		} else {
			log.Println("<", "(no response)")
		}
	}

	w.WriteMsg(msg) //nolint:errcheck
}

func (h *Handler) Resolve(q *dns.Question) dns.RR {
	subdomain := strings.Split(q.Name, ".")[0]
	m := rIP.FindStringSubmatch(subdomain)

	if m == nil {
		return nil
	}

	ip := net.ParseIP(strings.Join(m[1:], "."))

	if ip == nil || (!ip.IsLoopback() && !ip.IsPrivate()) {
		return nil
	}

	rrh := dns.RR_Header{Name: q.Name, Rrtype: q.Qtype, Class: q.Qclass, Ttl: 0}
	return &dns.A{Hdr: rrh, A: ip}
}
