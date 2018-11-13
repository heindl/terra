package terra

import (
	"fmt"

	"github.com/dhconnelly/rtreego"
	"github.com/paulsmith/gogeos/geos"
	"github.com/saleswise/errors/errors"
)

type Feature struct {
	ID          string
	Type        string
	Properties  map[string]interface{}
	Geometry    *geos.Geometry
}

type FeatureCollection []*Feature

// FIXME: READ ABOUT WKB AND WKT

func NewFeature() (feature *Feature) {

	feature = &Feature{}

	// AddUsage Random ID string if non-existant.
	feature.ID = generateKey()

	return

}

func NewPoint(latitude float64, longitude float64) (*Feature, error) {
	// Return a new Feature that is a point.
	feat := NewFeature()
	point, err := geos.NewPoint(geos.NewCoord(longitude, latitude))
	if err != nil {
		return nil, errors.Wrap(err, "could not create point")
	}
	if err := feat.SetGeometry("Point", point); err != nil {
		return nil, errors.Wrap(err, "could not set geometry")
	}

	return feat, nil
}

func NewPolygon(polygons [][][]float64) (*Feature, error) {
	// Return a new Feature that is a point.

	feat := NewFeature()

	var iPolygons []interface{}
	for _, linestrings := range polygons {
		var iLinestrings []interface{}
		for _, coordinates := range linestrings {
			var iCoordinates []interface{}
			for i := range coordinates {
				iCoordinates = append(iCoordinates, interface{}(coordinates[i]))
			}
			iLinestrings = append(iLinestrings, iCoordinates)
		}
		iPolygons = append(iPolygons, iLinestrings)
	}

	polygon, err := decodePolygon(iPolygons)
	if err != nil {
		return nil, err
	}
	if err := feat.SetGeometry("Polygon", polygon); err != nil {
		return nil, err
	}

	return feat, nil
}

func (feat Feature) Bounds() *rtreego.Rect {

	typer, err := feat.Geometry.Type()
	if err != nil {
		// return nil, errors.Wrap(err, "could not get type")
		return nil
	}
	if typer == geos.POINT {
		x, err := feat.Geometry.X()
		if err != nil {
			// return nil, errors.Wrap(err, "could not get geometry x")
			return nil
		}
		y, err := feat.Geometry.Y()
		if err != nil {
			// return nil, errors.Wrap(err, "could not get geometry y")
			return nil
		}
		rect, err := rtreego.NewRect(rtreego.Point{x, y}, []float64{0.00001, 0.00001})
		if err != nil {
			// return nil, errors.Wrap(err, "could not get new rectangle")
			return nil
		}
		return rect
		// return rect, nil
	}

	envelope, err := feat.Geometry.Envelope()
	if err != nil {
		// return nil, errors.Wrap(err, "could not get envelope")
		return nil
	}
	envelope, err = envelope.Shell()
	if err != nil {
		// return nil, errors.Wrap(err, "could not get shell")
		return nil
	}
	coords, er := envelope.Coords()
	if er != nil {
		// return nil, errors.Wrap(err, "could not get coords")
		return nil
	}
	if len(coords) != 5 {
		// return nil, errors.New("coords not equal to five")
		return nil
	}

	// FIXME: GEOJSON somehow reverses these coordinates.
	height := coords[1].X - coords[0].X
	width := coords[3].Y - coords[0].Y
	rect, err := rtreego.NewRect(rtreego.Point{coords[0].X, coords[0].Y}, []float64{height, width})
	if err != nil {
		// return nil, errors.Wrap(err, "could not get rectangle")
		return nil
	}

	return rect
}

func (feat *Feature) SetGeometry(typer string, geometry *geos.Geometry) error {

	if typer != "Polygon" && typer != "Point" && typer != "MultiPolygon" {
		return errors.New("Presently, geostore only accepts GeoJSON types Point, LineString, Polygon and Multipolygon.")
	}

	feat.Type = typer

	if geometry == nil {
		return nil
	}

	feat.Geometry = geometry

	return nil

}

func (feat *Feature) Property(name string) (property interface{}) {

	//FIXME: Properly handle subarrays or subobjects.

	property, ok := feat.Properties[name]
	if !ok {
		return
	}

	return

	/*
		if floatValue, ok := feat.Properties[name].(float64); ok {
			property = strconv.FormatFloat(floatValue, 'f', 6, 64)
			return
		}

		if boolValue, ok := feat.Properties[name].(bool); ok {
			// Fixme: More Tersley Convert Boolean to string
			property = strconv.FormatBool(boolValue)
			return
		}

		if property, ok = feat.Properties[name].(string); ok {
			// Fixme: More Tersley Convert Boolean to string
			return
		}
	*/

	return
}

func (feat *Feature) SetProperty(name string, property interface{}) {
	if feat.Properties == nil {
		feat.Properties = make(map[string]interface{})
	}
	feat.Properties[name] = property
}

func (feat *Feature) Contains(subfeat *Feature) (bool, error) {

	contains, err := feat.Geometry.Contains(subfeat.Geometry)
	if err != nil {
		return false, errors.Wrap(err, "could not check geometry")
	}
	return contains, nil
}

func (feat *Feature) Within(subfeat *Feature) (bool, error) {
	within, err := feat.Geometry.Within(subfeat.Geometry);
	if err != nil {
		return false, errors.Wrap(err, "could ot check subfeature within")
	}
	return within, nil
}

func (feat *Feature) PointCoords() (x float64, y float64, er error) {

	if feat.Type != "Point" {
		er = fmt.Errorf("The feature geometry must be a point to return an x,y coordinate")
	}

	if x, er = feat.Geometry.X(); er != nil {
		return 0, 0, errors.Wrap(er, "could not get geometry x")
	}

	if y, er = feat.Geometry.Y(); er != nil {
		return 0, 0, errors.Wrap(er, "could not get geometry y")
	}

	return
}

func (feat *Feature) IsEmpty() (bool, error) {

	if feat == nil {
		return true, nil
	}

	if feat.Geometry == nil {
		return true, nil
	}

	empty, er := feat.Geometry.IsEmpty()
	if er != nil {
		return false, errors.Wrap(er, "could not get geometry")
	}

	if empty {
		return true, nil
	}

	if feat.ID == "" {
		return true, nil
	}

	return false, nil
}
