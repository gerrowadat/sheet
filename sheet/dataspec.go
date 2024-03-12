package sheet

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type DataRange struct {
	StartRow int
	StartCol int
	EndRow   int
	EndCol   int
}

func colToLetter(col int) string {
	// Convert a column number to a letter.
	// e.g. 1 -> A, 2 -> B, 27 -> AA, 28 -> AB
	ret := ""
	for col > 0 {
		col--
		ret = string(rune('A'+col%26)) + ret
		col = col / 26
	}
	return ret
}

func letterToCol(letter string) int {
	// Convert a letter to a column number.
	// e.g. A -> 1, B -> 2, AA -> 27, AB -> 28
	ret := 0
	for _, c := range letter {
		ret = ret*26 + int(c) - int('A') + 1
	}
	return ret
}

func (d *DataRange) String() string {
	ret := colToLetter(d.StartCol)
	if d.StartRow > 0 {
		ret += fmt.Sprintf("%v", d.StartRow)
	}
	ret += ":"
	ret += colToLetter(d.EndCol)
	if d.EndRow > 0 {
		ret += fmt.Sprintf("%v", d.EndRow)
	}
	return ret
}

func splitRangeFragment(s string) (int, int, error) {
	// Given a string like "A1" or "B2", return the column and row.
	colstr := ""
	rowstr := ""
	for _, c := range s {
		if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') {
			colstr += string(c)
		} else if c >= '0' && c <= '9' {
			rowstr += string(c)
		} else {
			return 0, 0, fmt.Errorf("invalid range fragment: %v", s)
		}
	}

	// Handle "A:D" or "1:2" or "A:2" or "B:A" etc.
	var row, col int
	var err error
	if len(rowstr) > 0 {
		row, err = strconv.Atoi(rowstr)
		if err != nil {
			return 0, 0, err
		}
	}
	if len(colstr) == 0 {
		col = 0
	} else {
		col = letterToCol(colstr)
	}

	return col, row, nil
}

func (d *DataRange) FromString(s string) (*DataRange, error) {
	fragments := strings.Split(s, ":")
	if len(fragments) != 2 {
		return nil, fmt.Errorf("invalid range: %v", s)
	}
	startc, startr, err := splitRangeFragment(fragments[0])
	if err != nil {
		return nil, err
	}
	endc, endr, err := splitRangeFragment(fragments[1])
	if err != nil {
		return nil, err
	}
	d.StartRow = startr
	d.StartCol = startc
	d.EndRow = endr
	d.EndCol = endc
	return d, nil
}

func (d *DataRange) SizeXY() (int, int) {
	var col, row int
	// For a whole row or column, return 0 for the end of the range
	if d.EndCol > 0 {
		col = d.EndCol - d.StartCol + 1
	}
	if d.EndRow > 0 {
		row = d.EndRow - d.StartRow + 1
	}
	return col, row
}

type DataSpec struct {
	Workbook  string
	Worksheet string
	Range     DataRange
}

func (d *DataSpec) GetInSheetDataSpec() string {
	// Return a string that can be used to reference this DataSpec in a sheet.
	// e.g. "Sheet1!A1:B2"
	if d.Worksheet != "" {
		if d.Range != (DataRange{}) {
			return fmt.Sprintf("%v!%v", d.Worksheet, d.Range.String())
		} else {
			return d.Worksheet
		}
	} else {
		return d.Range.String()
	}
}

func (d *DataSpec) IsWorkbook() bool {
	return (d.Workbook != "" && d.Worksheet == "" && d.Range == DataRange{})
}

func (d *DataSpec) IsWorksheet() bool {
	return (d.Workbook != "" && d.Worksheet != "" && d.Range == DataRange{})
}

func (d *DataSpec) IsRange() bool {
	return (d.Workbook != "" && d.Worksheet != "" && d.Range != DataRange{})
}

func (d *DataSpec) String() string {
	ret := []string{}
	if d.Workbook != "" {
		ret = append(ret, "Workbook: "+d.Workbook)
	}
	if d.Worksheet != "" {
		ret = append(ret, "Worksheet: "+d.Worksheet)
	}
	if d.IsRange() {
		ret = append(ret, "Range: "+d.Range.String())
	}
	return strings.Join(ret, ", ")
}

func (d *DataSpec) FromString(s string) *DataSpec {
	// This will always be datasheet, or datasheet!range
	if strings.Contains(s, "!") {
		fragments := strings.Split(s, "!")
		d.Worksheet = fragments[0]
		_, err := d.Range.FromString(fragments[1])
		if err != nil {
			return nil
		}
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

	alias_prefix := viper.GetString("alias-spec-prefix")

	if alias_prefix == "" {
		alias_prefix = "@"
	}

	if len(args) == 1 {
		// If there's only one argument, it can be an alias or a workbook ID.
		if strings.HasPrefix(args[0], alias_prefix) {
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
		if strings.HasPrefix(arg, alias_prefix) {
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
		if (spec.Range != DataRange{}) {
			if (ret.Range != DataRange{}) {
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
		_, err := ret.Range.FromString(fragments[1])
		if err != nil {
			return nil, err
		}
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
			_, err := ret.Range.FromString(v.(string))
			if err != nil {
				return nil, err
			}
		}
	}

	// Do this check at the end -- if somehow we're specifying @myalias!range
	// and @myalias is a workbook alias, we've got an incomplete dataspec.
	if (ret.Range != DataRange{}) {
		if ret.Worksheet == "" {
			return nil, fmt.Errorf("invalid alias for ! notation: %v", aliasname)
		}
	}
	return &ret, nil
}
