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
		var rrSet []dns.RR

		if rrs := h.Zone.Resolve(&q); rrs != nil {
			rrSet = rrs
		} else {
			if q.Qtype == dns.TypeA {
				if rrs := h.ResolveDynamic(&q); rrs != nil {
					rrSet = rrs
				}
			}
		}

		if len(rrSet) >= 1 {
			slog.Info("OK", "question", q.String(), "answer", rrSet)
			msg.Answer = append(msg.Answer, rrSet...)
		} else {
			slog.Debug("NOT FOUND", "question", q.String())
		}
	}

	w.WriteMsg(msg) //nolint:errcheck
}

func (h *Handler) ResolveDynamic(q *dns.Question) []dns.RR {
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
	return []dns.RR{&dns.A{Hdr: rrh, A: ip}}
}
