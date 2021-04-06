package kvstore

import (
	"github.com/etcd-io/bbolt"
	"github.com/pkg/errors"
	"github.com/zmon-deploy/zmon-common-go/misc"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type KvStore interface {
	Init(buckets []string) error
	Close() error
	Write(bucket string, batches map[string][]byte) error
	Read(bucket, key string) ([]byte, error)
}

type kvStore struct {
	sync.RWMutex
	db *bbolt.DB
}

func NewKvStore(dbFile string) (KvStore, error) {
	db, err := newBoltDB(dbFile)
	if err != nil {
		return nil, err
	}

	return &kvStore{db: db}, nil
}

func (s *kvStore) Init(buckets []string) error {
	s.Lock()
	defer s.Unlock()

	return s.db.Update(func(tx *bbolt.Tx) error {
		for _, bucket := range buckets {
			if _, err := tx.CreateBucketIfNotExists([]byte(bucket)); err != nil {
				return errors.Wrap(err, "failed to create bucket")
			}
		}
		return nil
	})
}

func (s *kvStore) Close() error {
	s.Lock()
	defer s.Unlock()

	if err := s.db.Close(); err != nil {
		return errors.Wrap(err, "failed to close db")
	}
	if err := os.Remove(s.db.Path()); err != nil {
		return errors.Wrap(err, "failed to remove db file")
	}
	return nil
}

func (s *kvStore) Write(bucket string, batches map[string][]byte) error {
	s.Lock()
	defer s.Unlock()

	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		for key, value := range batches {
			if err := b.Put([]byte(key), value); err != nil {
				return errors.Wrap(err, "failed to write")
			}
		}
		return nil
	})
}

func (s *kvStore) Read(bucket, key string) ([]byte, error) {
	s.RLock()
	defer s.RUnlock()

	var value []byte

	if err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		value = b.Get([]byte(key))
		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "failed to view db")
	}

	return value, nil
}

func newBoltDB(dbFile string) (*bbolt.DB, error) {
	misc.EnsureDirExist(filepath.Dir(dbFile))

	db, err := bbolt.Open(dbFile, 0600, &bbolt.Options{
		Timeout: 10 * time.Second,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to open bolt db")
	}

	return db, nil
}
