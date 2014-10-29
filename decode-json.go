package terra

import (
	"crypto/rand"
	"encoding/json"
	"fmt"

	"github.com/paulsmith/gogeos/geos"
)

func NewFeatureFromJSON(request []byte) (feature *Feature, er error) {

	var g map[string]interface{}
	if er = json.Unmarshal(request, &g); er != nil {
		log.Error(er.Error())
		return
	}
	feature, er = decodeFeature(g)
	return

}

func NewFeatureCollectionFromJSON(request []byte) (features FeatureCollection, er error) {

	type geoJSONFeatureType struct {
		Features []map[string]interface{} `json:"features"`
	}

	var geo geoJSONFeatureType
	if er = json.Unmarshal(request, &geo); er != nil {
		log.Error(er.Error())
		return
	}

	for i := range geo.Features {
		var feat *Feature
		if feat, er = decodeFeature(geo.Features[i]); er != nil {
			return
		}
		features = append(features, feat)
	}

	return
}

func decodeFeature(geo map[string]interface{}) (feature *Feature, er error) {

	feature = &Feature{}
	var ok bool

	// ADD RANDOM ID STRING IF NONEXISTANT
	if _, ok = geo["id"]; !ok || feature.ID == "" {
		feature.ID = generateKey()
	} else {
		feature.ID = geo["id"].(string)
	}

	// DECODE PROPERTIES
	if p, ok := geo["properties"]; ok && p != nil {
		feature.Properties = geo["properties"].(map[string]interface{})
	}

	// DECODE CLASSIFIERS
	if _, ok = geo["classifiers"]; ok {
		if classifiers, ok := geo["classifiers"].([]interface{}); ok {
			for i := range classifiers {
				if class, ok := classifiers[i].(map[string]interface{}); ok {

					c := &Classifier{}

					if _, ok = class["key"]; ok {
						c.Key = class["key"].(string)
					}

					if _, ok = class["value"]; ok {
						c.Value = class["value"].(float64)
					}

					if _, ok = class["startdate"]; ok {
						c.StartDate = class["startdate"].(string)
					}

					if _, ok = class["enddate"]; ok {
						c.EndDate = class["enddate"].(string)
					}

					if _, ok = class["type"]; ok {
						c.Type = class["type"].(string)
					}

					feature.Classifiers = append(feature.Classifiers, c)

				}
			}
		}
	}

	if _, ok = geo["geometry"]; !ok {
		er = fmt.Errorf("Missing a geoJSON geometry property.")
		log.Error(er.Error())
		return
	}

	var geometry map[string]interface{}
	if geometry, ok = geo["geometry"].(map[string]interface{}); !ok {
		er = fmt.Errorf("Geometry property is malformed: %s.", geo["geometry"])
		log.Error(er.Error())
		return
	}

	if _, ok = geometry["type"]; !ok {
		er = fmt.Errorf("A Geometry Type property is required for decoding geoJSON.")
		log.Error(er.Error())
		return
	}

	if feature.Type, ok = geometry["type"].(string); !ok {
		er = fmt.Errorf("The geoJSON Geometry Type property is expected to be a string.")
		log.Error(er.Error())
		return
	}

	if _, ok = geometry["coordinates"]; !ok {
		er = fmt.Errorf("GeoJSON Geometry Coordinates property is required.")
		log.Error(er.Error())
		return
	}

	var coordinates []interface{}
	if coordinates, ok = geometry["coordinates"].([]interface{}); !ok {
		er = fmt.Errorf("Geometry Coordinates property values are are malformed: %s.", geometry["coordinates"])
		log.Error(er.Error())
		return
	}

	switch {
	case feature.Type == "Point":
		if feature.Geometry, er = decodePoint(coordinates); er != nil {
			return
		}
	case feature.Type == "LineString":
		if feature.Geometry, er = decodeLineString(coordinates); er != nil {
			return
		}
	case feature.Type == "Polygon":
		if feature.Geometry, er = decodePolygon(coordinates); er != nil {
			return
		}
	case feature.Type == "MultiPolygon":
		if feature.Geometry, er = decodeMultiPolygon(coordinates); er != nil {
			return
		}
	default:
		er = fmt.Errorf("Currently, GeoJSON must be type Point, Linestring, Polygon, or Multipolygon. Found %s.", feature.Type)
		log.Error(er.Error())
		return
	}

	return

}

func decodePoint(coordinates []interface{}) (response *geos.Geometry, er error) {

	var latitude, longitude float64
	var ok bool

	if latitude, ok = coordinates[0].(float64); !ok {
		er = fmt.Errorf("First element of Point array should be a float64 %s.", coordinates[0])
		log.Error(er.Error())
		return
	}

	if longitude, ok = coordinates[1].(float64); !ok {
		er = fmt.Errorf("Second element of Point array should be a float64 %s.", coordinates[0])
		log.Error(er.Error())
		return
	}

	if response, er = geos.NewPoint(geos.NewCoord(latitude, longitude)); er != nil {
		log.Error(er.Error())
	}

	return

}

func decodeLineString(coordinates []interface{}) (response *geos.Geometry, er error) {

	var coords []geos.Coord

	for _, coordinate := range coordinates {

		var (
			points    []interface{}
			latitude  float64
			longitude float64
			ok        bool
		)

		if points, ok = coordinate.([]interface{}); !ok {
			er = fmt.Errorf("Expect each element in a LineString property array to be a coordinate array.")
			log.Error(er.Error())
			return
		}

		if latitude, ok = points[0].(float64); !ok {
			er = fmt.Errorf("First element of Point array should be a float64 %s.", points[0])
			log.Error(er.Error())
			return
		}

		if longitude, ok = points[1].(float64); !ok {
			er = fmt.Errorf("Second element of Point array should be a float64 %s.", points[1])
			log.Error(er.Error())
			return
		}

		coords = append(coords, geos.NewCoord(latitude, longitude))
	}

	if response, er = geos.NewLineString(coords...); er != nil {
		log.Error(er.Error())
	}

	return

}

func decodePolygon(coordinates []interface{}) (response *geos.Geometry, er error) {

	var contours [][]geos.Coord
	var ok bool

	for _, coordinate := range coordinates {

		var coords []geos.Coord
		var linestrings []interface{}

		if linestrings, ok = coordinate.([]interface{}); !ok {
			er = fmt.Errorf("Expect each sub-element in a Polygon property array to be a LineString array.")
			log.Error(er.Error())
			return
		}

		for _, linestring := range linestrings {

			var latitude, longitude float64
			var points []interface{}

			if points, ok = linestring.([]interface{}); !ok {
				er = fmt.Errorf("Expect each sub-element in a Polygon Linestring sub-array to be a coordinates array.")
				log.Error(er.Error())
				return
			}

			if latitude, ok = points[0].(float64); !ok {
				er = fmt.Errorf("First element of Point array should be a float64 %s.", points[0])
				log.Error(er.Error())
				return
			}

			if longitude, ok = points[1].(float64); !ok {
				er = fmt.Errorf("Second element of Point array should be a float64 %s.", points[1])
				log.Error(er.Error())
				return
			}

			coords = append(coords, geos.NewCoord(latitude, longitude))
		}

		contours = append(contours, coords)

	}

	if response, er = geos.NewPolygon(contours[0], contours[1:]...); er != nil {
		log.Error(er.Error())
	}

	return

}

func decodeMultiPolygon(coordinates []interface{}) (response *geos.Geometry, er error) {

	var geometries []*geos.Geometry
	var ok bool

	for _, coordinate := range coordinates {
		var polygon []interface{}
		if polygon, ok = coordinate.([]interface{}); !ok {
			er = fmt.Errorf("Error decoding MultiPolygon.")
			log.Error(er.Error())
			return
		}
		var g *geos.Geometry
		if g, er = decodePolygon(polygon); er != nil {
			return
		}
		geometries = append(geometries, g)
	}

	if response, er = geos.NewCollection(geos.MULTIPOLYGON, geometries...); er != nil {
		log.Error(er.Error())
	}

	return
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
