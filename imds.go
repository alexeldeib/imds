// This package provides a few simples types around the data returned by Azure IMDS.
// The client hard codes the API version to 2019-03-11, so the types will correspond to that format.
package imds

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

func New() (*Metadata, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://169.254.169.254/metadata/instance", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Metadata", "True")

	q := req.URL.Query()
	q.Add("format", "json")
	q.Add("api-version", "2019-08-15")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	data := &Metadata{}
	err = json.Unmarshal(body, data)
	data.ParsedTags = parseTags(data.Compute.Tags)
	return data, err
}

func parseTags(input string) map[string]string {
	tags := map[string]string{}
	tagString := strings.Split(input, ";")
	for _, tag := range tagString {
		tuple := strings.Split(tag, ":")
		tags[tuple[0]] = tuple[1]
	}
	return tags
}

type Metadata struct {
	Compute    Compute           `json:"compute"`
	Network    Network           `json:"network"`
	ParsedTags map[string]string `json:"parsedTags"`
}

type Plan struct {
	Name      string `json:"name"`
	Product   string `json:"product"`
	Publisher string `json:"publisher"`
}

type PublicKeys struct {
	KeyData string `json:"keyData"`
	Path    string `json:"path"`
}

type Image struct {
	URI string `json:"uri"`
}

type ManagedDisk struct {
	ID                 string `json:"id"`
	StorageAccountType string `json:"storageAccountType"`
}

type VHD struct {
	URI string `json:"uri"`
}

type DataDisks struct {
	Caching                 string      `json:"caching"`
	CreateOption            string      `json:"createOption"`
	DiskSizeGB              string      `json:"diskSizeGB"`
	Image                   Image       `json:"image"`
	LUN                     string      `json:"lun"`
	ManagedDisk             ManagedDisk `json:"managedDisk"`
	Name                    string      `json:"name"`
	VHD                     VHD         `json:"vhd"`
	WriteAcceleratorEnabled string      `json:"writeAcceleratorEnabled"`
}

type ImageReference struct {
	ID        string `json:"id"`
	Offer     string `json:"offer"`
	Publisher string `json:"publisher"`
	Sku       string `json:"sku"`
	Version   string `json:"version"`
}

type DiffDiskSettings struct {
	Option string `json:"option"`
}

type EncryptionSettings struct {
	Enabled string `json:"enabled"`
}

type OsDisk struct {
	Caching                 string             `json:"caching"`
	CreateOption            string             `json:"createOption"`
	DiffDiskSettings        DiffDiskSettings   `json:"diffDiskSettings"`
	DiskSizeGB              string             `json:"diskSizeGB"`
	EncryptionSettings      EncryptionSettings `json:"encryptionSettings"`
	Image                   Image              `json:"image"`
	ManagedDisk             ManagedDisk        `json:"managedDisk"`
	Name                    string             `json:"name"`
	OsType                  string             `json:"osType"`
	VHD                     VHD                `json:"vhd"`
	WriteAcceleratorEnabled string             `json:"writeAcceleratorEnabled"`
}

type StorageProfile struct {
	DataDisks      []DataDisks    `json:"dataDisks"`
	ImageReference ImageReference `json:"imageReference"`
	OsDisk         OsDisk         `json:"osDisk"`
}

type Compute struct {
	AzEnvironment        string         `json:"azEnvironment"`
	CustomData           string         `json:"customData"`
	Location             string         `json:"location"`
	Name                 string         `json:"name"`
	Offer                string         `json:"offer"`
	OsType               string         `json:"osType"`
	PlacementGroupID     string         `json:"placementGroupId"`
	Plan                 Plan           `json:"plan"`
	PlatformFaultDomain  string         `json:"platformFaultDomain"`
	PlatformUpdateDomain string         `json:"platformUpdateDomain"`
	Provider             string         `json:"provider"`
	PublicKeys           []PublicKeys   `json:"publicKeys"`
	Publisher            string         `json:"publisher"`
	ResourceGroupName    string         `json:"resourceGroupName"`
	ResourceID           string         `json:"resourceId"`
	Sku                  string         `json:"sku"`
	StorageProfile       StorageProfile `json:"storageProfile"`
	SubscriptionID       string         `json:"subscriptionId"`
	Tags                 string         `json:"tags"`
	TagsList             []interface{}  `json:"tagsList"`
	Version              string         `json:"version"`
	VMID                 string         `json:"vmId"`
	VMScaleSetName       string         `json:"vmScaleSetName"`
	VMSize               string         `json:"vmSize"`
	Zone                 string         `json:"zone"`
}

type IPAddress struct {
	PrivateIPAddress string `json:"privateIpAddress"`
	PublicIPAddress  string `json:"publicIpAddress"`
}

type Subnet struct {
	Address string `json:"address"`
	Prefix  string `json:"prefix"`
}

type IPv4 struct {
	IPAddress []IPAddress `json:"ipAddress"`
	Subnet    []Subnet    `json:"subnet"`
}

type IPv6 struct {
	IPAddress []interface{} `json:"ipAddress"`
}

type Interface struct {
	IPv4       IPv4   `json:"ipv4"`
	IPv6       IPv6   `json:"ipv6"`
	MacAddress string `json:"macAddress"`
}

type Network struct {
	Interface []Interface `json:"interface"`
}
