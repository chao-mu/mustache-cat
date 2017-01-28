package main

import (
	"net/http"
	"regexp"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/unrolled/render"
)

var reUserAgent = regexp.MustCompile(`User-Agent: ([^\r\n]{1,600})`)

type DataHandler struct {
	packets []*Packet
	rend    *render.Render
}

func NewDataHandler(iface string) (*DataHandler, error) {
	dh := &DataHandler{}
	dh.rend = render.New()

	handle, err := pcap.OpenLive(iface, 262144, true, pcap.BlockForever)
	if err != nil {
		return nil, err
	}

	err = handle.SetBPFFilter("tcp and port 80")
	if err != nil {
		return nil, err
	}

	pktSrc := gopacket.NewPacketSource(handle, handle.LinkType())
	go func() {
		for packet := range pktSrc.Packets() {
			dh.packets = append(dh.packets, NewPacket(packet))
		}
	}()

	return dh, nil
}

func (dh *DataHandler) Handler() http.Handler {
	m := http.NewServeMux()

	m.HandleFunc("/data/dump", dh.serveData)

	return m
}

func (dh *DataHandler) serveData(rw http.ResponseWriter, req *http.Request) {
	dh.rend.JSON(rw, http.StatusOK, dh.packets)
}

type Packet struct {
	IP struct {
		Version    uint8
		IHL        uint8
		TOS        uint8
		Length     uint16
		Id         uint16
		Flags      uint8
		FragOffset uint16
		TTL        uint8
		Protocol   uint8
	}
	TCP struct {
		Seq                                        uint32
		Ack                                        uint32
		DataOffset                                 uint8
		FIN, SYN, RST, PSH, ACK, URG, ECE, CWR, NS bool
		Window                                     uint16
		Urgent                                     uint16
	}
	UserAgent string
}

func NewPacket(goPkt gopacket.Packet) *Packet {
	p := &Packet{}

	if tcpLayer := goPkt.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)
		p.TCP.Seq = tcp.Seq
		p.TCP.Ack = tcp.Ack
		p.TCP.Window = tcp.Window
		p.TCP.Urgent = tcp.Urgent
		p.TCP.FIN = tcp.FIN
		p.TCP.SYN = tcp.SYN
		p.TCP.RST = tcp.RST
		p.TCP.PSH = tcp.PSH
		p.TCP.ACK = tcp.ACK
		p.TCP.URG = tcp.URG
		p.TCP.ECE = tcp.ECE
		p.TCP.CWR = tcp.CWR
		p.TCP.NS = tcp.NS
	}

	if ipLayer := goPkt.Layer(layers.LayerTypeIPv4); ipLayer != nil {
		ip, _ := ipLayer.(*layers.IPv4)
		p.IP.Flags = uint8(ip.Flags)
		p.IP.FragOffset = ip.FragOffset
		p.IP.IHL = ip.IHL
		p.IP.Id = ip.Id
		p.IP.Length = ip.Length
		p.IP.Protocol = uint8(ip.Protocol)
		p.IP.TOS = ip.TOS
		p.IP.TTL = ip.TTL
		p.IP.Version = ip.Version
	}

	ua := reUserAgent.FindString(string(goPkt.Data()))
	if ua != "" {
		p.UserAgent = ua
	}

	return p
}
