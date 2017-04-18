package terra

import (
	"crypto/rand"
	"encoding/json"
	"github.com/paulsmith/gogeos/geos"
	"github.com/saleswise/errors/errors"
	"bitbucket.org/heindl/logkeys"
	. "bitbucket.org/heindl/malias"
	"strconv"
)

func NewFeatureFromJSON(request []byte) (*Feature, error) {
	var g map[string]interface{}
	if err := json.Unmarshal(request, &g); err != nil {
		return nil, errors.Wrap(err, "could not unmarshal json for feature")
	}
	return decodeFeature(g)
}

func NewFeatureCollectionFromJSON(request []byte) (FeatureCollection, error) {

	type geoJSONFeatureType struct {
		Features []map[string]interface{} `json:"features"`
	}

	var geo geoJSONFeatureType
	if err := json.Unmarshal(request, &geo); err != nil {
		return nil, errors.Wrap(err, "could not unmarhsal geojson feature type")
	}

	var coll FeatureCollection
	for i := range geo.Features {
		feat, err := decodeFeature(geo.Features[i])
		if err != nil {
			return nil, err
		}
		coll = append(coll, feat)
	}

	return coll, nil
}

func decodeFeature(geo map[string]interface{}) (*Feature, error) {
	// ADD RANDOM ID STRING IF NONEXISTANT

	feature := Feature{}

	if _, ok := geo["id"]; !ok || feature.ID == "" {
		feature.ID = generateKey()
	} else {
		feature.ID = geo["id"].(string)
	}

	// DECODE PROPERTIES
	if p, ok := geo["properties"]; ok && p != nil {
		feature.Properties = geo["properties"].(map[string]interface{})
	}

	g, ok := geo["geometry"]
	if !ok {
		return nil, errors.New("Missing a geoJSON geometry property.")
	}

	geometry, ok := g.(map[string]interface{})
	if !ok {
		return nil, errors.Newf("Geometry property is malformed: %s.", geo["geometry"])
	}

	geometryType, ok := geometry["type"]
	if !ok {
		return nil, errors.New("A Geometry Type property is required for decoding geoJSON.")
	}

	feature.Type, ok = geometryType.(string)
	if !ok {
		return nil, errors.New("The geoJSON Geometry Type property is expected to be a string.")
	}

	coords, ok := geometry["coordinates"]
	if !ok {
		return nil, errors.New("GeoJSON Geometry Coordinates property is required.")
	}

	coordinates, ok := coords.([]interface{})
	if !ok {
		return nil, errors.Newf("Geometry Coordinates property values are are malformed: %s.", geometry["coordinates"])
	}

	var err error
	switch {
	case feature.Type == "Point":
		feature.Geometry, err = decodePoint(coordinates)
	case feature.Type == "LineString":
		feature.Geometry, err = decodeLineString(coordinates)
	case feature.Type == "Polygon":
		feature.Geometry, err = decodePolygon(coordinates)
	case feature.Type == "MultiPolygon":
		feature.Geometry, err = decodeMultiPolygon(coordinates)
	default:
		return nil, errors.Newf("Unsupported type: %s.Currently, GeoJSON must be type Point, Linestring, Polygon, or Multipolygon.", feature.Type)
	}
	if err != nil {
		return nil, err
	}

	return &feature, nil

}

func decodePoint(coordinates []interface{}) (*geos.Geometry, error) {

	var (
		latitude, longitude float64
		err error
	)

	if lng, ok := coordinates[0].(string); !ok || lng == "" {
		return nil, errors.New("first element of Point array should be a float64 longitude").SetState(M{logkeys.StringValue: coordinates[0]})
	} else {
		longitude, err = strconv.ParseFloat(lng, 64)
		if err != nil {
			return nil, errors.Wrap(err, "could not parse longitude").SetState(M{logkeys.StringValue: lng})
		}
	}

	if lat, ok := coordinates[1].(string); !ok || lat == "" {
		return nil, errors.New("second element of Point array should be a float64 latitude").SetState(M{logkeys.StringValue: coordinates[1]})
	} else {
		latitude, err = strconv.ParseFloat(lat, 64)
		if err != nil {
			return nil, errors.Wrap(err, "could not parse latitude").SetState(M{logkeys.StringValue: lat})
		}
	}

	response, err := geos.NewPoint(geos.NewCoord(longitude, latitude))
	if err != nil {
		return nil, errors.Wrap(err, "could not create new coordinates")
	}

	return response, nil

}

func decodeLineString(coordinates []interface{}) (*geos.Geometry, error) {

	var coords []geos.Coord

	for _, coordinate := range coordinates {

		var (
			points    []interface{}
			latitude  float64
			longitude float64
			ok        bool
		)

		if points, ok = coordinate.([]interface{}); !ok {
			return nil, errors.New("Expect each element in a LineString property array to be a coordinate array.")
		}

		if latitude, ok = points[1].(float64); !ok {
			return nil, errors.Newf("First element of Point array should be a float64 %s.", points[0])
		}

		if longitude, ok = points[0].(float64); !ok {
			return nil, errors.Newf("Second element of Point array should be a float64 %s.", points[1])
		}

		coords = append(coords, geos.NewCoord(latitude, longitude))
	}

	response, err := geos.NewLineString(coords...)
	if err != nil {
		return nil, errors.Wrap(err, "could not create line string")
	}

	return response, nil

}

func decodePolygon(coordinates []interface{}) (*geos.Geometry, error) {

	var contours [][]geos.Coord

	for _, coordinate := range coordinates {

		linestrings, ok := coordinate.([]interface{})
		if !ok {
			return nil, errors.New("Expect each sub-element in a Polygon property array to be a LineString array.")
		}

		coords := []geos.Coord{}

		for _, linestring := range linestrings {

			points, ok := linestring.([]interface{})
			if !ok {
				return nil, errors.New("Expect each sub-element in a Polygon Linestring sub-array to be a coordinates array.")
			}

			latitude, ok := points[1].(float64)
			if !ok {
				return nil, errors.Newf("First element of Point array should be a float64 %s.", points[0])
			}

			longitude, ok := points[0].(float64)
			if !ok {
				return nil, errors.Newf("Second element of Point array should be a float64 %s.", points[1])
			}
			coords = append(coords, geos.NewCoord(longitude, latitude))
		}

		contours = append(contours, coords)

	}

	response, err := geos.NewPolygon(contours[0], contours[1:]...)
	if err != nil {
		return nil, errors.Wrap(err, "could not create new polygon")
	}

	return response, nil

}

func decodeMultiPolygon(coordinates []interface{}) (*geos.Geometry, error) {

	geometries := []*geos.Geometry{}

	for _, coordinate := range coordinates {
		polygon, ok := coordinate.([]interface{})
		if !ok {
			return nil, errors.New("coordinate polygon not an interface array")
		}
		g, err := decodePolygon(polygon)
		if err != nil {
			return nil, errors.Wrap(err, "could not decode polygon")
		}
		geometries = append(geometries, g)
	}

	response, err := geos.NewCollection(geos.MULTIPOLYGON, geometries...)
	if err != nil {
		return nil, errors.Wrap(err, "could not create new collection")
	}

	return response, nil
}

func generateKey() string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, 15)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}
