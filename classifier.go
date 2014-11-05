package terra

import (
	"errors"
	"fmt"
	"time"

	"github.com/cpucycle/astrotime"
)

type Classifier struct {
	StartDate string  `json:"startdate"`
	EndDate   string  `json:"enddate"`
	Key       string  `json:"key"`
	Value     float64 `json:"value"`
	Type      string  `json:"type"`
}

var ErrClassifierNotFound = errors.New("classifier is not present in feature")

func (feat *Feature) SetClassifier(class *Classifier) error {

	response, er := feat.getClassifier(class)
	if er == ErrClassifierNotFound {
		feat.Classifiers = append(feat.Classifiers, class)
		return nil
	}
	if er != nil {
		return er
	}

	response.Value = class.Value
	return nil

}

func (feat *Feature) HasClassifier(class *Classifier) (exists bool) {

	_, er := feat.getClassifier(class)
	if er != nil {
		return
	}
	exists = true
	return
}

func (feat *Feature) GetClassifier(class *Classifier) (er error) {

	response, er := feat.getClassifier(class)
	if er != nil && er == ErrClassifierNotFound {
		return
	}
	class.Value = response.Value
	return
}

func (feat *Feature) getClassifier(class *Classifier) (response *Classifier, er error) {

	if feat.IsEmpty() {
		er = fmt.Errorf("Asked to find a classifier (%s, %s) for an empty feature.", class.Key, class.StartDate)
		log.Warning(er.Error())
		return
	}

	// FIXME: Should possible only accept go timestamps.
	class.StartDate = class.StartDate[:10]
	if class.EndDate == "" {
		class.EndDate = class.StartDate
	}
	class.EndDate = class.EndDate[:10]

	for _, value := range feat.Classifiers {
		if class.Key != value.Key {
			continue
		}
		if class.StartDate < value.StartDate {
			continue
		}
		if class.EndDate > value.EndDate {
			continue
		}
		response = value
		return
	}

	er = ErrClassifierNotFound
	return
}

// GetClassifiers
// Requires that all classifiers are found to return without error.
func (g *Geostore) GetClassifiers(feat *Feature, filters []*Classifier) (er error) {

	x, y, er := feat.PointCoords()
	if er != nil {
		log.Error(er.Error())
		return
	}

	// fmt.Printf("Getting %d classifiers for point %v (%f,%f) on date %v.", len(filters), feat.Property("recordTitle"), x, y, feat.Property("recordDate"))

	list, er := g.Contains(feat)
	if er != nil {
		log.Error(er.Error())
		return
	}

	for i := range filters {

		if filters[i].Key == "LONGITUDE" {
			filters[i].Value = x
			continue
		}

		if filters[i].Key == "LATITUDE" {
			filters[i].Value = y
			continue
		}

		// FIXME: TEST THIS
		if filters[i].Key == "DAY_LENGTH" {
			var date time.Time
			if date, er = time.Parse("2006-01-02", filters[i].StartDate[:10]); er != nil {
				log.Error(er.Error())
				return
			}
			sunrise := astrotime.CalcSunrise(date, x, y)
			sunset := astrotime.CalcSunset(date, x, y)

			filters[i].Value = sunset.Sub(sunrise).Seconds()
			if filters[i].Value == 0 {
				er = fmt.Errorf("Length of day set to zero for %v (%f,%f) on date %v.", feat.Property("recordTitle"), x, y, feat.Property("recordDate"))
				log.Error(er.Error())
				return
			}
			continue
		}

		if len(list) > 0 {
			for k := range list {
				if er = list[k].GetClassifier(filters[i]); er != ErrClassifierNotFound {
					break
				}
			}
		} else {
			er = ErrClassifierNotFound
		}

		if er != nil && er != ErrClassifierNotFound {
			log.Error("For point %v (%f,%f) on date %v: %s", er.Error(), feat.Property("recordTitle"), x, y, feat.Property("recordDate"))
			return
		}

		// Assigns a non-existant bool to classifier.
		if filters[i].Type == "BOOL" && er == ErrClassifierNotFound {
			er = nil
			filters[i].Value = 0
		}

		if filters[i].Type != "BOOL" && er == ErrClassifierNotFound {
			log.Warning("%s classifier not found for point %v (%f,%f) on date %v.", filters[i].Key, feat.Property("recordTitle"), x, y, feat.Property("recordDate"))
			return
		}

	}

	// log.Info("%d classifiers returning for point. All classifiers returned: %t.", len(filters), er != ErrClassifierNotFound)

	return
}
