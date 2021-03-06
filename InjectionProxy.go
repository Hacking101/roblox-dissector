package main

/*
import (
	"context"
	"fmt"
	"net"

	"github.com/Gskartwii/roblox-dissector/peer"
	"github.com/olebedev/emitter"
)

func captureFromInjectionProxy(captureJobContext context.Context, src string, dst string, window *DissectorWindow) {
	srcAddr, err := net.ResolveUDPAddr("udp", src)
	if err != nil {
		return err
	}
	dstAddr, err := net.ResolveUDPAddr("udp", dst)
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", srcAddr)
	if err != nil {
		return err
	}
	defer conn.Close()
	dstConn, err := net.DialUDP("udp", nil, dstAddr)
	if err != nil {
		return err
	}
	defer dstConn.Close()

	// srcAddr = client listen address
	// dstAddr = server connection address

	commContext.Client = srcAddr
	commContext.Server = dstAddr
	proxyWriter := peer.NewProxyWriter(commContext)
	proxyWriter.ServerAddr = dstAddr
	proxyWriter.SecuritySettings = peer.Win10Settings()
	proxyWriter.RuntimeContext, proxyWriter.CancelFunc = context.WithCancel(captureJobContext)

	proxyWriter.ClientHalf.Output.On("udp", func(e *emitter.Event) {
		p := e.Args[0].([]byte)
		_, err := conn.WriteToUDP(p, proxyWriter.ClientAddr)
		if err != nil {
			fmt.Println("write fail to client %s: %s", proxyWriter.ClientAddr.String(), err.Error())
			return
		}
	}, emitter.Void)
	proxyWriter.ServerHalf.Output.On("udp", func(e *emitter.Event) {
		p := e.Args[0].([]byte)
		_, err := dstConn.Write(p)
		if err != nil {
			fmt.Println("write fail to server: %s", err.Error())
			return
		}
	}, emitter.Void)

	var n int
	packetChan := make(chan ProxiedPacket, 100)
	go func() {
		for {
			payload := make([]byte, 1500)
			n, proxyWriter.ClientAddr, err = conn.ReadFromUDP(payload)
			if err != nil {
				fmt.Println("readfromudp fail: %s", err.Error())
				return
			}
			layers := &peer.PacketLayers{
				Root: peer.RootLayer{
					Source:      srcAddr,
					Destination: dstAddr,
					FromClient:  true,
				},
			}
			if payload[0] > 0x8 {
				packetChan <- ProxiedPacket{Layers: layers, Payload: payload[:n]}
			} else { // Need priority for join packets
				proxyWriter.ProxyClient(payload[:n], layers)
			}
		}
	}()
	go func() {
		for {
			payload := make([]byte, 1500)
			n, addr, err := dstConn.ReadFromUDP(payload)
			if err != nil {
				fmt.Println("readfromudp fail %s: %s", addr.String(), err.Error())
				return
			}
			layers := &peer.PacketLayers{
				Root: peer.RootLayer{
					Source:      dstAddr,
					Destination: srcAddr,
					FromServer:  true,
				},
			}
			if payload[0] > 0x8 {
				packetChan <- ProxiedPacket{Layers: layers, Payload: payload[:n]}
			} else { // Need priority for join packets
				proxyWriter.ProxyServer(payload[:n], layers)
			}
		}
	}()
	for {
		select {
		case newPacket := <-packetChan:
			if newPacket.Layers.Root.FromClient {
				proxyWriter.ProxyClient(newPacket.Payload, newPacket.Layers)
			} else {
				proxyWriter.ProxyServer(newPacket.Payload, newPacket.Layers)
			}
		case _ = <-injectPacket:
			//proxyWriter.InjectServer(injectedPacket)
			println("Attempt to inject packet not implemented")
		case <-captureJobContext.Done():
			proxyWriter.CancelFunc()
			return
		}
	}
	return
}
*/
