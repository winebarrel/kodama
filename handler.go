package kodama

import (
	"log/slog"
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
		var rrs []dns.RR
		switch q.Qtype {
		case dns.TypeNS:
			rrs = h.ResolveNS(&q)
		case dns.TypeA:
			if rr := h.ResolveA(&q); rr != nil {
				rrs = []dns.RR{rr}
			}
		case dns.TypeTXT:
			if rr := h.ResolveTXT(&q); rr != nil {
				rrs = []dns.RR{rr}
			}
		case dns.TypeCNAME:
			if rr := h.ResolveCNAME(&q); rr != nil {
				rrs = []dns.RR{rr}
			}
		}

		if len(rrs) >= 1 {
			slog.Info("OK", "question", q.String(), "answer", rrs)
			msg.Answer = append(msg.Answer, rrs...)
		} else {
			slog.Debug("NOT FOUND", "question", q.String(), "answer", rrs)
		}
	}

	w.WriteMsg(msg) //nolint:errcheck
}

func (h *Handler) ResolveA(q *dns.Question) dns.RR {
	qname := dns.CanonicalName(q.Name)

	for name, ipstr := range h.NS {
		name = dns.CanonicalName(name)

		if name != qname {
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

	rrh := dns.RR_Header{Name: q.Name, Rrtype: q.Qtype, Class: q.Qclass, Ttl: h.TTL}
	return &dns.A{Hdr: rrh, A: ip}
}

func (h *Handler) ResolveNS(q *dns.Question) []dns.RR {
	qname := dns.CanonicalName(q.Name)

	if len(h.NS) == 0 || dns.CanonicalName(h.Domain) != qname {
		return nil
	}

	rrs := []dns.RR{}

	for name := range h.NS {
		rrh := dns.RR_Header{Name: q.Name, Rrtype: q.Qtype, Class: q.Qclass, Ttl: 300}
		rrs = append(rrs, &dns.NS{Hdr: rrh, Ns: dns.CanonicalName(name)})
	}

	return rrs
}

func (h *Handler) ResolveTXT(q *dns.Question) dns.RR {
	qname := dns.CanonicalName(q.Name)
	rrh := dns.RR_Header{Name: q.Name, Rrtype: q.Qtype, Class: q.Qclass, Ttl: 300}

	if dns.CanonicalName(h.Domain) == qname {
		return &dns.TXT{Hdr: rrh, Txt: []string{"kodama-version=" + h.Version}}
	}

	for name, value := range h.TXT {
		name := dns.CanonicalName(name)

		if name == qname {
			return &dns.TXT{Hdr: rrh, Txt: []string{value}}
		}
	}

	return nil
}

func (h *Handler) ResolveCNAME(q *dns.Question) dns.RR {
	qname := dns.CanonicalName(q.Name)
	rrh := dns.RR_Header{Name: q.Name, Rrtype: q.Qtype, Class: q.Qclass, Ttl: 300}

	for name, value := range h.CNAME {
		name := dns.CanonicalName(name)

		if name == qname {
			return &dns.CNAME{Hdr: rrh, Target: dns.CanonicalName(value)}
		}
	}

	return nil
}
