package megaport

// MVEOrderConfig represents a request to buy an MVE from the Megaport Products API.
type MVEOrderConfig struct {
	LocationID        int                    `json:"locationId"`
	Name              string                 `json:"productName"`
	Term              int                    `json:"term"`
	ProductType       string                 `json:"productType"`
	DiversityZone     string                 `json:"diversityZone"`
	NetworkInterfaces []MVENetworkInterface  `json:"vnics"`
	VendorConfig     VendorConfig `json:"VendorConfig"`
}

// VendorConfig is an interface for MVE vendor configuration.
type VendorConfig interface {
	IsVendorConfig()
}

// ArubaConfig represents the configuration for an Aruba MVE.
type ArubaConfig struct {
	VendorConfig
	Vendor string `json:"vendor"`
	ImageID int `json:"imageId"`
	ProductSize string `json:"productSize"`
	AccountName string `json:"accountName"`
	AccountKey string `json:"accountKey"`
}

// CiscoConfig represents the configuration for a Cisco MVE.
type CiscoConfig struct {
	VendorConfig
	Vendor string `json:"vendor"`
	ImageID int `json:"imageId"`
	ProductSize string `json:"productSize"`
	AdminSSHPublicKey string `json:"adminSshPublicKey"`
	CloudInit string `json:"cloudInit"`
}

// FortinetConfig represents the configuration for a Fortinet MVE.
type FortinetConfig struct {
	VendorConfig
	Vendor string `json:"vendor"`
	ImageID int `json:"imageId"`
	ProductSize string `json:"productSize"`
	AdminSSHPublicKey string `json:"adminSshPublicKey"`
	LicenseData string `json:"licenseData"`
}

// PaloAltoConfig represents the configuration for a Palo Alto MVE.
type PaloAltoConfig struct {
	VendorConfig
	Vendor string `json:"vendor"`
	ImageID int `json:"imageId"`
	ProductSize string `json:"productSize"`
	AdminSSHPublicKey string `json:"adminSshPublicKey"`
	AdminPasswordHash string `json:"adminPasswordHash"`
	LicenseData string `json:"licenseData"`
}

// VersaConfig represents the configuration for a Versa MVE.
type VersaConfig struct {
	VendorConfig
	Vendor string `json:"vendor"`
	ImageID int `json:"imageId"`
	ProductSize string `json:"productSize"`
	DirectorAddress string `json:"directorAddress"`
	ControllerAddress string `json:"controllerAddress"`
	LocalAuth string `json:"localAuth"`
	RemoteAuth string `json:"remoteAuth"`
	SerialNumber string `json:"serialNumber"`
}

// VmwareConfig represents the configuration for a VMware MVE.
type VmwareConfig struct {
	VendorConfig
	Vendor string `json:"vendor"`
	ImageID int `json:"imageId"`
	ProductSize string `json:"productSize"`
	AdminSSHPublicKey string `json:"adminSshPublicKey"`
	VcoAddress string `json:"vcoAddress"`
	VcoActivationCode string `json:"vcoActivationCode"`
}

// MVENetworkInterface represents a vNIC.
type MVENetworkInterface struct {
	Description string `json:"description"`
	VLAN        int    `json:"vlan"`
}

// InstanceSize encodes the available MVE instance sizes.
type MVEInstanceSize string

// MVE instance sizes.
const (
	MVE_SMALL  MVEInstanceSize = "SMALL"
	MVE_MEDIUM MVEInstanceSize = "MEDIUM"
	MVE_LARGE  MVEInstanceSize = "LARGE"
	MVE_XLARGE MVEInstanceSize = "X_LARGE_12"
)

// MVEOrderConfirmation represents the response to an MVE order request.
type MVEOrderConfirmation struct {
	TechnicalServiceUID string `json:"technicalServiceUid"`
}

// MVE represents a Megaport Virtual Edge from the Megaport MVE API.
type MVE struct {
	ID                    int                    `json:"productId"`
	UID                   string                 `json:"productUid"`
	Name                  string                 `json:"productName"`
	Type                  string                 `json:"productType"`
	ProvisioningStatus    string                 `json:"provisioningStatus"`
	CreateDate            *Time                  `json:"createDate"`
	CreatedBy             string                 `json:"createdBy"`
	TerminateDate         *Time                  `json:"terminateDate"`
	LiveDate              int                    `json:"liveDate"`
	Market                string                 `json:"market"`
	LocationID            int                    `json:"locationId"`
	UsageAlgorithm        string                 `json:"usageAlgorithm"`
	MarketplaceVisibility bool                   `json:"marketplaceVisibility"`
	VXCPermitted          bool                   `json:"vxcpermitted"`
	VXCAutoApproval       bool                   `json:"vxcAutoApproval"`
	SecondaryName         string                 `json:"secondaryName"`
	CompanyUID            string                 `json:"companyUid"`
	CompanyName           string                 `json:"companyName"`
	ContractStartDate     *Time                  `json:"contractStartDate"`
	ContractEndDate       *Time                  `json:"contractEndDate"`
	ContractTermMonths    int                    `json:"contractTermMonths"`
	AttributeTags         map[string]string      `json:"attributeTags"`
	Virtual               bool                   `json:"virtual"`
	BuyoutPort            bool                   `json:"buyoutPort"`
	Locked                bool                   `json:"locked"`
	AdminLocked           bool                   `json:"adminLocked"`
	Cancelable            bool                   `json:"cancelable"`
	Resources             *MVEResources `json:"resources"`
	Vendor                string                 `json:"vendor"`
	Size                  string                 `json:"mveSize"`
	NetworkInterfaces     []*MVENetworkInterface `json:"vnics"`
}

// MVEResources represents the resources associated with an MVE.
type MVEResources struct {
	Interface *PortInterface `json:"interface"`
	VirtualMachines []*MVEVirtualMachine `json:"virtual_machine"`
}

// MVEVirtualMachine represents a virtual machine associated with an MVE.
type MVEVirtualMachine struct {
	ID int `json:"id"`
	CpuCount int `json:"cpu_count"`
	Image *MVEVirtualMachineImage `json:"image"`
	ResourceType string `json:"resource_type"`
	Up bool `json:"up"`
	Vnics []*MVENetworkInterface `json:"vnics"`
}

// MVVEVirtualMachineImage represents the image associated with an MVE virtual machine.
type MVEVirtualMachineImage struct {
	ID int `json:"id"`
	Vendor string `json:"vendor"`
	Product string `json:"product"`
	Version string `json:"version"`
}

// MVEOrderResponse represents the response to an MVE order request.
type MVEOrderResponse struct {
	Message string                  `json:"message"`
	Terms   string                  `json:"terms"`
	Data    []*MVEOrderConfirmation `json:"data"`
}

// MVEResponse represents the response to an MVE request.
type MVEResponse struct {
	Message string `json:"message"`
	Terms   string `json:"terms"`
	Data    *MVE   `json:"data"`
}
