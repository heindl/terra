package terra

import (
	"encoding/json"
	"fmt"

	"github.com/paulsmith/gogeos/geos"
)

type geoJSONEncodeType struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Properties  map[string]interface{} `json:"properties"`
	Classifiers []*Classifier          `json:"classifiers"`
	Geometry    struct {
		Type        string        `json:"type"`
		Coordinates []interface{} `json:"coordinates"`
	} `json:"geometry"`
}

func (feat *Feature) ToJSON() (geojson []byte, er error) {

	if feat.IsEmpty() {
		er = fmt.Errorf("The feature is empty, with nothing to encode into GeoJSON.")
		log.Warning(er.Error())
		return
	}

	var construct = &geoJSONEncodeType{
		ID:          feat.ID,
		Type:        "Feature",
		Properties:  feat.Properties,
		Classifiers: feat.Classifiers,
	}
	construct.Geometry.Type = feat.Type

	switch {
	case feat.Type == "Point":
		var coords []geos.Coord
		if coords, er = feat.Geometry.Coords(); er != nil {
			log.Error(er.Error())
			return
		}
		construct.Geometry.Coordinates = encodeCoord(coords[0])
	case feat.Type == "LineString":
		if construct.Geometry.Coordinates, er = encodeLineString(feat.Geometry); er != nil {
			return
		}
	case feat.Type == "Polygon":
		if construct.Geometry.Coordinates, er = encodePolygon(feat.Geometry); er != nil {
			return
		}
	case feat.Type == "MultiPolygon":
		if construct.Geometry.Coordinates, er = encodeMultiPolygon(feat.Geometry); er != nil {
			return
		}
	default:
		er = fmt.Errorf("Currently GeoJSON must be type Point, Linestring, Polygon, or Multipolygon. Found %s.", feat.Type)
		log.Error(er.Error())
		return
	}

	if geojson, er = json.Marshal(construct); er != nil {
		log.Error(er.Error())
	}

	return

}

func encodeCoord(coord geos.Coord) []interface{} {
	return []interface{}{coord.X, coord.Y}
}

func encodeLineString(geometry *geos.Geometry) (response []interface{}, er error) {

	var coords []geos.Coord
	if coords, er = geometry.Coords(); er != nil {
		log.Error(er.Error())
		return
	}

	for i := range coords {
		response = append(response, encodeCoord(coords[i]))
	}

	return
}

func encodePolygon(geometry *geos.Geometry) (response []interface{}, er error) {

	var geometries []*geos.Geometry

	shell, er := geometry.Shell()
	if er != nil {
		log.Error(er.Error())
		return
	}
	geometries = append(geometries, shell)

	holes, er := geometry.Holes()
	if er != nil {
		log.Error(er.Error())
		return
	}
	geometries = append(geometries, holes...)

	for i := range geometries {
		var coordProgression []interface{}
		if coordProgression, er = encodeLineString(geometries[i]); er != nil {
			return
		}
		response = append(response, coordProgression)
	}

	return
}

func encodeMultiPolygon(multipolygon *geos.Geometry) (response []interface{}, er error) {

	var collectionLength int
	if collectionLength, er = multipolygon.NGeometry(); er != nil {
		log.Error(er.Error())
		return
	}

	for i := 0; i < collectionLength; i++ {
		var g *geos.Geometry
		var compiled []interface{}
		if g, er = multipolygon.Geometry(i); er != nil {
			log.Error(er.Error())
			return
		}
		if compiled, er = encodePolygon(g); er != nil {
			return
		}
		response = append(response, compiled)
	}

	return
}
