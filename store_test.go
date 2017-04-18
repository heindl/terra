package terra

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestStore(t *testing.T) {

	Convey("should create a new polygon", t, func() {

		// Broadly represents North America.
		polygon, err := NewPolygon([][][]float64{
			[][]float64{
				{-139.75, 55.03},
				{-51.28, 55.03},
				{-51.28, 23.73},
				{-139.75, 23.73},
				{-139.75, 55.03},
			},
		})
		So(err, ShouldBeNil)

		point, err := NewPoint(33.7489954, -84.3879824)
		So(err, ShouldBeNil)

		within, err := polygon.Contains(point)
		So(err, ShouldBeNil)
		So(within, ShouldBeTrue)
	})

	Convey("given an environment", t, func() {

		store, err := Open("./geostore")
		So(err, ShouldBeNil)

		Convey("should save and renew from store", func() {

			length, err := store.Length()
			So(err, ShouldBeNil)
			So(length, ShouldEqual, 0)

			point, err := NewPoint(18.2324388, -63.0419003)
			So(err, ShouldBeNil)

			distantPolygon, err := NewFeatureFromJSON([]byte(`{ "type": "Feature", "properties": { "scalerank": 3, "featurecla": "Admin-0 country", "labelrank": 5, "sovereignt": "Netherlands", "sov_a3": "NL1", "adm0_dif": 1, "level": 2, "type": "Country", "admin": "Aruba", "adm0_a3": "ABW", "geou_dif": 0, "geounit": "Aruba", "gu_a3": "ABW", "su_dif": 0, "subunit": "Aruba", "su_a3": "ABW", "brk_diff": 0, "name": "Aruba", "name_long": "Aruba", "brk_a3": "ABW", "brk_name": "Aruba", "brk_group": null, "abbrev": "Aruba", "postal": "AW", "formal_en": "Aruba", "formal_fr": null, "note_adm0": "Neth.", "note_brk": null, "name_sort": "Aruba", "name_alt": null, "mapcolor7": 4, "mapcolor8": 2, "mapcolor9": 2, "mapcolor13": 9, "pop_est": 103065, "gdp_md_est": 2258, "pop_year": -99, "lastcensus": 2010, "gdp_year": -99, "economy": "6. Developing region", "income_grp": "2. High income: nonOECD", "wikipedia": -99, "fips_10": null, "iso_a2": "AW", "iso_a3": "ABW", "iso_n3": "533", "un_a3": "533", "wb_a2": "AW", "wb_a3": "ABW", "woe_id": -99, "adm0_a3_is": "ABW", "adm0_a3_us": "ABW", "adm0_a3_un": -99, "adm0_a3_wb": -99, "continent": "North America", "region_un": "Americas", "subregion": "Caribbean", "region_wb": "Latin America & Caribbean", "name_len": 5, "long_len": 5, "abbrev_len": 5, "tiny": 4, "homepart": -99 }, "geometry": { "type": "Polygon", "coordinates": [ [ [ -69.899121093749997, 12.452001953124991 ], [ -69.895703125, 12.422998046874994 ], [ -69.942187499999989, 12.438525390624989 ], [ -70.004150390625, 12.50048828125 ], [ -70.066113281249997, 12.546972656249991 ], [ -70.050878906249991, 12.597070312499994 ], [ -70.035107421874997, 12.614111328124991 ], [ -69.97314453125, 12.567626953125 ], [ -69.911816406249997, 12.48046875 ], [ -69.899121093749997, 12.452001953124991 ] ] ] } }`))
			So(err, ShouldBeNil)

			additions, err := store.Add(distantPolygon)
			So(err, ShouldBeNil)
			So(len(additions), ShouldEqual, 1)

			length, err = store.Length()
			So(err, ShouldBeNil)
			So(length, ShouldEqual, 1)
			//if length != 1 {
			//	t.Errorf("Expected geostore to have two elements.")
			//}

			contains, err := store.Contains(point)
			So(err, ShouldBeNil)
			So(len(contains), ShouldEqual, 0)

			polygon, err := NewFeatureFromJSON([]byte(`{ "type": "Feature", "properties": { "scalerank": 1, "featurecla": "Admin-0 country", "labelrank": 6, "sovereignt": "United Kingdom", "sov_a3": "GB1", "adm0_dif": 1, "level": 2, "type": "Dependency", "admin": "Anguilla", "adm0_a3": "AIA", "geou_dif": 0, "geounit": "Anguilla", "gu_a3": "AIA", "su_dif": 0, "subunit": "Anguilla", "su_a3": "AIA", "brk_diff": 0, "name": "Anguilla", "name_long": "Anguilla", "brk_a3": "AIA", "brk_name": "Anguilla", "brk_group": null, "abbrev": "Ang.", "postal": "AI", "formal_en": null, "formal_fr": null, "note_adm0": "U.K.", "note_brk": null, "name_sort": "Anguilla", "name_alt": null, "mapcolor7": 6, "mapcolor8": 6, "mapcolor9": 6, "mapcolor13": 3, "pop_est": 14436, "gdp_md_est": 108.90000000000001, "pop_year": -99, "lastcensus": -99, "gdp_year": -99, "economy": "6. Developing region", "income_grp": "3. Upper middle income", "wikipedia": -99, "fips_10": null, "iso_a2": "AI", "iso_a3": "AIA", "iso_n3": "660", "un_a3": "660", "wb_a2": "-99", "wb_a3": "-99", "woe_id": -99, "adm0_a3_is": "AIA", "adm0_a3_us": "AIA", "adm0_a3_un": -99, "adm0_a3_wb": -99, "continent": "North America", "region_un": "Americas", "subregion": "Caribbean", "region_wb": "Latin America & Caribbean", "name_len": 8, "long_len": 8, "abbrev_len": 4, "tiny": -99, "homepart": -99 }, "geometry": { "type": "Polygon", "coordinates": [ [ [ -63.001220703125, 18.221777343749991 ], [ -63.160009765624991, 18.17138671875 ], [ -63.1533203125, 18.200292968749991 ], [ -63.026025390624994, 18.269726562499997 ], [ -62.979589843749991, 18.264794921874994 ], [ -63.001220703125, 18.221777343749991 ] ] ] } }`))
			So(err, ShouldBeNil)

			additions, err = store.Add(polygon)
			So(err, ShouldBeNil)
			So(len(additions), ShouldEqual, 1)

			length, err = store.Length()
			So(err, ShouldBeNil)
			So(length, ShouldEqual, 2)

			contains, err = store.Contains(point)
			So(err, ShouldBeNil)
			So(len(contains), ShouldEqual, 1)

		})

		Convey("should save polygon to store", func() {

			length, err := store.Length()
			So(err, ShouldBeNil)
			So(length, ShouldEqual, 0)

			point, err := NewPoint(37.865101, -119.538329)
			So(err, ShouldBeNil)

			polygon, err := NewFeatureFromJSON([]byte(`{
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
			So(err, ShouldBeNil)

			within, err := polygon.Contains(point)
			So(err, ShouldBeNil)
			So(within, ShouldBeTrue)

			additions, err := store.Add(polygon)
			So(err, ShouldBeNil)
			So(len(additions),ShouldEqual, 1)

			contains, err := store.Contains(point)
			So(err, ShouldBeNil)
			So(len(contains), ShouldEqual, 1)
		})



		Reset(func() {
			So(store.Clear(), ShouldBeNil)
			So(store.Close(), ShouldBeNil)
		})

	})

}
