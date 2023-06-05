package storage

import "github.com/boltdb/bolt"

type BoltStorage struct {
	dbName     string
	bucketName []byte
}

func NewBoltStorage(dbName string) (*BoltStorage, error) {
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	bucketName := []byte("urlBucket")
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &BoltStorage{
		dbName:     dbName,
		bucketName: bucketName,
	}, nil
}

func (s *BoltStorage) Save(shortPath, srcURL string) error {
	db, err := bolt.Open(s.dbName, 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(s.bucketName))

		err := b.Put([]byte(shortPath), []byte(srcURL))
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (s *BoltStorage) GetSourceURL(shortPath string) (string, error) {
	db, err := bolt.Open(s.dbName, 0600, nil)
	if err != nil {
		return "", err
	}
	defer db.Close()

	var srcURL string

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(s.bucketName))

		srcURL = string(b.Get([]byte(shortPath)))

		return nil
	})

	return srcURL, nil
}

func (s *BoltStorage) GetShortPath(srcURL string) (string, error) {
	db, err := bolt.Open(s.dbName, 0600, nil)
	if err != nil {
		return "", err
	}
	defer db.Close()

	var shortPath string
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(s.bucketName))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			str := string(v)
			if str == srcURL {
				shortPath = string(k)
			}
		}

		return nil
	})

	return shortPath, nil
}

func (s *BoltStorage) DeleteSourceURL(srcURL string) error {

	return nil
}
