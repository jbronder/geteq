# geteq: A USGS Earthquake Event CLI
Retreive real-time and historical earthquake event records from the CLI.


## Tools & Technologies
- [Cobra CLI](https://cobra.dev/)


## About
This is an ongoing project to construct a simple command-line application in Go
using Cobra as the scaffold for the CLI. It parses commands and flags from users
to send a GET request to USGS's real-time and historical earthquake data API and
formats the response into a readable table to the terminal. It also forwards the
request to return CSV and JSON formats of the data for further data processing
and analysis.


## Sources
- [CSV Source](https://earthquake.usgs.gov/earthquakes/feed/v1.0/csv.php)
- [JSON Source](https://earthquake.usgs.gov/earthquakes/feed/v1.0/geojson.php)
- [FDSN API Source](https://earthquake.usgs.gov/fdsnws/event/1/)

