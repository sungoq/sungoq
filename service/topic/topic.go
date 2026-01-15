package topic

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/dgraph-io/badger/v3"
	"github.com/sungoq/sungoq/model"
)

type TopicService struct {
	storageLocationPrefix string
}

func New() (*TopicService, error) {
	service := &TopicService{
		storageLocationPrefix: "/tmp/sungoq",
	}

	return service, nil
}

func (service *TopicService) storage() (*badger.DB, error) {
	storageOpt := badger.DefaultOptions(
		fmt.Sprintf("%s/%s", service.storageLocationPrefix, "sungoq"),
	)
	store, err := badger.Open(
		storageOpt,
	)
	if err != nil {
		return nil, err
	}

	return store, nil
}

func (service *TopicService) GetAll() ([]string, error) {
	topics := make([]string, 0)
	serviceStorage, err := service.storage()
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = serviceStorage.Close()
	}()

	err = serviceStorage.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			topics = append(topics, string(k))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return topics, nil
}

func (service *TopicService) Create(name string) error {
	serviceStorage, err := service.storage()
	if err != nil {
		return err
	}
	defer func() {
		_ = serviceStorage.Close()
	}()

	err = serviceStorage.Update(func(txn *badger.Txn) error {

		_, err := txn.Get([]byte(name))

		if err != nil {

			if errors.Is(err, badger.ErrKeyNotFound) {
				err = txn.Set([]byte(name), []byte(name))
				if err != nil {
					return err
				}
			}

			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	storage, err := badger.Open(
		badger.DefaultOptions(
			fmt.Sprintf("%s/%s", service.storageLocationPrefix, name),
		),
	)

	if err != nil {
		return err
	}

	defer func() {
		_ = storage.Close()
	}()

	return nil
}

func (service *TopicService) Delete(name string) error {
	serviceStorage, err := service.storage()
	if err != nil {
		return err
	}
	defer func() {
		_ = serviceStorage.Close()
	}()

	err = os.RemoveAll(fmt.Sprintf("%s/%s", service.storageLocationPrefix, name))
	if err != nil {
		return err
	}

	err = serviceStorage.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(name))
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (service TopicService) Publish(topic string, message interface{}) (model.Message, error) {
	storage, err := badger.Open(
		badger.DefaultOptions(
			fmt.Sprintf("%s/%s", service.storageLocationPrefix, topic),
		),
	)
	if err != nil {
		return model.Message{}, err
	}

	defer func() {
		_ = storage.Close()
	}()

	newMessage := model.NewMessage(message)

	err = storage.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(newMessage.ID), newMessage.ToJSON())
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return model.Message{}, err
	}

	return newMessage, nil
}

func (service TopicService) GetAllMessages(topic string) ([]model.Message, error) {
	storage, err := badger.Open(
		badger.DefaultOptions(
			fmt.Sprintf("%s/%s", service.storageLocationPrefix, topic),
		),
	)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = storage.Close()
	}()

	messagesRaw := make([][]byte, 0)

	err = storage.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			err := item.Value(func(v []byte) error {
				messagesRaw = append(messagesRaw, v)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	messages := make(model.Messages, 0)

	for _, mraw := range messagesRaw {
		m := model.Message{}
		if err := json.Unmarshal(mraw, &m); err != nil {
			continue
		}

		messages = append(messages, m)
	}

	sort.Sort(messages)

	return messages, nil
}

func (service TopicService) DeleteMessage(topic string, id string) error {
	storage, err := badger.Open(
		badger.DefaultOptions(
			fmt.Sprintf("%s/%s", service.storageLocationPrefix, topic),
		),
	)
	if err != nil {
		return err
	}

	defer func() {
		_ = storage.Close()
	}()

	err = storage.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(id))
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
