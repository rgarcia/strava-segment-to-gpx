package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"os"

	strava "github.com/strava/go.strava"
)

type gpx struct {
	Version        string   `xml:"version,attr"`
	Creator        string   `xml:"creator,attr"`
	Xmlns          string   `xml:"xmlns,attr"`
	Xsi            string   `xml:"xsi,attr"`
	SchemaLocation string   `xml:"schemaLocation,attr"`
	Metadata       Metadata `xml:"metadata"`
	Trk            Trk      `xml:"trk"`
}
type Metadata struct {
	Name   string `xml:"name"`
	Author Author `xml:"author"`
	Link   Link   `xml:"link"`
}
type Link struct {
	Href string `xml:"href,attr"`
}
type Author struct {
	Name string `xml:"name"`
	Link Link   `xml:"link"`
}
type Trk struct {
	Name   string `xml:"name"`
	Link   Link   `xml:"link"`
	Type   string `xml:"type"`
	TrkSeg TrkSeg `xml:"trkseg"`
}
type TrkSeg struct {
	TrkPts []TrkPt `xml:"trkpt"`
}

type TrkPt struct {
	Lon string `xml:"lon,attr"`
	Lat string `xml:"lat,attr"`
	Ele string `xml:"ele,omitempty"`
}

// GpxDefaults contains boilerplate that stays the same for all segments.
var GpxDefaults = gpx{
	Version:        "1.1",
	Creator:        "strava-segment-to-gpx",
	Xmlns:          "http://www.topografix.com/GPX/1/1",
	Xsi:            "http://www.w3.org/2001/XMLSchema-instance",
	SchemaLocation: "http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/1/1/gpx.xsd http://www.garmin.com/xmlschemas/GpxExtensions/v3 http://www.garmin.com/xmlschemas/GpxExtensionsv3.xsd http://www.garmin.com/xmlschemas/TrackPointExtension/v1 http://www.garmin.com/xmlschemas/TrackPointExtensionv1.xsd http://www.garmin.com/xmlschemas/GpxExtensions/v3 http://www.garmin.com/xmlschemas/GpxExtensionsv3.xsd http://www.garmin.com/xmlschemas/TrackPointExtension/v1 http://www.garmin.com/xmlschemas/TrackPointExtensionv1.xsd http://www.garmin.com/xmlschemas/GpxExtensions/v3 http://www.garmin.com/xmlschemas/GpxExtensionsv3.xsd http://www.garmin.com/xmlschemas/TrackPointExtension/v1 http://www.garmin.com/xmlschemas/TrackPointExtensionv1.xsd",
	Metadata: Metadata{
		Author: Author{
			Name: "strava-segment-to-gpx",
			Link: Link{Href: "https://github.com/rgarcia/strava-segment-to-gpx"},
		},
	},
}

func main() {
	var segmentID int64
	var accessToken string
	flag.Int64Var(&segmentID, "id", 229781, "Strava Segment Id")
	flag.StringVar(&accessToken, "token", "", "Access Token")
	flag.Parse()
	if accessToken == "" {
		fmt.Fprintf(os.Stderr, "Please provide an access_token, one can be found at https://www.strava.com/settings/api\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	client := strava.NewClient(accessToken)
	segment, err := strava.NewSegmentsService(client).Get(segmentID).Do()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get segment from Strava API: %s", err)
		os.Exit(1)
	}

	gpx := GpxDefaults
	gpx.Metadata.Name = segment.Name
	gpx.Metadata.Link = Link{Href: fmt.Sprintf("https://www.strava.com/segments/%d", segment.Id)}
	gpx.Trk.Name = segment.Name
	gpx.Trk.Link = Link{Href: fmt.Sprintf("https://www.strava.com/segments/%d", segment.Id)}
	gpx.Trk.Type = string(segment.ActivityType)
	for _, pt := range segment.Map.Polyline.Decode() {
		gpx.Trk.TrkSeg.TrkPts = append(gpx.Trk.TrkSeg.TrkPts, TrkPt{
			Lat: fmt.Sprintf("%f", pt[0]),
			Lon: fmt.Sprintf("%f", pt[1]),
		})
	}

	fmt.Fprintf(os.Stdout, xml.Header)
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", " ")
	if err := enc.Encode(gpx); err != nil {
		fmt.Printf("error: %v\n", err)
	}
}
