package kodama

import (
	"strings"

	"github.com/miekg/dns"
)

type Zone struct {
	m map[string]map[uint16][]dns.RR
}

type RRQ struct {
	dns.RR
	Qname string
}

func (rrq *RRQ) Header() *dns.RR_Header {
	h := rrq.RR.Header()
	h.Name = rrq.Qname
	return h
}

func NewZone(src string) (*Zone, error) {
	parser := dns.NewZoneParser(strings.NewReader(src), "", "")
	rrs := []dns.RR{}

	for rr, ok := parser.Next(); ok; rr, ok = parser.Next() {
		rrs = append(rrs, rr)
	}

	if err := parser.Err(); err != nil {
		return nil, err
	}

	m := map[string]map[uint16][]dns.RR{}

	for _, rr := range rrs {
		h := rr.Header()
		rrname := h.Name
		rrtype := h.Rrtype
		rrsByType, ok := m[rrname]

		if !ok {
			rrsByType = map[uint16][]dns.RR{}
			m[rrname] = rrsByType
		}

		rrSet, ok := rrsByType[rrtype]

		if !ok {
			rrSet = []dns.RR{}
		}

		rrSet = append(rrSet, rr)
		rrsByType[rrtype] = rrSet
	}

	z := &Zone{
		m: m,
	}

	return z, nil
}

func (z *Zone) Resolve(q *dns.Question) []dns.RR {
	if z.m == nil {
		return nil
	}

	qname := dns.CanonicalName(q.Name)
	qtype := q.Qtype
	rrsByType, ok := z.m[qname]

	if !ok {
		return nil
	}

	rrs, ok := rrsByType[qtype]

	if !ok {
		return nil
	}

	newRrs := []dns.RR{}

	for _, rr := range rrs {
		rrq := &RRQ{RR: rr, Qname: q.Name}
		newRrs = append(newRrs, rrq)
	}

	return newRrs
}
