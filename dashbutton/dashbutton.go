package dashbutton

import (
	"log"
	"strings"

	"github.com/krolaw/dhcp4"
)

const DashMacAddress = "44:65:0d:4a:e2:b4"

type DhcpHandler struct {
	interrupt func()
}

func New(interrupt func()) {

	handler := &DhcpHandler{interrupt: interrupt}

	err := dhcp4.ListenAndServe(handler)
	if err != nil {
		log.Println(err)
	}

}

func (mh *DhcpHandler) ServeDHCP(req dhcp4.Packet, msgType dhcp4.MessageType, options dhcp4.Options) dhcp4.Packet {

	addr := req.CHAddr()

	if addr != nil && DashMacAddress == strings.ToLower(addr.String()) {
		mh.interrupt()
	}

	return nil
}
