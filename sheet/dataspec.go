package sheet

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type DataSpec struct {
	Workbook  string
	Worksheet string
	Range     string
}

func (d *DataSpec) GetInSheetDataSpec() string {
	if d.Worksheet != "" && d.Range != "" {
		return fmt.Sprintf("%v!%v", d.Worksheet, d.Range)
	} else {
		return fmt.Sprintf("%v%v", d.Worksheet, d.Range)
	}
}

func (d *DataSpec) FromString(s string) *DataSpec {
	// This will always be datasheet, or datasheet!range
	if strings.Contains(s, "!") {
		fragments := strings.Split(s, "!")
		d.Worksheet = fragments[0]
		d.Range = fragments[1]
	} else {
		d.Worksheet = s
	}
	return d
}

func ExpandArgsToDataSpec(args []string) (*DataSpec, error) {
	// Expands an alias within an argument list (if it exists).
	// 'args' is the set of arguments that should represent a DataSpec (see dataspec.go)
	// If the first argument is an alias, it is expanded into a DataSpec and returned.
	if len(args) > 2 {
		return nil, fmt.Errorf("too many arguments when expanding DataSpec: %v", args)
	}

	if len(args) == 0 {
		return &DataSpec{}, nil
	}

	if len(args) == 1 {
		// If there's only one argument, it can be an alias or a workbook ID.
		if strings.HasPrefix(args[0], "@") {
			// Expand alias, if it exists.
			return dataSpecFromAlias(args[0][1:])
		} else {
			// First non-alias argument is always a workbook ID.
			return &DataSpec{Workbook: args[0]}, nil
		}
	}
	// If there are two arguments, each of them may be an alias. We partially populate a DataSpec
	// from all aliases and bare arguments -- they cannot overlap.
	//
	// Possible 2-arg examples:
	// - wb ws
	// - wb ws!r
	// - @wb ws
	// - @wb ws!r
	// ... and so forth.
	specs := []*DataSpec{}
	for i, arg := range args {
		if strings.HasPrefix(arg, "@") {
			spec, err := dataSpecFromAlias(args[0][1:])
			if err != nil {
				return nil, err
			}
			specs = append(specs, spec)
		} else {
			spec := &DataSpec{}
			if i == 0 {
				spec.Workbook = arg
			} else {
				spec.FromString(arg)
			}
			specs = append(specs, spec)
		}
	}
	return mergeDataSpecs(specs)
}

func mergeDataSpecs(specs []*DataSpec) (*DataSpec, error) {
	// Merge all DataSpecs into one.
	// If there any fields overlapping, return an error.
	ret := DataSpec{}
	for _, spec := range specs {
		if spec.Workbook != "" {
			if ret.Workbook != "" {
				return nil, fmt.Errorf("multiple workbooks in DataSpecs: %v", specs)
			}
			ret.Workbook = spec.Workbook
		}
		if spec.Worksheet != "" {
			if ret.Worksheet != "" {
				return nil, fmt.Errorf("multiple worksheets in DataSpecs: %v", specs)
			}
			ret.Worksheet = spec.Worksheet
		}
		if spec.Range != "" {
			if ret.Range != "" {
				return nil, fmt.Errorf("multiple ranges in DataSpecs: %v", specs)
			}
			ret.Range = spec.Range
		}
	}
	return &ret, nil
}

func dataSpecFromAlias(aliasname string) (*DataSpec, error) {
	ret := DataSpec{}

	// Handle @myalias!range
	if strings.Contains(aliasname, "!") {
		fragments := strings.Split(aliasname, "!")
		aliasname = fragments[0]
		ret.Range = fragments[1]
	}

	all := viper.GetStringMap("aliases")
	alias := map[string]interface{}{}
	for k := range all {
		if k == aliasname {
			alias = all[k].(map[string]interface{})
		}
	}
	if len(alias) == 0 {
		return nil, fmt.Errorf("alias not found: %v", aliasname)
	}

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

	// Do this check at the end -- if somehow we're specifying @myalias!range
	// and @myalias is a workbook alias, we've got an incomplete dataspec.
	if ret.Range != "" {
		if ret.Worksheet == "" {
			return nil, fmt.Errorf("invalid alias for ! notation: %v", aliasname)
		}
	}
	return &ret, nil
}
