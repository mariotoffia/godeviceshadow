package mempersistence

import (
	"fmt"

	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
)

func (s *Store) List(pk string) persistencemodel.ListResults {
	s.mu.RLock()
	defer s.mu.RUnlock()

	toListResult := func(id, name string, mt persistencemodel.ModelType, entry *modelEntry) persistencemodel.ListResult {
		return persistencemodel.ListResult{
			ID:          persistencemodel.PersistenceID{ID: id, Name: name, ModelType: mt},
			Version:     entry.version,
			TimeStamp:   entry.timestamp,
			ClientToken: entry.clientToken,
		}
	}

	if pk == "" {
		res := make([]persistencemodel.ListResult, 0, len(s.partitions))

		for id, partition := range s.partitions {
			for name, entry := range partition {
				if entry.modelType == 0 /*combined*/ {
					if entry.reported != nil {
						res = append(res, toListResult(id, name, persistencemodel.ModelTypeReported, entry))
					}

					if entry.desired != nil {
						res = append(res, toListResult(id, name, persistencemodel.ModelTypeDesired, entry))
					}
				} else {
					res = append(res, toListResult(id, name, entry.modelType, entry))
				}
			}
		}

		return persistencemodel.ListResults{Items: res}
	}

	if partition, ok := s.partitions[pk]; ok {
		res := make([]persistencemodel.ListResult, 0, len(partition))

		for name, entry := range partition {
			if entry.modelType == 0 /*combined*/ {
				if entry.reported != nil {
					res = append(res, toListResult(pk, name, persistencemodel.ModelTypeReported, entry))
				}

				if entry.desired != nil {
					res = append(res, toListResult(pk, name, persistencemodel.ModelTypeDesired, entry))
				}
			} else {
				res = append(res, toListResult(pk, name, entry.modelType, entry))
			}
		}

		return persistencemodel.ListResults{Items: res}
	}

	return persistencemodel.ListResults{}
}

func (s *Store) DeleteEntry(mt persistencemodel.ModelType, pk, sk string, version int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if partition, ok := s.partitions[pk]; ok {
		if entry, ok := partition[renderSortKey(mt, sk)]; ok {
			if version > 0 && entry.version != version {
				return persistencemodel.Error409("Conflict, version mismatch")
			}

			delete(partition, sk)
			return nil
		}

		return persistencemodel.Error404(fmt.Sprintf("entry: %s - Not found", sk))
	}

	return persistencemodel.Error404(fmt.Sprintf("partition: %s - Not found", pk))
}

func (s *Store) GetEntry(mt persistencemodel.ModelType, pk string, sk string, version int64) (*modelEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if partition, ok := s.partitions[pk]; ok {
		if entry, ok := partition[renderSortKey(mt, sk)]; ok {
			if version > 0 && entry.version != version {
				return nil, persistencemodel.Error409("Conflict, version mismatch")
			}

			return entry, nil
		}
	}

	return nil, persistencemodel.Error404("Not found")
}

func (s *Store) StoreEntry(mt persistencemodel.ModelType, pk, sk string, entry *modelEntry) (*modelEntry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if p, ok := s.partitions[pk]; ok {
		if s, ok := p[renderSortKey(mt, sk)]; ok {
			if entry.version > 0 && s.version != entry.version {
				return nil, persistencemodel.Error409("Conflict, version mismatch")
			}

			entry.version++
			entry.modelType = mt

			p[sk] = entry
		} else {
			entry.version = 1
			p[sk] = entry
		}
	} else {
		entry.version = 1
		entry.modelType = mt
		s.partitions[pk] = Partition{renderSortKey(mt, sk): entry}
	}

	return entry, nil
}
