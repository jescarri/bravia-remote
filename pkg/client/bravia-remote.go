package bravia

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"io/ioutil"
	"net/http"
	"strings"
)

type Remote struct {
	IpAddress string
	Secret    string
	IRCCUrl   string
}

const GetIrcCodes string = `{"method":"getRemoteControllerInfo","params":[],"id":10,"version":"1.0"}`

func (r *Remote) NewRemote(ipAddress string, secret string) *Remote {
	r.IpAddress = ipAddress
	r.Secret = secret
	r.IRCCUrl = fmt.Sprintf("http://%s/sony/IRCC", r.IpAddress)
	return r
}

func (r *Remote) SendCode(code string) error {
	cmd := fmt.Sprintf("<?xml version=\"1.0\"?><s:Envelope xmlns:s=\"http://schemas.xmlsoap.org/soap/envelope/\" s:encodingStyle=\"http://schemas.xmlsoap.org/soap/encoding/\"><s:Body><u:X_SendIRCC xmlns:u=\"urn:schemas-sony-com:service:IRCC:1\"><IRCCCode>%s</IRCCCode></u:X_SendIRCC></s:Body></s:Envelope>", code)
	body := strings.NewReader(cmd)
	fmt.Println(cmd)
	fmt.Println(r.IRCCUrl)
	req, err := http.NewRequest("POST", r.IRCCUrl, body)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "\"text/xml; charset=UTF-8\"")
	req.Header.Add("SOAPACTION", "\"urn:schemas-sony-com:service:IRCC:1#X_SendIRCC\"")
	req.Header.Add("X-Auth-PSK", r.Secret)
	spew.Dump(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	ret, err := ioutil.ReadAll(resp.Body)
	spew.Dump(ret)
	if resp.StatusCode != 200 {
		return fmt.Errorf("TV Returned error %d", resp.StatusCode)
	}
	return nil

}
