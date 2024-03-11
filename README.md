# sheet - a cli thing for messing with google sheets

### Getting Started

Install `sheet` into `$GOROOT/bin/`

```
go install github.com/gerrowadat/sheet@latest
```

Follow the instructions [here](https://developers.google.com/identity/protocols/oauth2) to obtain client credentials.
Make sure it has access to the Sheets API. Download the client secret file somewhere it can't be read by anyone else.

Set up `sheet` to point your config at the client secret file.

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

### Usage Examples
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

#### Configuration

```
# See configuration items available.
sheet config get

# Get one config items
sheet config get read-chunksize

# Set config items.
sheet config set read-chunksize 500
```

#### Workbook/Worksheet etc. metadata
```
# Get the list of worksheet in a workbook
sheet ls SpReAdShEeTiDfRoMUrL 
```

#### Reading Data
```
# Get a range and spit it out as CSV
sheet get SpReAdShEeTiDfRoMUrL 'myworksheet!B3:F8'

# Print the last 5 populated rows of a worksheet
sheet tail SpReAdShEeTiDfRoMUrL 'myworksheet' --lines=5

# Output an entire worksheet
sheet cat SpReAdShEeTiDfRoMUrL myworksheet
```

#### Modifying Spreadsheet Info
```
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

### Aliases

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

### Coming....when.

TODO (See [issues](https://github.com/gerrowadat/sheet/issues) for tracking.:

```
# Writing
sheet put <id> <datarange>
sheet replace/append <id> <worksheet>

# Etc.
sheet rm
sheet cp
```
