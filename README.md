# strava-segment-to-gpx

Converts a Strava segment to gpx.

## Installation

```
go get github.com/rgarcia/strava-segment-to-gpx
```

This will put `strava-segment-to-gpx` in `$(GOPATH)/bin`.

## Usage

You'll need to create a Strava app to get an [access token](https://www.strava.com/settings/api).

```
strava-segment-to-gpx -token <strava access token> -id <segment ID> > output.gpx
```
