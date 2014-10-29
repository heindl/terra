package terra

import "testing"

func catchError(er error, t *testing.T) {
	if er != nil {
		t.Errorf(er.Error())
		panic(er.Error())
	}
}

func TestContains(t *testing.T) {

	var point *Feature
	var er error
	var store *Geostore

	store, er = Open("./geostore")
	catchError(er, t)
	defer store.Close()
	store.Clear()

	length, er := store.Length()
	catchError(er, t)
	if length != 0 {
		t.Errorf("Expected geostore to be empty.")
	}

	point, er = NewPoint(-80.22, 44.23)
	catchError(er, t)

	distantPolygon, er := NewFeatureFromJSON([]byte(`{
		"type": "Feature",
		"geometry": {
			"type": "Polygon",
			"coordinates": [
				[
					[-90.283590,44.184382], [-90.283689,44.275583], [-90.156311,44.275583], [-90.156410,44.184382], [-90.283590, 44.184382]
				]
			]
		}
	}`))

	additions, er := store.Add(distantPolygon)
	catchError(er, t)

	if len(additions) != 1 {
		t.Errorf("Expected 1 feature to be added.")
	}

	length, er = store.Length()
	catchError(er, t)
	if length != 1 {
		t.Errorf("Expected geostore to have two elements.")
	}

	contains, er := store.Contains(point)
	catchError(er, t)

	if len(contains) != 0 {
		t.Errorf("Expected store to not contain the point.")
	}

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

	additions, er = store.Add(polygon)
	catchError(er, t)

	if len(additions) != 1 {
		t.Errorf("Expected 1 feature to be added.")
	}

	length, er = store.Length()
	catchError(er, t)
	if length != 2 {
		t.Errorf("Expected geostore to have two elements.")
	}

	contains, er = store.Contains(point)
	catchError(er, t)

	if len(contains) != 1 {
		t.Errorf("Expected store to contain one value.")
	}

	return

}

func TestNewPolygon(t *testing.T) {

	point, er := NewPoint(-80.22, 44.23)
	catchError(er, t)

	coordinates := make([][][]float64, 1)
	coordinates[0] = append(coordinates[0], []float64{-80.283590, 44.184382})
	coordinates[0] = append(coordinates[0], []float64{-80.283689, 44.275583})
	coordinates[0] = append(coordinates[0], []float64{-80.156311, 44.275583})
	coordinates[0] = append(coordinates[0], []float64{-80.156410, 44.184382})
	coordinates[0] = append(coordinates[0], []float64{-80.283590, 44.184382})

	polygon, er := NewPolygon(coordinates)
	catchError(er, t)

	within, er := polygon.Contains(point)
	catchError(er, t)
	if !within {
		t.Errorf("Expected yosemite to contain point.")
	}

	return

}

func TestYosemiteContains(t *testing.T) {

	var point *Feature
	var er error
	var store *Geostore

	store, er = Open("./geostore")
	catchError(er, t)
	defer store.Close()
	store.Clear()

	length, er := store.Length()
	catchError(er, t)
	if length != 0 {
		t.Errorf("Expected geostore to be empty.")
	}

	point, er = NewPoint(-119.538329, 37.865101)
	catchError(er, t)

	polygon, er := NewFeatureFromJSON([]byte(`{
        "type": "Feature",
        "properties": {
            "unit_code": "YOSE",
            "unit_name":
            "Yosemite NP",
            "unit_type": "National Park",
            "nps_region": "Pacific West",
            "scalerank": 3,
            "featurecla": "National Park Service",
            "note": null,
            "name": "Yosemite"
        },
        "geometry": {
            "type": "Polygon",
            "coordinates": [
                [
                    [ -119.542195638020843, 38.1513671875 ],
                    [ -119.5035400390625, 38.136800130208336 ],
                    [ -119.498738606770843, 38.156168619791671 ],
                    [ -119.469767252604171, 38.127156575520836 ],
                    [ -119.4600830078125, 38.098225911458336 ],
                    [ -119.42626953125, 38.117513020833336 ],
                    [ -119.344197591145843, 38.083699544270836 ],
                    [ -119.310384114583343, 38.045084635416671 ],
                    [ -119.305582682291671, 38.011271158854171 ],
                    [ -119.320027669270843, 37.967814127604171 ],
                    [ -119.310384114583343, 37.948527018229171 ],
                    [ -119.266927083333343, 37.92919921875 ],
                    [ -119.262125651041671, 37.909871419270836 ],
                    [ -119.228312174479171, 37.909871419270836 ],
                    [ -119.199300130208343, 37.885701497395836 ],
                    [ -119.218668619791671, 37.847127278645836 ],
                    [ -119.194498697916671, 37.842244466145836 ],
                    [ -119.218668619791671, 37.818115234375 ],
                    [ -119.199300130208343, 37.798787434895836 ],
                    [ -119.237955729166671, 37.769856770833336 ],
                    [ -119.266927083333343, 37.740885416666671 ],
                    [ -119.257283528645843, 37.707071940104171 ],
                    [ -119.2862548828125, 37.687744140625 ],
                    [ -119.324869791666671, 37.634602864583336 ],
                    [ -119.368367513020843, 37.629801432291671 ],
                    [ -119.3876953125, 37.591145833333336 ],
                    [ -119.3876953125, 37.552530924479171 ],
                    [ -119.4166259765625, 37.557373046875 ],
                    [ -119.445597330729171, 37.538045247395836 ],
                    [ -119.576009114583343, 37.533243815104171 ],
                    [ -119.576009114583343, 37.494588216145836 ],
                    [ -119.658040364583343, 37.499430338541671 ],
                    [ -119.677408854166671, 37.538045247395836 ],
                    [ -119.7015380859375, 37.552530924479171 ],
                    [ -119.696695963541671, 37.629801432291671 ],
                    [ -119.720865885416671, 37.629801432291671 ],
                    [ -119.716023763020843, 37.658772786458336 ],
                    [ -119.759480794270843, 37.658772786458336 ],
                    [ -119.759480794270843, 37.687744140625 ],
                    [ -119.77880859375, 37.707071940104171 ],
                    [ -119.77880859375, 37.740885416666671 ],
                    [ -119.851236979166671, 37.740885416666671 ],
                    [ -119.851236979166671, 37.760172526041671 ],
                    [ -119.8753662109375, 37.789143880208336 ],
                    [ -119.870564778645843, 37.8084716796875 ],
                    [ -119.827107747395843, 37.832600911458336 ],
                    [ -119.827107747395843, 37.890584309895836 ],
                    [ -119.885050455729171, 37.890584309895836 ],
                    [ -119.885050455729171, 37.991984049479171 ],
                    [ -119.86572265625, 38.069254557291671 ],
                    [ -119.8319091796875, 38.093343098958336 ],
                    [ -119.802978515625, 38.088541666666671 ],
                    [ -119.7353515625, 38.098225911458336 ],
                    [ -119.69189453125, 38.131998697916671 ],
                    [ -119.624267578125, 38.1513671875 ],
                    [ -119.5904541015625, 38.185139973958336 ],
                    [ -119.576009114583343, 38.156168619791671 ],
                    [ -119.542195638020843, 38.1513671875 ]
                ]
            ]
        }
    }`))

	within, er := polygon.Contains(point)
	catchError(er, t)
	if !within {
		t.Errorf("Expected yosemite to contain point.")
	}

	additions, er := store.Add(polygon)
	catchError(er, t)

	if len(additions) != 1 {
		t.Errorf("Expected 1 feature to be added.")
	}

	contains, er := store.Contains(point)
	catchError(er, t)

	if len(contains) != 1 {
		t.Errorf("Expected store to contain one value.")
	}

	return
}
