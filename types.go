package explorerclient

// This files contains an unsorted collection used as sub types of the API resources

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

var (
	dateRe = regexp.MustCompile(`^(?:(\d{2,4})/)?(\d{2})/(\d{2,4})(?:\s+(\d{1,2})(am|pm)?:(\d{1,2}))?$`)
)

// WorkloadAmount represents an amount of workload
type WorkloadAmount struct {
	Network        uint16 `bson:"network" json:"network"`
	Volume         uint16 `bson:"volume" json:"volume"`
	ZDBNamespace   uint16 `bson:"zdb_namespace" json:"zdb_namespace"`
	Container      uint16 `bson:"container" json:"container"`
	K8sVM          uint16 `bson:"k8s_vm" json:"k8s_vm"`
	Proxy          uint16 `bson:"proxy" json:"proxy"`
	ReverseProxy   uint16 `bson:"reverse_proxy" json:"reverse_proxy"`
	Subdomain      uint16 `bson:"subdomain" json:"subdomain"`
	DelegateDomain uint16 `bson:"delegate_domain" json:"delegate_domain"`
}

// HTTPError is the error type returned by the client
// it contains the error and the HTTP response
type HTTPError struct {
	resp *http.Response
	err  error
}

func (h HTTPError) Error() string {
	return fmt.Sprintf("%v status:%s", h.err, h.resp.Status)
}

// Response return the HTTP response that trigger this error
func (h HTTPError) Response() http.Response {
	return *h.resp
}

// WalletAddress represents a wallet address
type WalletAddress struct {
	Asset   string ` json:"asset"`
	Address string `json:"address"`
}

// NodeResourcePrice represents a node resource price
type NodeResourcePrice struct {
	Currency PriceCurrencyEnum `json:"currency"`
	Cru      float64           `json:"cru"`
	Mru      float64           `json:"mru"`
	Hru      float64           `json:"hru"`
	Sru      float64           `json:"sru"`
	Nru      float64           `json:"nru"`
}

// PriceCurrencyEnum represents currencies
type PriceCurrencyEnum uint8

const (
	// PriceCurrencyEUR represents EUR
	PriceCurrencyEUR PriceCurrencyEnum = iota
	// PriceCurrencyUSD Represents USD
	PriceCurrencyUSD
	// PriceCurrencyTFT represents TFT
	PriceCurrencyTFT
	// PriceCurrencyAED represents AED
	PriceCurrencyAED
	// PriceCurrencyGBP represents GBP
	PriceCurrencyGBP
)

func (e PriceCurrencyEnum) String() string {
	switch e {
	case PriceCurrencyEUR:
		return "EUR"
	case PriceCurrencyUSD:
		return "USD"
	case PriceCurrencyTFT:
		return "TFT"
	case PriceCurrencyAED:
		return "AED"
	case PriceCurrencyGBP:
		return "GBP"
	}
	return "UNKNOWN"
}

// MacAddress type
type MacAddress struct{ net.HardwareAddr }

// MarshalText marshals MacAddress type to a string
func (mac MacAddress) MarshalText() ([]byte, error) {
	if mac.HardwareAddr == nil {
		return nil, nil
	} else if mac.HardwareAddr.String() == "" {
		return nil, nil
	}
	return []byte(mac.HardwareAddr.String()), nil
}

// UnmarshalText loads a macaddress from a string
func (mac *MacAddress) UnmarshalText(addr []byte) error {
	if len(addr) == 0 {
		return nil
	}
	addr, err := net.ParseMAC(string(addr))
	if err != nil {
		return err
	}
	mac.HardwareAddr = addr
	return nil
}

// IPRange type
type IPRange struct{ net.IPNet }

// UnmarshalText loads IPRange from string
func (i *IPRange) UnmarshalText(text []byte) error {
	v, err := ParseIPRange(string(text))
	if err != nil {
		return err
	}

	i.IPNet = v.IPNet
	return nil
}

// ParseIPRange parse iprange
func ParseIPRange(txt string) (r IPRange, err error) {
	//empty ip net value
	if len(txt) == 0 {
		return r, nil
	}
	ip, net, err := net.ParseCIDR(txt)
	if err != nil {
		return r, err
	}

	net.IP = ip
	r.IPNet = *net
	return
}

// MarshalJSON dumps IPRange as a string
func (i IPRange) MarshalJSON() ([]byte, error) {
	if len(i.IPNet.IP) == 0 {
		return []byte(`""`), nil
	}
	v := fmt.Sprint("\"", i.String(), "\"")
	return []byte(v), nil
}

func (i IPRange) String() string {
	return i.IPNet.String()
}

// PublicIface represents a public interface
type PublicIface struct {
	Master  string        `json:"master"`
	Type    IfaceTypeEnum `json:"type"`
	Ipv4    IPRange       `json:"ipv4"`
	Ipv6    IPRange       `json:"ipv6"`
	Gw4     net.IP        `json:"gw4"`
	Gw6     net.IP        `json:"gw6"`
	Version int64         `json:"version"`
}

// IfaceTypeEnum represents interface types
type IfaceTypeEnum uint8

const (
	// IfaceTypeMacvlan represents macvlan
	IfaceTypeMacvlan IfaceTypeEnum = iota
	// IfaceTypeVlan represents vlan
	IfaceTypeVlan
)

func (e IfaceTypeEnum) String() string {
	switch e {
	case IfaceTypeMacvlan:
		return "macvlan"
	case IfaceTypeVlan:
		return "vlan"
	}
	return "UNKNOWN"
}

// Iface represents an interface
type Iface struct {
	Name       string     `json:"name"`
	Addrs      []IPRange  `json:"addrs"`
	Gateway    []net.IP   `json:"gateway"`
	MacAddress MacAddress `json:"macaddress"`
}

// Proof represents proof?
type Proof struct {
	Created      Date                   `json:"created"`
	HardwareHash string                 `json:"hardware_hash"`
	DiskHash     string                 `json:"disk_hash"`
	Hardware     map[string]interface{} `json:"hardware"`
	Disks        map[string]interface{} `json:"disks"`
	Hypervisor   []string               `json:"hypervisor"`
}

// ResourceAmount contains an amount for each resource
type ResourceAmount struct {
	Cru uint64  `json:"cru"`
	Mru float64 `json:"mru"`
	Hru float64 `json:"hru"`
	Sru float64 `json:"sru"`
}

// Location represents a physical location
type Location struct {
	City      string  `json:"city"`
	Country   string  `json:"country"`
	Continent string  `json:"continent"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Date is a jumpscale date wrapper
type Date struct{ time.Time }

// UnmarshalJSON converts jumpscale data in json to Golang time.Time
func (d *Date) UnmarshalJSON(bytes []byte) error {
	var inI interface{}
	if err := json.Unmarshal(bytes, &inI); err != nil {
		return err
	}

	var in string
	switch v := inI.(type) {
	case int64:
		d.Time = time.Unix(v, 0).UTC()
		return nil
	case float64:
		d.Time = time.Unix(int64(v), 0).UTC()
		return nil
	case string:
		in = v
	default:
		return fmt.Errorf("unknown date format: %T(%s)", v, string(bytes))
	}

	if len(in) == 0 {
		//null date
		d.Time = time.Time{}
		return nil
	}

	m := dateRe.FindStringSubmatch(in)
	if m == nil {
		return fmt.Errorf("invalid date string '%s'", in)
	}

	first := m[1]
	month := m[2]
	last := m[3]

	hour := m[4]
	ampm := m[5]
	min := m[6]

	var year string
	var day string

	if first == "" {
		year = fmt.Sprint(time.Now().Year())
		day = last
	} else if len(first) == 4 && len(last) == 4 {
		return fmt.Errorf("invalid date format ambiguous year: %s", in)
	} else if len(last) == 4 {
		year = last
		day = first
	} else {
		// both ar 2 or first is 4 and last is 2
		year = first
		day = last
	}

	if hour == "" {
		hour = "0"
	}
	if min == "" {
		min = "0"
	}

	var values []int
	for _, str := range []string{year, month, day, hour, min} {
		value, err := strconv.Atoi(str)
		if err != nil {
			return fmt.Errorf("invalid integer value '%s' in date", str)
		}
		values = append(values, value)
	}

	if values[0] < 100 {
		values[0] += 2000
	}

	if ampm == "pm" {
		values[3] += 12
	}

	d.Time = time.Date(values[0], time.Month(values[1]), values[2], values[3], values[4], 0, 0, time.UTC)

	return nil
}

// MarshalJSON formatting to JSON will return the unix time stamp
func (d Date) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte(`0`), nil
	}
	return []byte(fmt.Sprintf(`%d`, d.Unix())), nil
}

// String implements stringer interface
func (d Date) String() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte(`""`), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, d.Format("02/01/2006 15:04"))), nil
}
