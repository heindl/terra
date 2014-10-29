package terra

import (
	"testing"
	"time"
)

func TestGetClassifier(t *testing.T) {

	coordinates := make([][][]float64, 1)
	coordinates[0] = append(coordinates[0], []float64{-80.283590, 44.184382})
	coordinates[0] = append(coordinates[0], []float64{-80.283689, 44.275583})
	coordinates[0] = append(coordinates[0], []float64{-80.156311, 44.275583})
	coordinates[0] = append(coordinates[0], []float64{-80.156410, 44.184382})
	coordinates[0] = append(coordinates[0], []float64{-80.283590, 44.184382})

	polygon, er := NewPolygon(coordinates)
	catchError(er, t)

	er = polygon.SetClassifier(&Classifier{
		StartDate: time.Now().Format("2006-01-02"),
		Key:       "TMAX",
		Value:     98,
	})
	catchError(er, t)

	er = polygon.SetClassifier(&Classifier{
		StartDate: time.Now().Format("2006-01-02"),
		Key:       "TMIN",
		Value:     50,
	})
	catchError(er, t)

	if len(polygon.Classifiers) != 2 {
		t.Errorf("Expected polygon to have two classifiers.")
	}

	c := &Classifier{
		StartDate: time.Now().Format("2006-01-02"),
		Key:       "TMIN",
	}
	er = polygon.GetClassifier(c)
	catchError(er, t)

	if c.Value != 50 {
		t.Errorf("feature.GetClassifier did not return the correct value")
	}

	return

}

func TestGetStoreClassifiers(t *testing.T) {

	var er error

	store, er := Open("./geostore")
	catchError(er, t)
	defer store.Close()
	store.Clear()

	length, er := store.Length()
	catchError(er, t)
	if length != 0 {
		t.Errorf("Expected geostore to be empty.")
	}

	point, er := NewPoint(-80.22, 44.23)
	catchError(er, t)

	polygon, er := NewFeatureFromJSON([]byte(`{
        "type": "Feature",
        "geometry": {
            "type": "Polygon",
            "coordinates": [
                [
                    [-80.283590,44.184382], [-80.283689,44.275583], [-80.156311,44.275583], [-80.156410,44.184382], [-80.283590, 44.184382]
                ]
            ]
        },
        "classifiers": [
            {
                "key": "TMAX",
                "value": 98,
                "startdate": "2014-01-01",
                "enddate": "2014-01-01"
            },
            {
                "key": "TMIN",
                "value": 50,
                "startdate": "2014-01-01",
                "enddate": "2014-01-01"
            }
        ]
    }`))

	contains, er := store.Contains(point)
	catchError(er, t)

	if len(contains) != 0 {
		t.Errorf("Expected store to not contain point.")
	}

	additions, er := store.Add(polygon)
	catchError(er, t)

	if len(additions) != 1 {
		t.Errorf("Expected 1 feature to be added.")
	}

	contains, er = store.Contains(point)
	catchError(er, t)

	if len(contains) != 1 {
		t.Errorf("Expected store to contain one value.")
	}

	var classes []*Classifier
	classes = append(classes, &Classifier{
		Key:       "TMAX",
		StartDate: "2014-01-01",
	})
	classes = append(classes, &Classifier{
		Key:       "TMIN",
		StartDate: "2014-01-01",
	})

	er = store.GetClassifiers(point, classes)
	catchError(er, t)

	if classes[0].Value != 98 {
		t.Errorf("feature.GetClassifier did not return the correct value")
	}

	if classes[1].Value != 50 {
		t.Errorf("feature.GetClassifier did not return the correct value")
	}

	return

}
