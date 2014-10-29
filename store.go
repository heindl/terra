package terra

import (
	"fmt"
	"path"

	"github.com/dhconnelly/rtreego"
	"github.com/syndtr/goleveldb/leveldb"
)

var cache *leveldb.DB
var tree *rtreego.Rtree

// Geostore represents ...
type Geostore struct {
	directory string
}

// OpenGeostore creates ...
func Open(directory string) (store *Geostore, er error) {

	store = &Geostore{}

	if directory != "" {
		store.directory = path.Join(directory, "features")
	} else {
		store.directory = path.Join(".", "terrastore", "features")
	}

	if cache, er = leveldb.OpenFile(store.directory, nil); er != nil {
		log.Error("Unable to open the directory: %s. %s", store.directory, er)
	}

	tree = rtreego.NewTree(2, 25, 50)

	iter := cache.NewIterator(nil, nil)
	var feature *Feature
	for iter.Next() {
		if feature, er = NewFeatureFromJSON(iter.Value()); er != nil {
			log.Error("While creating a new feature from initial store data: %s.", er.Error())
			return
		}
		tree.Insert(feature)
	}
	iter.Release()
	if iter.Error() != nil {
		er = fmt.Errorf("From the initial store data iterator: %s.", iter.Error())
		log.Error(er.Error())
	}
	return
}

func (g *Geostore) Close() {
	cache.Close()
}

// Add includes a new Geometry element into the geostore.
// TODO: Open a new cache based on the country, for faster parsing.
func (g *Geostore) Add(features ...*Feature) (keys []string, er error) {

	for i := range features {

		if features[i].IsEmpty() {
			continue
		}

		var value []byte
		if value, er = features[i].ToJSON(); er != nil {
			return
		}

		if er = cache.Put([]byte(features[i].ID), value, nil); er != nil {
			log.Error(er.Error())
			return
		}

		tree.Insert(features[i])

		keys = append(keys, features[i].ID)

	}

	return
}

// Update ...
func (g *Geostore) Update(key []byte, feature *Feature) (er error) {

	var value []byte
	if value, er = feature.ToJSON(); er != nil {
		return
	}

	feat, _ := g.Get(key)
	if !feat.IsEmpty() {
		deleted := tree.Delete(feat)
		if !deleted {
			er = fmt.Errorf("Feature %s not found in rtree, and cannot be deleted.")
			log.Error(er.Error())
			return
		}
	}

	if er = cache.Put(key, value, nil); er != nil {
		log.Error(er.Error())
		return
	}
	tree.Insert(feature)

	return
}

// Remove ...
func (g *Geostore) Remove(key []byte) (er error) {

	feat, _ := g.Get(key)
	if !feat.IsEmpty() {
		deleted := tree.Delete(feat)
		if !deleted {
			er = fmt.Errorf("Feature %s not found in rtree, and cannot be deleted.")
			log.Error(er.Error())
			return
		}
	}

	if er = cache.Delete(key, nil); er != nil {
		log.Error(er.Error())
	}

	return

}

// Get ...
func (g *Geostore) Get(key []byte) (feature *Feature, er error) {
	var response []byte
	if response, er = cache.Get(key, nil); er != nil {
		log.Error(er.Error())
		return
	}
	feature, er = NewFeatureFromJSON(response)
	return
}

func (g *Geostore) Contains(feat *Feature) (list []*Feature, er error) {
	response := tree.SearchIntersect(feat.Bounds())
	for i := range response {
		var ok bool
		var f *Feature
		// if f, ok = response[i].(Feature); !ok {
		// 	log.Warning("Response from rtree could not be typed as a feature.")
		// 	continue
		// }
		f, _ = response[i].(*Feature)
		if ok, er = f.Contains(feat); er != nil {
			return
		}
		if !ok {
			continue
		}
		list = append(list, f)
	}
	return
}

// Contains ...
// func (g *Geostore) Contains(feat *Feature) (list []*Feature, er error) {
//
// 	x, y, er := feat.PointCoords()
// 	if er != nil {
// 		log.Error(er.Error())
// 		return
// 	}
//
// 	// log.Info("Searching geostore for features containing point %v (%f,%f) on date %v.", feat.Property("recordTitle"), x, y, feat.Property("recordDate"))
//
// 	iter := cache.NewIterator(nil, nil)
//
// 	var contained bool
//
// 	for iter.Next() {
// 		var item *Feature
// 		if item, er = NewFeatureFromJSON(iter.Value()); er != nil {
// 			return
// 		}
// 		if contained, er = item.Contains(feat); er != nil {
// 			return
// 		}
// 		if !contained {
// 			continue
// 		}
// 		list = append(list, item)
// 	}
// 	iter.Release()
// 	if iter.Error() != nil {
// 		er = iter.Error()
// 		log.Error(er.Error())
// 		return
// 	}
// 	if len(list) == 0 {
// 		log.Warning("Zero geostore features contained point %v (%f,%f) on date %v.", feat.Property("recordTitle"), x, y, feat.Property("recordDate"))
// 	}
//
// 	return
// }

func (g *Geostore) Clear() (er error) {
	batch := new(leveldb.Batch)
	iter := cache.NewIterator(nil, nil)
	for iter.Next() {
		batch.Delete(iter.Key())
	}
	iter.Release()
	if iter.Error() != nil {
		er = iter.Error()
		return
	}
	er = cache.Write(batch, nil)
	tree = rtreego.NewTree(2, 25, 50)
	return
}

func (g *Geostore) Length() (count int, er error) {
	iter := cache.NewIterator(nil, nil)
	for iter.Next() {
		count = count + 1
	}
	iter.Release()
	if iter.Error() != nil {
		er = iter.Error()
	}
	treeSize := tree.Size()
	if treeSize != count {
		er = fmt.Errorf("Expected both the store and the rtree to have length %d.", count)
	}
	return
}
