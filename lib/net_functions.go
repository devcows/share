package lib

// Ask the kernel for a free open port that is ready to use
import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/prestonTao/upnp"
)

func RandomFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return -1, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return -1, err
	}
	defer l.Close()

	return l.Addr().(*net.TCPAddr).Port, nil
}

func OpenUpnpPort(port int) bool {
	mapping := new(upnp.Upnp)
	err := mapping.AddPortMapping(port, port, "TCP")
	if err == nil {
		log.Info("Upnp port opened!")
		return true
		// remove port mapping in gatway
		// mapping.Reclaim()
	}

	log.Info(fmt.Sprintf("Fail open upnp port opened. Error: %s", err.Error()))
	return false
}

func GetLocalIps(port int) []string {
	listIps := []string{}

	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			listIps = append(listIps, fmt.Sprintf("http://%s:%v", ipv4, port))
		}
	}

	return listIps
}

func GetPublicIps(port int) []string {
	listIps := []string{}
	resp, err := http.Get("http://myexternalip.com/raw")
	if err == nil {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if err == nil {
			lines := strings.Split(string(body), "\n")
			listIps = append(listIps, fmt.Sprintf("http://%s:%v", lines[0], port))
		}
	}

	return listIps
}
