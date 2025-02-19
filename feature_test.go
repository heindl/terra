package terra

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestFeatureCollection(t *testing.T) {

	t.Parallel()

	Convey("should provide a list of waypoints", t, func() {

		aru := []byte(`
			{
				"type": "FeatureCollection",
				"features": [
					{
						"type": "Feature",
						"properties": {
							"scalerank": 3,
							"featurecla": "Admin-0 country",
							"labelrank": 5,
							"sovereignt": "Netherlands",
							"sov_a3": "NL1",
							"adm0_dif": 1,
							"level": 2,
							"type": "Country",
							"admin": "Aruba",
							"adm0_a3": "ABW",
							"geou_dif": 0,
							"geounit": "Aruba",
							"gu_a3": "ABW",
							"su_dif": 0,
							"subunit": "Aruba",
							"su_a3": "ABW",
							"brk_diff": 0,
							"name": "Aruba",
							"name_long": "Aruba",
							"brk_a3": "ABW",
							"brk_name": "Aruba",
							"brk_group": null,
							"abbrev": "Aruba",
							"postal": "AW",
							"formal_en": "Aruba",
							"formal_fr": null,
							"note_adm0": "Neth.",
							"note_brk": null,
							"name_sort": "Aruba",
							"name_alt": null,
							"mapcolor7": 4,
							"mapcolor8": 2,
							"mapcolor9": 2,
							"mapcolor13": 9,
							"pop_est": 103065,
							"gdp_md_est": 2258,
							"pop_year": -99,
							"lastcensus": 2010,
							"gdp_year": -99,
							"economy": "6. Developing region",
							"income_grp": "2. High income: nonOECD",
							"wikipedia": -99,
							"fips_10": null,
							"iso_a2": "AW",
							"iso_a3": "ABW",
							"iso_n3": "533",
							"un_a3": "533",
							"wb_a2": "AW",
							"wb_a3": "ABW",
							"woe_id": -99,
							"adm0_a3_is": "ABW",
							"adm0_a3_us": "ABW",
							"adm0_a3_un": -99,
							"adm0_a3_wb": -99,
							"continent": "North America",
							"region_un": "Americas",
							"subregion": "Caribbean",
							"region_wb": "Latin America & Caribbean",
							"name_len": 5,
							"long_len": 5,
							"abbrev_len": 5,
							"tiny": 4,
							"homepart": -99
						},
						"geometry": { "type": "Polygon", "coordinates": [ [ [ -69.899121093749997, 12.452001953124991 ], [ -69.895703125, 12.422998046874994 ], [ -69.942187499999989, 12.438525390624989 ], [ -70.004150390625, 12.50048828125 ], [ -70.066113281249997, 12.546972656249991 ], [ -70.050878906249991, 12.597070312499994 ], [ -70.035107421874997, 12.614111328124991 ], [ -69.97314453125, 12.567626953125 ], [ -69.911816406249997, 12.48046875 ], [ -69.899121093749997, 12.452001953124991 ] ] ] }
					}
				]
			}
		`)

		features, err := NewFeatureCollectionFromJSON(aru)
		So(err, ShouldBeNil)

		pf, err := NewPoint(12.5362871,-70.0133061)
		So(err, ShouldBeNil)

		within, err := pf.Within(features[0])
		So(err, ShouldBeNil)

		So(within, ShouldBeTrue)

	})
}

