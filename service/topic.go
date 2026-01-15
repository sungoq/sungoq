package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/dgraph-io/badger/v4"
	"github.com/sungoq/sungoq/model"
)

const storageLocationPrefix = "/tmp/sungoq"

var store *badger.DB

func topicConfigure() {
	storageOpt := badger.DefaultOptions(
		fmt.Sprintf("%s/%s", storageLocationPrefix, "sungoq"),
	)
	_store, err := badger.Open(
		storageOpt,
	)
	if err != nil {
		panic(err)
	}

	store = _store
}

func TopicGetAll() ([]string, error) {
	topics := make([]string, 0)

	err := store.View(func(txn *badger.Txn) error {
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

func TopicCreate(name string) error {
	err := store.Update(func(txn *badger.Txn) error {
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
			fmt.Sprintf("%s/%s", storageLocationPrefix, name),
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

func TopicDelete(name string) error {
	err := os.RemoveAll(fmt.Sprintf("%s/%s", storageLocationPrefix, name))
	if err != nil {
		return err
	}

	err = store.Update(func(txn *badger.Txn) error {
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

func TopicPublish(topic string, message any) (model.Message, error) {
	storage, err := badger.Open(
		badger.DefaultOptions(
			fmt.Sprintf("%s/%s", storageLocationPrefix, topic),
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

func TopicGetAllMessages(topic string) ([]model.Message, error) {
	storage, err := badger.Open(
		badger.DefaultOptions(
			fmt.Sprintf("%s/%s", storageLocationPrefix, topic),
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

func TopicDeleteMessage(topic string, id string) error {
	storage, err := badger.Open(
		badger.DefaultOptions(
			fmt.Sprintf("%s/%s", storageLocationPrefix, topic),
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
