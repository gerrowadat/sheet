package sheet

import (
	"fmt"

	"github.com/spf13/viper"
)

func aliasValid(name string, spec *DataSpec) error {
	// TODO
	return nil
}

func SetAlias(name string, spec *DataSpec) error {
	if err := aliasValid(name, spec); err != nil {
		return err
	}
	// Remove the alias if it exists
	if err := DeleteAlias(name); err != nil {
		return err
	}
	if spec.Workbook != "" {
		viper.Set("aliases."+name+".workbook", spec.Workbook)
	}
	if spec.Worksheet != "" {
		viper.Set("aliases."+name+".worksheet", spec.Worksheet)
	}
	if spec.Range != "" {
		viper.Set("aliases."+name+".range", spec.Range)
	}
	viper.WriteConfig()
	return nil
}

func GetAlias(name string) (*DataSpec, error) {
	ret := &DataSpec{}
	all := viper.GetStringMap("aliases")
	for k := range all {
		if k == name {
			alias := all[k].(map[string]interface{})
			for k, v := range alias {
				if k == "workbook" {
					ret.Workbook = v.(string)
				}
				if k == "worksheet" {
					ret.Worksheet = v.(string)
				}
				if k == "range" {
					ret.Range = v.(string)
				}
			}
			return ret, nil
		}
	}
	return nil, fmt.Errorf("alias not found: %v", name)
}

func GetAllAliases() map[string]*DataSpec {
	ret := map[string]*DataSpec{}
	all := viper.GetStringMap("aliases")
	for k := range all {
		alias := all[k].(map[string]interface{})
		ret[k] = &DataSpec{}
		for k, v := range alias {
			if k == "workbook" {
				ret[k].Workbook = v.(string)
			}
			if k == "worksheet" {
				ret[k].Worksheet = v.(string)
			}
			if k == "range" {
				ret[k].Range = v.(string)
			}
		}
	}
	return ret
}

func DeleteAlias(name string) error {
	viper.Set("aliases."+name, nil)
	return nil
}
