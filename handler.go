package kodama

import (
	"net"

	"github.com/miekg/dns"
)

type Handler struct{}

func (h *Handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := &dns.Msg{}
	msg.SetReply(r)
	msg.Authoritative = true

	for _, q := range r.Question {
		rr := h.Resolve(&q)

		if rr != nil {
			msg.Answer = append(msg.Answer, rr)
		}
	}

	w.WriteMsg(msg)
}

// TODO:
func (*Handler) Resolve(q *dns.Question) *dns.A {
	ip := net.ParseIP("127.0.0.1")

	if ip == nil {
		return nil
	}

	h := dns.RR_Header{Name: q.Name, Rrtype: q.Qtype, Class: q.Qclass, Ttl: 0}
	return &dns.A{Hdr: h, A: ip}
}
