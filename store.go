package terra

import (
	"fmt"
	"path"

	"github.com/dhconnelly/rtreego"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/saleswise/errors/errors"
)



// Geostore represents ...
type Geostore struct {
	directory string
	cache *leveldb.DB
	tree *rtreego.Rtree
}

// OpenGeostore creates ...
func Open(directory string) (*Geostore, error) {

	store := Geostore{}

	if directory != "" {
		store.directory = path.Join(directory, "features")
	} else {
		store.directory = path.Join(".", "terrastore", "features")
	}

	var err error
	store.cache, err = leveldb.OpenFile(store.directory, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to open the directory: %s.", store.directory)
	}

	store.tree = rtreego.NewTree(2, 25, 50)

	iter := store.cache.NewIterator(nil, nil)
	defer iter.Release()
	for iter.Next() {
		feature, err := NewFeatureFromJSON(iter.Value())
		if err != nil {
			return nil, err
		}
		store.tree.Insert(feature)
	}
	if iter.Error() != nil {
		return nil, errors.Wrapf(err, "error iterating through store")
	}
	return &store, nil
}

func (g *Geostore) Close() error {
	//if err := g.cache.Close(); err != nil && err != leveldb.ErrClosed {
	//	return errors.Wrap(err, "could not close geostore")
	//}
	return g.cache.Close()
}

func (g *Geostore) IsClosed() bool {
	return (g == nil || g == &Geostore{})
}

// Add includes a new Geometry element into the geostore.
// TODO: Open a new cache based on the country, for faster parsing.
func (g *Geostore) Add(features ...*Feature) ([]string, error) {

	keys := []string{}
	for i := range features {

		empty, err := features[i].IsEmpty()
		if err != nil {
			return nil, err
		}
		if empty {
			continue
		}

		value, err := features[i].ToJSON()
		if err != nil {
			return nil, err
		}

		if err := g.cache.Put([]byte(features[i].ID), value, nil); err != nil {
			return nil, errors.Wrap(err, "could not put feature in cache")
		}

		g.tree.Insert(features[i])

		keys = append(keys, features[i].ID)

	}

	return keys, nil
}

// Update ...
func (g *Geostore) Update(key []byte, feature *Feature) error {

	value, err := feature.ToJSON()
	if err != nil {
		return err
	}

	feat, _ := g.Get(key)
	empty, err := feat.IsEmpty()
	if err != nil {
		return err
	}
	if !empty && !g.tree.Delete(feat) {
		return errors.Newf("Feature %s not found in rtree, and cannot be deleted.", feat)
	}

	if err := g.cache.Put(key, value, nil); err != nil {
		return errors.Wrap(err, "could not add feature to cache")
	}
	g.tree.Insert(feature)

	return nil
}

// Remove ...
func (g *Geostore) Remove(key []byte) error {

	feat, _ := g.Get(key)
	empty, err := feat.IsEmpty()
	if err != nil {
		return err
	}
	if !empty && !g.tree.Delete(feat) {
		return errors.Newf("Feature %s not found in rtree, and cannot be deleted.", feat)
	}

	if err := g.cache.Delete(key, nil); err != nil {
		return errors.Wrap(err, "could not delete key")
	}

	return nil

}

// Get ...
func (g *Geostore) Get(key []byte) (*Feature, error) {
	response, err := g.cache.Get(key, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not get data in cache")
	}
	return NewFeatureFromJSON(response)
}

func (g *Geostore) Contains(feat *Feature) ([]*Feature, error) {
	response := g.tree.SearchIntersect(feat.Bounds())
	list := []*Feature{}
	for i := range response {
		// if f, ok = response[i].(Feature); !ok {
		// 	log.Warning("Response from rtree could not be typed as a feature.")
		// 	continue
		// }
		f := response[i].(*Feature)
		ok, err := f.Contains(feat)
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}
		list = append(list, f)
	}
	return list, nil
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
	iter := g.cache.NewIterator(nil, nil)
	for iter.Next() {
		batch.Delete(iter.Key())
	}
	iter.Release()
	if iter.Error() != nil {
		er = iter.Error()
		return
	}
	er = g.cache.Write(batch, nil)
	g.tree = rtreego.NewTree(2, 25, 50)
	return
}

func (g *Geostore) Length() (int, error) {
	var count int
	iter := g.cache.NewIterator(nil, nil)
	for iter.Next() {
		count = count + 1
	}
	iter.Release()
	if iter.Error() != nil {
		return 0, iter.Error()
	}
	treeSize := g.tree.Size()
	if treeSize != count {
		return 0, fmt.Errorf("Expected both the store and the rtree to have length %d.", count)
	}
	return count, nil

}
