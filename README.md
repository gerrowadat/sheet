# sheet - a CLI tool and library for working with Google Sheets

## Getting Started

Install `sheet` into `$GOROOT/bin/`:

```
go install github.com/gerrowadat/sheet@latest
```

Follow the instructions [here](https://developers.google.com/identity/protocols/oauth2) to obtain client credentials.
Make sure it has access to the Sheets API. Download the client secret file somewhere it can't be read by anyone else.

Set up `sheet` to point your config at the client secret file:

```
sheet config set clientsecretfile /path/to/clientsecrets.json
```

Also set up somewhere to save your authentication token (you don't have one yet):
```
sheet config set authtokenfile /path/to/authtoken.json
```

You should then be set up with access. The first time you issue a command that tries to reach Sheets,
you'll be pointed at a URL to visit as the logged-in user - the approval flow will redirect to a localhost URL
that will have a token in it -- paste this token into the CLI when asked.

## Using `sheet` as a Library

The `lib` directory provides the `sheet` package, which you can use directly in your own Go programs to read, write, and manage Google Sheets data. See `examples/library_usage.go` for a complete working example.

### Installation

```bash
go get github.com/gerrowadat/sheet
```

Then import the library:

```go
import "github.com/gerrowadat/sheet/lib"
```

The package name is `sheet`, so you'll reference it as `sheet.GetService()`, `sheet.DataSpec{}`, etc.

### Authentication

There are two ways to authenticate:

**Option 1: Using viper configuration (same as the CLI)**

If you configure viper with `clientsecretfile` and `authtokenfile` keys (as the CLI does), you can call `GetService()` directly:

```go
srv, err := sheet.GetService()
if err != nil {
    log.Fatal(err)
}
```

**Option 2: Direct authentication (no viper dependency)**

If you want to manage credentials yourself, use `GetClient()` to get an HTTP client, then create the Sheets service:

```go
import (
    "context"
    "google.golang.org/api/option"
    "google.golang.org/api/sheets/v4"
)

client := sheet.GetClient("/path/to/client_secret.json", "/path/to/token.json")
srv, err := sheets.NewService(context.Background(), option.WithHTTPClient(client))
if err != nil {
    log.Fatal(err)
}
```

### Core Types

#### DataSpec

`DataSpec` is the central type for referencing locations in Google Sheets. It represents a workbook, worksheet, or range:

```go
// Reference a workbook
spec := &sheet.DataSpec{Workbook: "spreadsheet-id"}
spec.IsWorkbook()  // true

// Reference a worksheet
spec := &sheet.DataSpec{Workbook: "spreadsheet-id", Worksheet: "Sheet1"}
spec.IsWorksheet() // true

// Reference a range
spec := &sheet.DataSpec{
    Workbook:  "spreadsheet-id",
    Worksheet: "Sheet1",
    Range:     sheet.RangeFromString("A1:C10"),
}
spec.IsRange() // true

// Get the in-sheet reference string (e.g., "Sheet1!A1:C10")
ref := spec.GetInSheetDataSpec()
```

You can also parse arguments the same way the CLI does:

```go
spec, err := sheet.ExpandArgsToDataSpec([]string{"spreadsheet-id", "Sheet1!A1:C10"})
```

#### DataRange

`DataRange` represents a cell range in spreadsheet notation:

```go
r := sheet.RangeFromString("A1:C10")
cols, rows := r.SizeXY()     // 3, 10
fixed := r.IsFixedSize()     // true (both row and col bounds are set)
s := r.String()              // "A1:C10"
```

#### DataFormat

`DataFormat` represents CSV or TSV:

```go
sheet.CsvFormat // "csv"
sheet.TsvFormat // "tsv"
```

### Reading Data

Use the Google Sheets service with a `DataSpec` to read data, then format it:

```go
srv, _ := sheet.GetService()
spec := &sheet.DataSpec{Workbook: "spreadsheet-id", Worksheet: "Sheet1"}

resp, err := srv.Spreadsheets.Values.Get(spec.Workbook, spec.GetInSheetDataSpec()).Do()
if err != nil {
    log.Fatal(err)
}

// Format as CSV or TSV
csvOutput := sheet.FormatValues(resp, sheet.CsvFormat)
fmt.Print(csvOutput)

// Or print directly to stdout
sheet.PrintValues(resp, sheet.TsvFormat)
```

### Writing Data

```go
data := [][]string{
    {"Name", "Age", "City"},
    {"Alice", "30", "NYC"},
    {"Bob", "25", "SF"},
}

// Write to an entire worksheet (clears existing data first)
// The protect and force flags control worksheet protection behavior
err := sheet.WriteDataToWorksheet(srv, spec, data, false, false)

// Write to a specific range (data must fit within the range)
spec.Range = sheet.RangeFromString("A1:C3")
err = sheet.WriteDataToRange(srv, spec, data)
```

### Clearing Data

```go
// Clear an entire worksheet (respects protection settings)
err := sheet.ClearWorksheet(srv, spec, false, false)

// Clear a specific range
err = sheet.ClearRange(srv, spec)
```

### Data Format Conversion

Convert between CSV/TSV strings and Go data structures:

```go
import (
    "bufio"
    "strings"
)

// Parse CSV/TSV input into [][]string
reader := bufio.NewReader(strings.NewReader("a,b,c\nd,e,f"))
data, err := sheet.ScanValues(reader, sheet.CsvFormat)
// data: [][]string{{"a", "b", "c"}, {"d", "e", "f"}}

// Format Google Sheets ValueRange as CSV/TSV string
output := sheet.FormatValues(valueRange, sheet.TsvFormat)
```

### Aliases

Aliases provide named shortcuts to workbooks, worksheets, and ranges (stored via viper config):

```go
// Set an alias
spec := &sheet.DataSpec{Workbook: "spreadsheet-id", Worksheet: "MyData"}
err := sheet.SetAlias("mydata", spec)

// Get an alias
spec, err := sheet.GetAlias("mydata")

// List all aliases
for name, spec := range sheet.GetAllAliases() {
    fmt.Printf("%s => %s\n", name, spec.String())
}

// Delete an alias
err = sheet.DeleteAlias("mydata")
```

## CLI Usage Examples
(These examples will get more useful as functionality improves)

```
# Create a new workbook, add a worksheet, and set an alias to a given range
WORKBOOK_ID=`sheet touch workbook "My Customers"`
sheet touch $WORKBOOK_ID "Customer Data"

sheet alias client_table $WORKBOOK_ID "Customer Data" A1:C100

# ... go populate some data somehow.

# Later on, go get it.
sheet get @client_table
```


### Commands

A 'workbook' is a top-level spreadsheet (identified by the ID from the URL).
A 'worksheet' is a tabbed sheet within a workbook.
A 'range' is a range, seriously.

#### Configuration - `config set`/`config get`

```
# See configuration items available.
sheet config get

# Get one config items
sheet config get read-chunksize

# Set config items.
sheet config set read-chunksize 500
```

#### Workbook/Worksheet etc. metadata - `ls`
```
# Get the list of worksheet in a workbook
sheet ls SpReAdShEeTiDfRoMUrL 
```

#### Reading Data - `get`/`tail`/`cat`
```
# Get a range and spit it out as CSV
sheet get SpReAdShEeTiDfRoMUrL 'myworksheet!B3:F8'

# Print the last 5 populated rows of a worksheet (default is 10)
sheet tail SpReAdShEeTiDfRoMUrL 'myworksheet' --lines=5

# Output an entire worksheet
sheet cat SpReAdShEeTiDfRoMUrL myworksheet
```

#### Modifying Spreadsheet Info - `touch`/`rm`
```
# touch
# Create a new workbook or worksheet, and output the ID
sheet touch workbook "My Workbook"
sheet touch workbook # Will use --default-workbook-title or your config key 'default-workbook-title' (in that order).

# Create a new worksheet inside an existing workbook
sheet touch MyWorkBoOk mynewsheet
sheet touch @mywb mynewsheet
# Or you can do this (aliases and refer to non-existent entities)
sheet alias set mons MyWorkBoOk myothernewsheet
sheet touch @mons
```

```
# rm
# rm doesn't work on workbooks, but does on worksheets and ranges.

# Delete the named sheet.
sheet rm mywrokbook mysheet
sheet rm @mysheetalias

# Clear the data inside this range (stuff like conditional formatting etc. stays).
sheet rm @myworkbook 'junk!A10:F100'
```

#### Modifying Data - `put`
```
# put
# Put by default reads data in the format specified from stdin and writes it to the specified place.
# Only works on worksheets and ranges. Will delete all data in the sheet or range first
# (Yes, even if you're specifying less data than was present before)

# Write a csv file into a sheet.
sheet put MyWorkBoOk mysheet < data.csv

# Write a couple of values into a range
echo "a,b,c" | sheet put MyWoRkBoOk 'mysheet!A1:C1'

# This won't work -- the specified data must fit in the specified range
echo "a,b,c" | sheet put MyWoRkBoOk 'mysheet!A1:B1'

# This will work, and will clear the 'c' from cell C1 
echo "a,b" | sheet put MyWoRkBoOk 'mysheet!A1:C1'

# Hey, let's put something in Cell C1!
echo "bees" | sheet put MyWoRkBoOk 'mysheet!C1:C1'

# This will copy the cells we're working on to the row below
sheet get MyWoRkBoOk 'mysheet!A1:C1' | sheet put MyWoRkBoOk 'mysheet!A2:C2'
```


### Aliases - `alias get`/`alias set`

You can set aliases for workbooks, worksheets and even ranges, then refer to them with the @ prefix (you can also configure this with `--alias-spec-prefix` or `sheet config set alias-spec-prefix &` or whatever)

```
# See my current aliases
sheet alias get

# Set an alias to a workbook, worksheet or range
sheet alias mywork SprEaDsHeTiD
sheet alias mysheet SpReAdShEeTiD clients
sheet alias mybestbits SpReAdShEeTiD 'clients!A1:B5'

# Now I can do these instead of pasting in the spreadsheet ID every time or relying on shell history:
sheet ls @mywork # show all sheet names in SprEaDsHeTiD
sheet cat @mysheet # Output all of the 'clients' sheet
sheet get @mybestbits # Output just the specified range

# I can even do this!
sheet get '@mysheet!A1:B10'
```

### Flags and Config

The following flags are generally supported, they override the config entries of the same name
(without the `--`, of course) that can be set with `sheet config set` per above.

#### `--input-format` and `--output-format`

As you might guess, specifies the format to use for input and output. Currently only supports 'csv' and 'tsv'.

#### `--authtokenfile` and `--clientsecretfile`

Specify where your oauth 2.0 client secrets and token file go.

#### `--read-chunksize` and `--write-chunksize`

Specify the amount of data to be read from or written to a sheet at a time, in rows.

### Cookbook

A couple of terrifying use cases I've either used or have considered using.

```
# Storing a blob in a spreadsheet cell for whatever reason
MY_WBID=`sheet touch workbook blobs`
sheet touch worksheet $MY_WBID files
sheet alias set blobcell $MY_WBID 'files!A1:A1'
base64 -w0 mybinaryfile | sheet put @blobcell
# Get the blob later
sheet get @blobcell | base64 -d -
```

```
# Dump a disk usage report to a sheet because reasons I guess?
# @myreport is an alias to a worksheet
du -S . | tr '\t' ',' | sheet put @myreport
# Actually, no need for tr.
du -S . | sheet put @myreport --input-format=tsv
```

### Coming....when.

TODO (See [issues](https://github.com/gerrowadat/sheet/issues) for tracking.:

```
# Writing
sheet replace/append <id> <worksheet>

# Etc.
sheet cp
```
