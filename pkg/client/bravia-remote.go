package bravia

import (
	"encoding/json"
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

type Remote struct {
	Host          string
	Secret        string
	IRCCUrl       string
	SystemUrl     string
	AvailabeCodes map[string]string
}

type Result struct {
	Codes []interface{} `json:"result"`
}

const GetIrcCodes string = `{"method":"getRemoteControllerInfo","params":[],"id":10,"version":"1.0"}`
const GetTVStatus string = `{"method":"getPowerStatus","params":[],"id":1,"version":"1.0"}`

func (r *Remote) NewRemote(host string, secret string) (*Remote, error) {
	r.Host = host
	r.Secret = secret
	r.IRCCUrl = fmt.Sprintf("http://%s/sony/IRCC", r.Host)
	r.SystemUrl = fmt.Sprintf("http://%s/sony/system", r.Host)
	c, err := r.GetSupportedActions()
	if err != nil {
		return &Remote{}, err
	}
	r.AvailabeCodes = c
	return r, nil
}

func (r *Remote) GetSupportedActions() (map[string]string, error) {
	body := strings.NewReader(GetIrcCodes)
	req, err := http.NewRequest("POST", r.SystemUrl, body)
	if err != nil {
		return map[string]string{}, err
	}
	req.Header.Add("X-Auth-PSK", r.Secret)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return map[string]string{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return map[string]string{}, fmt.Errorf("TV Returned error %d", resp.StatusCode)
	}
	supportedActions, err := flatenCodesResponse(resp.Body)
	if err != nil {
		return map[string]string{}, err
	}
	return supportedActions, nil

}

func flatenCodesResponse(r io.Reader) (map[string]string, error) {
	a := Result{}

	codes := map[string]string{}
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return map[string]string{}, err
	}
	err = json.Unmarshal(body, &a)
	if err != nil {
		return map[string]string{}, err
	}
	for _, v := range a.Codes {
		switch reflect.TypeOf(v).Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(v)
			for i := 0; i < s.Len(); i++ {
				if reflect.TypeOf(s.Index(i)).Kind() == reflect.Struct {
					ds := s.Index(i).Interface().(map[string]interface{})
					if name, ok := ds["name"]; ok {
						codes[name.(string)] = ds["value"].(string)
					}

				}
			}
		}
	}
	return codes, nil
}

func (r *Remote) SendCode(code string) error {
	cmd := fmt.Sprintf("<?xml version=\"1.0\"?><s:Envelope xmlns:s=\"http://schemas.xmlsoap.org/soap/envelope/\" s:encodingStyle=\"http://schemas.xmlsoap.org/soap/encoding/\"><s:Body><u:X_SendIRCC xmlns:u=\"urn:schemas-sony-com:service:IRCC:1\"><IRCCCode>%s</IRCCCode></u:X_SendIRCC></s:Body></s:Envelope>", code)
	body := strings.NewReader(cmd)
	req, err := http.NewRequest("POST", r.IRCCUrl, body)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "\"text/xml; charset=UTF-8\"")
	req.Header.Add("SOAPACTION", "\"urn:schemas-sony-com:service:IRCC:1#X_SendIRCC\"")
	req.Header.Add("X-Auth-PSK", r.Secret)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("TV Returned error %d", resp.StatusCode)
	}
	return nil

}

func (r *Remote) GetCodeForAction(action string) (string, error) {
	if code, ok := r.AvailabeCodes[action]; ok {
		return code, nil
	}
	return "", fmt.Errorf("No code found for Action %s", action)
}

func (r *Remote) Do(action string) error {
	code, err := r.GetCodeForAction(action)
	if err != nil {
		return err
	}
	err = r.SendCode(code)
	if err != nil {
		return err
	}
	return nil
}
