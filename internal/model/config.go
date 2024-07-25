package model

import (
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Email types.String `tfsdk:"email" yaml:"email"`
	Token types.String `tfsdk:"token" yaml:"token"`
}

type LocalConfig struct {
	Spec LocalConfigSpec `yaml:"spec"`
}

type LocalConfigSpec struct {
	Email string `yaml:"email"`
	Token string `yaml:"token"`
}

func (c *Config) From(d diag.Diagnostics) {
	p, err := homedir.Expand("~/.plural/config.yml")
	if err != nil {
		d.AddError("Client Error", fmt.Sprintf("Could not find local plural config: %s", err))
		return
	}

	res, err := os.ReadFile(p)
	if err != nil {
		d.AddError("Client Error", fmt.Sprintf("Could not read local plural config: %s", err))
		return
	}

	var config LocalConfig
	if err := yaml.Unmarshal(res, &config); err != nil {
		d.AddError("Client Error", fmt.Sprintf("Could not parse local plural config: %s", err))
		return
	}

	c.Email = types.StringValue(config.Spec.Email)
	c.Token = types.StringValue(config.Spec.Token)
}
