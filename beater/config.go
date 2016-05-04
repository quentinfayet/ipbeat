package beater

// IPConfig represents configuration for Ipbeat
type IPConfig struct {
	Period *int64
}

// ConfigSettings represents the content of ipbeat.yml
type ConfigSettings struct {
	Ipbeat *IPConfig `config:"ipbeat"`
	Input  *IPConfig `config:"input"`
}
