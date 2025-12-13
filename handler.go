package kodama

import (
	"github.com/miekg/dns"
)

type Handler struct{}

func (*Handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := &dns.Msg{}
	msg.SetReply(r)
	msg.Authoritative = true
	w.WriteMsg(msg)
}
