# sheet - a cli thing for messing with google sheets

### The Basics


```
# Pass --clientsecretfile and --authtokenfile to the following:

# Get a range and spit it out as CSV
sheet get SpReAdShEeTiDfRoMUrL 'myworksheet!B3:F8'

# Print the last 5 populated rows of a worksheet
sheet tail SpReAdShEeTiDfRoMUrL 'myworksheet' --lines=5

# List the worksheets in the specified sheet
sheet ls SpReAdShEeTiDfRoMUrL 

# Output an entire worksheet
sheet ls SpReAdShEeTiDfRoMUrL myworksheet
```

See `sheet help get` (and so on) for flags, you'll need a client secret file, per the docs.

### Aliases

You can set aliases for workbooks, worksheets and even ranges:

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
sheet touch
sheet rm
sheet cp
```
