package kodama

import (
	"fmt"
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

		var rrs []dns.RR
		switch q.Qtype {
		case dns.TypeNS:
			rrs = h.ResolveNS(&q)
		case dns.TypeA:
			if rr := h.ResolveA(&q); rr != nil {
				fmt.Println(rr)
				rrs = []dns.RR{rr}
			}
		}

		log.Printf("< %s", rrs)

		if len(rrs) >= 1 {
			msg.Answer = append(msg.Answer, rrs...)
		}
	}

	w.WriteMsg(msg) //nolint:errcheck
}

func (h *Handler) ResolveA(q *dns.Question) dns.RR {
	for name, ipstr := range h.NS {
		name = dns.CanonicalName(name)

		if name != q.Name {
			continue
		}

		ip := net.ParseIP(ipstr)
		rrh := dns.RR_Header{Name: q.Name, Rrtype: q.Qtype, Class: q.Qclass, Ttl: 300}
		return &dns.A{Hdr: rrh, A: ip}
	}

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

func (h *Handler) ResolveNS(q *dns.Question) []dns.RR {
	if len(h.NS) == 0 || dns.CanonicalName(h.Domain) != q.Name {
		return nil
	}

	rrs := []dns.RR{}

	for name := range h.NS {
		rrh := dns.RR_Header{Name: q.Name, Rrtype: q.Qtype, Class: q.Qclass, Ttl: 300}
		rrs = append(rrs, &dns.NS{Hdr: rrh, Ns: dns.CanonicalName(name)})
	}

	return rrs
}
