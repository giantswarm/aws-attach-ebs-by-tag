package routing

import (
	"bytes"
	"io/ioutil"
	"net"
	"text/template"

	"github.com/giantswarm/microerror"
)

const (
	eth1FileName = "/etc/systemd/network/10-eth1.network"
)

type params struct {
	ENIAddress    string
	ENIGateway    string
	ENISubnetSize int
}

func ConfigureNetworkRoutingForENI(eniIP string, eniSubnet *net.IPNet) error {

	p := params{
		ENIAddress:    eniIP,
		ENIGateway:    eniGateway(eniSubnet),
		ENISubnetSize: eniSubnetSize(eniSubnet),
	}

	err := renderRoutingNetworkdFile(p)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func renderRoutingNetworkdFile(p params) error {
	var buff bytes.Buffer
	t := template.Must(template.New("routing").Parse(networkRoutingTemplate))

	err := t.Execute(&buff, p)
	if err != nil {
		return microerror.Mask(err)
	}

	err = ioutil.WriteFile(eth1FileName, buff.Bytes(), 0644)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func eniGateway(ipNet *net.IPNet) string {
	// https://docs.aws.amazon.com/vpc/latest/userguide/VPC_Subnets.html
	gatewayAddressIP := dupIP(ipNet.IP)
	gatewayAddressIP.To4()
	gatewayAddressIP[3] += 1

	return gatewayAddressIP.String()
}

func eniSubnetSize(ipNet *net.IPNet) int {
	subnetSize, _ := ipNet.Mask.Size()

	return subnetSize
}

func dupIP(ip net.IP) net.IP {
	dup := make(net.IP, len(ip))
	copy(dup, ip)
	return dup
}
