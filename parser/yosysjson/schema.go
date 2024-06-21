package yosysjson

const (
	Input  string = "input"
	Output string = "output"
	InOut  string = "inout"
)

type YosysJson struct {
	Creator string            `json:"creator"`
	Modules map[string]Module `json:"modules"`
	Models  map[string]Model  `json:"models"`
}

type Module struct {
	Attributes Attribute `json:"attributes"`
	//恐らくstringとintを受けられるようにしないといけなさそう　Anyなど
	ParameterDefaultValues map[string]string  `json:"parameter_default_values"`
	Ports                  map[string]Port    `json:"ports"`
	Cells                  map[string]Cell    `json:"cells"`
	Memories               map[string]Memory  `json:"memories"`
	NetName                map[string]NetName `json:"netnames"`
}

type Attribute struct {
	Hdlname string `json:"hdlname"`
	Top     string `json:"top"`
	Src     string `json:"src"`
	Init    string `json:"init"`
}

type Parameter struct{}

type Port struct {
	Direction string `json:"direction"`
	Bits      []int  `json:"bits"`
}

type Cell struct {
	HideName       int               `json:"hide_name"`
	Type           string            `json:"type"`
	Parameters     Parameter         `json:"parameters"`
	Attributes     Attribute         `json:"attributes"`
	PortDirections map[string]string `json:"port_directions"`
	Connections    map[string][]int  `json:"connections"`
}

type Memory struct {
}

type NetName struct {
	HideName   int       `json:"hide_name"`
	Bits       []int     `json:"bits"`
	Attributes Attribute `json:"attributes"`
}

type Model struct {
}

type Experimental interface {
}
