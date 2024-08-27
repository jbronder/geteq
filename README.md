# geteq: A USGS Earthquake Event CLI
Retrieve real-time feed and historical earthquake event records from the CLI.


## Tools & Technologies
- [Cobra CLI Framework](https://cobra.dev/)


## About
This is an ongoing project to construct a simple command-line application in Go
using Cobra as the scaffold for the CLI. It parses commands and flags from users
to send an HTTP request to USGS's real-time feed and historical earthquake data
API and formats the response into a readable table to the terminal. It also
forwards the request to return CSV and JSON formats of the data for further data
processing and analysis.  data pipeline


## Sources
- [Real-time Feed CSV Source](https://earthquake.usgs.gov/earthquakes/feed/v1.0/csv.php)
- [Real-time Feed JSON Source](https://earthquake.usgs.gov/earthquakes/feed/v1.0/geojson.php)
- [Historical FDSN Source](https://earthquake.usgs.gov/fdsnws/event/1/)


## Real-Time Queries
Using the `realtime` subcommand, real-time data requests retrieve earthquake
events from as far back as 30 days from the current time of the request.
Real-time feed queries return a formatted table of events to the terminal. 

Events can be queried within the following time intervals:
- `-t {month, week, day, hour}`

Real-time queries have the following magnitude options available:
- `-m {all, 1.0, 2.5, 4.5, major}`

Queries output into the following formats:
- `-o {csv, json, table}` where `table` is a prettier format to view event
  records in the terminal


### Real-time Feed Query Examples
Retrieve records of significant earthquakes from the past month:
```bash
$ geteq realtime -m major -t month
$ geteq realtime # "major" and "month" are the current defaults
```

Retrieve records of earthquakes with magnitudes greater than or equal to 4.5 in
the past hour formatted to output JSON:
```bash
$ geteq rt -m 4.5 -t hour -o json
```


## Historical Queries
The `fdsn` subcommand currently allows for searching earthquake catalogs bounded
between date ranges and/or magnitudes or magnitude ranges. It also provides
support for querying individual earthquake records for detailed event
information such as the event's depth, Did You Feel It (DYFI) report counts, the
type of magnitude was computed, Modified Mercalli Intensity (MMI), etc.


### Historical Query Examples
Retrieve records of earthquake events between January 2, 2024 to January 3, 2024: 
```bash
$ geteq fdsn query -t 2024-01-02,2024-01-03
```

Retrieve records with magnitude range between 4.5 to 7.5 and formatted to output
JSON:
```bash
$ geteq fdsn q -m 4.5,7.5 -o json # q is an alias for query
```

Retrieve records with magnitudes greater than 4.5 and formatted to output CSV:
```bash
$ geteq fdsn q -m ">4.5" -o csv # NOTE: quotes are required
```
Retrieve details of a single event using an `eventid`:
```bash
$ geteq fdsn query event uw10530748 # where uw10530748 is an eventid
$ geteq fdsn q e uw10530748 # where e is an alias for event
```