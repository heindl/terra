package terra

import (
	"encoding/json"
	"github.com/paulsmith/gogeos/geos"
	"github.com/saleswise/errors/errors"
)

type geoJSONEncodeType struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Properties  map[string]interface{} `json:"properties"`
	Geometry    struct {
		Type        string        `json:"type"`
		Coordinates []interface{} `json:"coordinates"`
	} `json:"geometry"`
}

func (feat *Feature) ToJSON() ([]byte, error) {

	empty, err := feat.IsEmpty()
	if err != nil {
		return nil, err
	}
	if empty {
		return nil, errors.New("The feature is empty, with nothing to encode into GeoJSON.")
	}

	var construct = &geoJSONEncodeType{
		ID:          feat.ID,
		Type:        "Feature",
		Properties:  feat.Properties,
	}
	construct.Geometry.Type = feat.Type

	switch {
	case feat.Type == "Point":
		coords, err := feat.Geometry.Coords()
		if err != nil {
			return nil, errors.Wrap(err, "could not get geometry coords")
		}
		construct.Geometry.Coordinates = encodeCoord(coords[0])
	case feat.Type == "LineString":
		construct.Geometry.Coordinates, err = encodeLineString(feat.Geometry)
	case feat.Type == "Polygon":
		construct.Geometry.Coordinates, err = encodePolygon(feat.Geometry)
	case feat.Type == "MultiPolygon":
		construct.Geometry.Coordinates, err = encodeMultiPolygon(feat.Geometry)
	default:
		return nil, errors.Newf("Currently GeoJSON must be type Point, Linestring, Polygon, or Multipolygon. Found %s.", feat.Type)
	}
	if err != nil {
		return nil, err
	}

	geojson, err := json.Marshal(construct)
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal geojson")
	}

	return geojson, nil

}

func encodeCoord(coord geos.Coord) []interface{} {
	return []interface{}{coord.X, coord.Y}
}

func encodeLineString(geometry *geos.Geometry) ([]interface{}, error) {

	coords, err := geometry.Coords()
	if err != nil {
		return nil, err
	}

	var res []interface{}
	for i := range coords {
		res = append(res, encodeCoord(coords[i]))
	}

	return res, nil
}

func encodePolygon(geometry *geos.Geometry) ([]interface{}, error) {

	var geometries []*geos.Geometry
	shell, err := geometry.Shell()
	if err != nil {
		return nil, errors.Wrap(err, "could not get shell")
	}
	geometries = append(geometries, shell)

	holes, err := geometry.Holes()
	if err != nil {
		return nil, errors.Wrap(err, "could not get geometry holes")
	}
	geometries = append(geometries, holes...)

	res := []interface{}{}
	for i := range geometries {
		coordProgression, err := encodeLineString(geometries[i])
		if err != nil {
			return nil, err
		}
		res = append(res, coordProgression)
	}

	return res, nil
}

func encodeMultiPolygon(multipolygon *geos.Geometry) ([]interface{}, error) {

	collectionLength, err := multipolygon.NGeometry()
	if err != nil {
		return nil, errors.Wrap(err, "could not create new geometry")
	}
	res := []interface{}{}
	for i := 0; i < collectionLength; i++ {
		g, err := multipolygon.Geometry(i)
		if err != nil {
			return nil, errors.Wrap(err, "could not get polygon geometry")
		}
		compiled, err := encodePolygon(g)
		if err != nil {
			return nil, err
		}
		res = append(res, compiled)
	}

	return res, nil
}
