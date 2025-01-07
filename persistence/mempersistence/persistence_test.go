package mempersistence_test

import (
	"context"
	"testing"

	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
	"github.com/mariotoffia/godeviceshadow/persistence/mempersistence"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListEmptyStore(t *testing.T) {
	persistence := mempersistence.New()

	results, err := persistence.List(context.TODO(), persistencemodel.ListOptions{})

	assert.NoError(t, err)
	assert.Empty(t, results, "Results should be empty when the store is empty")
}

func TestWriteAndListSingleModel(t *testing.T) {
	persistence := mempersistence.New()

	writeResults := persistence.Write(context.TODO(), persistencemodel.WriteOptions{}, persistencemodel.WriteOperation{
		ID: persistencemodel.PersistenceID{
			ID:        "device123",
			Name:      "HomeHub",
			ModelType: persistencemodel.ModelTypeReported,
		},
		Model: map[string]any{
			"temperature": 22.5,
		},
	})

	assert.Len(t, writeResults, 1)
	assert.NoError(t, writeResults[0].Error)

	listResults, err := persistence.List(context.TODO(), persistencemodel.ListOptions{
		ID: "device123",
	})

	assert.NoError(t, err)

	assert.Len(t, listResults, 1, "There should be one model listed")
	assert.Equal(t, "device123", listResults[0].ID.ID, "ID should match")
	assert.Equal(t, "HomeHub", listResults[0].ID.Name, "Name should match")
	assert.Equal(t, persistencemodel.ModelTypeReported, listResults[0].ID.ModelType, "ModelType should match")
	assert.Greater(t, listResults[0].Version, int64(0), "Version should be greater than 0")
}

func TestWriteAndReadSingleModel(t *testing.T) {
	persistence := mempersistence.New()

	writeResults := persistence.Write(context.TODO(), persistencemodel.WriteOptions{}, persistencemodel.WriteOperation{
		ID: persistencemodel.PersistenceID{
			ID:        "device123",
			Name:      "HomeHub",
			ModelType: persistencemodel.ModelTypeReported,
		},
		Model: map[string]any{
			"temperature": 22.5,
		},
	})

	assert.Len(t, writeResults, 1)
	assert.NoError(t, writeResults[0].Error)

	readResults := persistence.Read(context.TODO(), persistencemodel.ReadOptions{}, persistencemodel.ReadOperation{
		ID: persistencemodel.PersistenceID{
			ID:        "device123",
			Name:      "HomeHub",
			ModelType: persistencemodel.ModelTypeReported,
		},
	})

	assert.Len(t, readResults, 1)
	assert.NoError(t, readResults[0].Error, "Read operation should not return an error")
	assert.Equal(t, "device123", readResults[0].ID.ID, "ID should match")
	assert.Equal(t, "HomeHub", readResults[0].ID.Name, "Name should match")
	assert.Equal(t, persistencemodel.ModelTypeReported, readResults[0].ID.ModelType, "ModelType should match")
	assert.NotNil(t, readResults[0].Model, "Model should not be nil")
	assert.Equal(t, map[string]any{"temperature": 22.5}, readResults[0].Model, "Model data should match")
}

func TestWriteAndDeleteSingleModel(t *testing.T) {
	persistence := mempersistence.New()

	writeResults := persistence.Write(context.TODO(), persistencemodel.WriteOptions{}, persistencemodel.WriteOperation{
		ID: persistencemodel.PersistenceID{
			ID:        "device123",
			Name:      "HomeHub",
			ModelType: persistencemodel.ModelTypeReported,
		},
		Model: map[string]any{
			"temperature": 22.5,
		},
	})

	assert.Len(t, writeResults, 1)
	assert.NoError(t, writeResults[0].Error)

	deleteResults := persistence.Delete(context.TODO(), persistencemodel.WriteOptions{}, persistencemodel.WriteOperation{
		ID: persistencemodel.PersistenceID{
			ID:        "device123",
			Name:      "HomeHub",
			ModelType: persistencemodel.ModelTypeReported,
		},
	})

	assert.Len(t, deleteResults, 1)
	assert.NoError(t, deleteResults[0].Error)

	readResults := persistence.Read(context.TODO(), persistencemodel.ReadOptions{}, persistencemodel.ReadOperation{
		ID: persistencemodel.PersistenceID{
			ID:        "device123",
			Name:      "HomeHub",
			ModelType: persistencemodel.ModelTypeReported,
		},
	})

	assert.Len(t, readResults, 1)
	assert.Error(t, readResults[0].Error, "Read operation should return an error for a deleted model")
	assert.Equal(t, 404, readResults[0].Error.(persistencemodel.PersistenceError).Code)
}

func TestWriteVersionConflict(t *testing.T) {
	persistence := mempersistence.New()

	writeResults := persistence.Write(context.TODO(), persistencemodel.WriteOptions{}, persistencemodel.WriteOperation{
		ID: persistencemodel.PersistenceID{
			ID:        "device123",
			Name:      "HomeHub",
			ModelType: persistencemodel.ModelTypeReported,
		},
		Model: map[string]any{
			"temperature": 22.5,
		},
	})

	assert.Len(t, writeResults, 1)
	assert.NoError(t, writeResults[0].Error)

	conflictResults := persistence.Write(context.TODO(), persistencemodel.WriteOptions{}, persistencemodel.WriteOperation{
		ID: persistencemodel.PersistenceID{
			ID:        "device123",
			Name:      "HomeHub",
			ModelType: persistencemodel.ModelTypeReported,
		},
		Model: map[string]any{
			"temperature": 23.0, // Updated value
		},
		Version: 99, // Incorrect version
	})

	assert.Len(t, conflictResults, 1, "There should be one result for the conflicting write operation")
	assert.Error(t, conflictResults[0].Error, "Write operation should return an error for version conflict")
	assert.Equal(t, 409, conflictResults[0].Error.(persistencemodel.PersistenceError).Code, "Error code should be 409 (Conflict)")
}

func TestDeleteVersionConflict(t *testing.T) {
	persistence := mempersistence.New()

	writeResults := persistence.Write(context.TODO(), persistencemodel.WriteOptions{}, persistencemodel.WriteOperation{
		ID: persistencemodel.PersistenceID{
			ID:        "device123",
			Name:      "HomeHub",
			ModelType: persistencemodel.ModelTypeReported,
		},
		Model: map[string]any{
			"temperature": 22.5,
		},
	})

	assert.Len(t, writeResults, 1)
	assert.NoError(t, writeResults[0].Error)

	deleteResults := persistence.Delete(context.TODO(), persistencemodel.WriteOptions{}, persistencemodel.WriteOperation{
		ID: persistencemodel.PersistenceID{
			ID:        "device123",
			Name:      "HomeHub",
			ModelType: persistencemodel.ModelTypeReported,
		},
		Version: 99, // Incorrect version
	})

	assert.Len(t, deleteResults, 1)
	assert.Error(t, deleteResults[0].Error, "Delete operation should return an error for version conflict")
	assert.Equal(t, 409, deleteResults[0].Error.(persistencemodel.PersistenceError).Code, "Error code should be 409 (Conflict)")
}

func TestDeleteWithoutVersionConstraint(t *testing.T) {
	persistence := mempersistence.New()

	writeResults := persistence.Write(context.TODO(), persistencemodel.WriteOptions{}, persistencemodel.WriteOperation{
		ID: persistencemodel.PersistenceID{
			ID:        "device123",
			Name:      "HomeHub",
			ModelType: persistencemodel.ModelTypeReported,
		},
		Model: map[string]any{
			"temperature": 22.5,
		},
	})

	assert.Len(t, writeResults, 1)
	assert.NoError(t, writeResults[0].Error)

	deleteResults := persistence.Delete(context.TODO(), persistencemodel.WriteOptions{}, persistencemodel.WriteOperation{
		ID: persistencemodel.PersistenceID{
			ID:        "device123",
			Name:      "HomeHub",
			ModelType: persistencemodel.ModelTypeReported,
		},
		Version: 0, // No version constraint
	})

	assert.Len(t, deleteResults, 1)
	assert.NoError(t, deleteResults[0].Error)

	readResults := persistence.Read(context.TODO(), persistencemodel.ReadOptions{}, persistencemodel.ReadOperation{
		ID: persistencemodel.PersistenceID{
			ID:        "device123",
			Name:      "HomeHub",
			ModelType: persistencemodel.ModelTypeReported,
		},
	})

	assert.Len(t, readResults, 1)
	assert.Error(t, readResults[0].Error, "Read operation should return an error for a deleted model")
	assert.Equal(t, 404, readResults[0].Error.(persistencemodel.PersistenceError).Code)
}

func TestListMultipleModels(t *testing.T) {
	persistence := mempersistence.New()

	writeResults := persistence.Write(context.TODO(), persistencemodel.WriteOptions{},
		persistencemodel.WriteOperation{
			ID: persistencemodel.PersistenceID{
				ID:        "device123",
				Name:      "HomeHub",
				ModelType: persistencemodel.ModelTypeReported,
			},
			Model: map[string]any{
				"temperature": 22.5,
			},
		},
		persistencemodel.WriteOperation{
			ID: persistencemodel.PersistenceID{
				ID:        "device124",
				Name:      "Car",
				ModelType: persistencemodel.ModelTypeDesired,
			},
			Model: map[string]any{
				"temperature": 18.0,
			},
		},
	)

	assert.Len(t, writeResults, 2, "There should be two results for the write operations")
	for _, result := range writeResults {
		assert.NoError(t, result.Error)
	}

	listResults, err := persistence.List(context.TODO(), persistencemodel.ListOptions{})

	assert.NoError(t, err)

	assert.Len(t, listResults, 2, "There should be two models listed")
	assert.ElementsMatch(t, []string{"device123", "device124"}, []string{
		listResults[0].ID.ID, listResults[1].ID.ID,
	}, "IDs of listed models should match")
	assert.ElementsMatch(t, []string{"HomeHub", "Car"}, []string{
		listResults[0].ID.Name, listResults[1].ID.Name,
	}, "Names of listed models should match")
}

func TestListWithIDFilter(t *testing.T) {
	persistence := mempersistence.New()

	persistence.Write(context.TODO(), persistencemodel.WriteOptions{},
		persistencemodel.WriteOperation{
			ID: persistencemodel.PersistenceID{
				ID:        "device123",
				Name:      "HomeHub",
				ModelType: persistencemodel.ModelTypeReported,
			},
			Model: map[string]any{
				"temperature": 22.5,
			},
		},
		persistencemodel.WriteOperation{
			ID: persistencemodel.PersistenceID{
				ID:        "device124",
				Name:      "Car",
				ModelType: persistencemodel.ModelTypeDesired,
			},
			Model: map[string]any{
				"temperature": 18.0,
			},
		},
	)

	listResults, err := persistence.List(context.TODO(), persistencemodel.ListOptions{
		ID: "device123",
	})

	assert.NoError(t, err)

	assert.Len(t, listResults, 1, "There should be one model listed")
	assert.Equal(t, "device123", listResults[0].ID.ID, "Listed model ID should match the filter")
	assert.Equal(t, "HomeHub", listResults[0].ID.Name, "Listed model Name should match")
	assert.Equal(t, persistencemodel.ModelTypeReported, listResults[0].ID.ModelType, "Listed model ModelType should match")
}

func TestReadNonExistentModel(t *testing.T) {
	persistence := mempersistence.New()

	readResults := persistence.Read(context.TODO(), persistencemodel.ReadOptions{}, persistencemodel.ReadOperation{
		ID: persistencemodel.PersistenceID{
			ID:        "nonexistent123",
			Name:      "HomeHub",
			ModelType: persistencemodel.ModelTypeReported,
		},
	})

	assert.Len(t, readResults, 1)
	assert.Error(t, readResults[0].Error, "Read operation should return an error for a non-existent model")
	assert.Equal(t, 404, readResults[0].Error.(persistencemodel.PersistenceError).Code)
	assert.Equal(t, "nonexistent123", readResults[0].ID.ID, "Returned ID should match the requested ID")
	assert.Equal(t, "HomeHub", readResults[0].ID.Name, "Returned Name should match the requested Name")
	assert.Equal(t, persistencemodel.ModelTypeReported, readResults[0].ID.ModelType, "Returned ModelType should match the requested ModelType")
}

func TestWriteUpdateExistingModel(t *testing.T) {
	persistence := mempersistence.New()

	initialWrite := persistence.Write(context.TODO(), persistencemodel.WriteOptions{}, persistencemodel.WriteOperation{
		ID: persistencemodel.PersistenceID{
			ID:        "device123",
			Name:      "HomeHub",
			ModelType: persistencemodel.ModelTypeReported,
		},
		Model: map[string]any{
			"temperature": 22.5,
		},
	})

	assert.Len(t, initialWrite, 1, "There should be one result for the initial write operation")
	assert.NoError(t, initialWrite[0].Error, "Initial write operation should not return an error")

	// Update the model
	updateWrite := persistence.Write(context.TODO(), persistencemodel.WriteOptions{}, persistencemodel.WriteOperation{
		ID: persistencemodel.PersistenceID{
			ID:        "device123",
			Name:      "HomeHub",
			ModelType: persistencemodel.ModelTypeReported,
		},
		Model: map[string]any{
			"temperature": 25.0,
		},
		Version: 1, // Use the correct version
	})

	assert.Len(t, updateWrite, 1, "There should be one result for the update write operation")
	assert.NoError(t, updateWrite[0].Error, "Update write operation should not return an error")

	readResults := persistence.Read(context.TODO(), persistencemodel.ReadOptions{}, persistencemodel.ReadOperation{
		ID: persistencemodel.PersistenceID{
			ID:        "device123",
			Name:      "HomeHub",
			ModelType: persistencemodel.ModelTypeReported,
		},
	})

	assert.Len(t, readResults, 1)
	assert.NoError(t, readResults[0].Error, "Read operation should not return an error")
	assert.Equal(t, map[string]any{"temperature": 25.0}, readResults[0].Model, "Updated model data should match")
	assert.Equal(t, readResults[0].Version, int64(2), "Since written twice")
}

func TestDeleteNonExistentModel(t *testing.T) {
	persistence := mempersistence.New()

	deleteResults := persistence.Delete(context.TODO(), persistencemodel.WriteOptions{}, persistencemodel.WriteOperation{
		ID: persistencemodel.PersistenceID{
			ID:        "nonexistent123",
			Name:      "HomeHub",
			ModelType: persistencemodel.ModelTypeReported,
		},
	})

	assert.Len(t, deleteResults, 1)
	assert.Error(t, deleteResults[0].Error, "Delete operation should return an error for a non-existent model")
	assert.Equal(t, 404, deleteResults[0].Error.(persistencemodel.PersistenceError).Code)
	assert.Equal(t, "nonexistent123", deleteResults[0].ID.ID, "Returned ID should match the requested ID")
	assert.Equal(t, "HomeHub", deleteResults[0].ID.Name, "Returned Name should match the requested Name")
	assert.Equal(t, persistencemodel.ModelTypeReported, deleteResults[0].ID.ModelType, "Returned ModelType should match the requested ModelType")
}

func TestListAfterDeletingAllModels(t *testing.T) {
	persistence := mempersistence.New()

	persistence.Write(context.TODO(), persistencemodel.WriteOptions{},
		persistencemodel.WriteOperation{
			ID: persistencemodel.PersistenceID{
				ID:        "device123",
				Name:      "HomeHub",
				ModelType: persistencemodel.ModelTypeReported,
			},
			Model: map[string]any{
				"temperature": 22.5,
			},
		},
		persistencemodel.WriteOperation{
			ID: persistencemodel.PersistenceID{
				ID:        "device124",
				Name:      "Car",
				ModelType: persistencemodel.ModelTypeDesired,
			},
			Model: map[string]any{
				"temperature": 18.0,
			},
		},
	)

	listResults, err := persistence.List(context.TODO(), persistencemodel.ListOptions{})
	require.NoError(t, err)
	require.Len(t, listResults, 2, "There should be two models listed before deletion")

	persistence.Delete(context.TODO(), persistencemodel.WriteOptions{},
		persistencemodel.WriteOperation{
			ID: persistencemodel.PersistenceID{
				ID:        "device123",
				Name:      "HomeHub",
				ModelType: persistencemodel.ModelTypeReported,
			},
		},
		persistencemodel.WriteOperation{
			ID: persistencemodel.PersistenceID{
				ID:        "device124",
				Name:      "Car",
				ModelType: persistencemodel.ModelTypeDesired,
			},
		},
	)

	listResults, err = persistence.List(context.TODO(), persistencemodel.ListOptions{})

	assert.NoError(t, err)
	assert.Empty(t, listResults, "Results should be empty after deleting all models")
}

func TestWriteIdenticalIDDifferentNames(t *testing.T) {
	persistence := mempersistence.New()

	writeResults := persistence.Write(context.TODO(), persistencemodel.WriteOptions{},
		persistencemodel.WriteOperation{
			ID: persistencemodel.PersistenceID{
				ID:        "device123",
				Name:      "HomeHub",
				ModelType: persistencemodel.ModelTypeReported,
			},
			Model: map[string]any{
				"temperature": 22.5,
			},
		},
		persistencemodel.WriteOperation{
			ID: persistencemodel.PersistenceID{
				ID:        "device123",
				Name:      "Car",
				ModelType: persistencemodel.ModelTypeDesired,
			},
			Model: map[string]any{
				"temperature": 20.0,
			},
		},
	)

	assert.Len(t, writeResults, 2, "There should be two results for the write operations")
	for _, result := range writeResults {
		assert.NoError(t, result.Error)
	}

	listResults, err := persistence.List(context.TODO(), persistencemodel.ListOptions{
		ID: "device123",
	})

	assert.NoError(t, err)
	assert.Len(t, listResults, 2, "There should be two models listed")
	assert.ElementsMatch(t, []string{"HomeHub", "Car"}, []string{
		listResults[0].ID.Name, listResults[1].ID.Name,
	}, "Names of listed models should match")
}

func TestDeleteByNameForIdenticalID(t *testing.T) {
	persistence := mempersistence.New()

	persistence.Write(context.TODO(), persistencemodel.WriteOptions{},
		persistencemodel.WriteOperation{
			ID: persistencemodel.PersistenceID{
				ID:        "device123",
				Name:      "HomeHub",
				ModelType: persistencemodel.ModelTypeReported,
			},
			Model: map[string]any{
				"temperature": 22.5,
			},
		},
		persistencemodel.WriteOperation{
			ID: persistencemodel.PersistenceID{
				ID:        "device123",
				Name:      "Car",
				ModelType: persistencemodel.ModelTypeDesired,
			},
			Model: map[string]any{
				"temperature": 20.0,
			},
		},
	)

	deleteResults := persistence.Delete(context.TODO(), persistencemodel.WriteOptions{}, persistencemodel.WriteOperation{
		ID: persistencemodel.PersistenceID{
			ID:        "device123",
			Name:      "HomeHub",
			ModelType: persistencemodel.ModelTypeReported,
		},
	})

	assert.Len(t, deleteResults, 1)
	assert.NoError(t, deleteResults[0].Error)

	listResults, err := persistence.List(context.TODO(), persistencemodel.ListOptions{
		ID: "device123",
	})

	assert.NoError(t, err)

	assert.Len(t, listResults, 1, "There should be one model remaining")
	assert.Equal(t, "Car", listResults[0].ID.Name, "Remaining model should have Name 'Desired'")
	assert.Equal(t, persistencemodel.ModelTypeDesired, listResults[0].ID.ModelType, "Remaining model should have ModelType 'Desired'")
}

func TestWriteUpdatesOnlySpecifiedModel(t *testing.T) {
	persistence := mempersistence.New()

	persistence.Write(context.TODO(), persistencemodel.WriteOptions{},
		persistencemodel.WriteOperation{
			ID: persistencemodel.PersistenceID{
				ID:        "device123",
				Name:      "HomeHub",
				ModelType: persistencemodel.ModelTypeReported,
			},
			Model: map[string]any{
				"temperature": 22.5,
			},
		},
		persistencemodel.WriteOperation{
			ID: persistencemodel.PersistenceID{
				ID:        "device123",
				Name:      "Car",
				ModelType: persistencemodel.ModelTypeDesired,
			},
			Model: map[string]any{
				"temperature": 20.0,
			},
		},
	)

	writeResults := persistence.Write(context.TODO(), persistencemodel.WriteOptions{}, persistencemodel.WriteOperation{
		ID: persistencemodel.PersistenceID{
			ID:        "device123",
			Name:      "HomeHub",
			ModelType: persistencemodel.ModelTypeReported,
		},
		Model: map[string]any{
			"temperature": 23.0, // New temperature
		},
		Version: 1, // Correct version for "HomeHub"
	})

	assert.Len(t, writeResults, 1)
	assert.NoError(t, writeResults[0].Error)

	readResults := persistence.Read(context.TODO(), persistencemodel.ReadOptions{},
		persistencemodel.ReadOperation{
			ID: persistencemodel.PersistenceID{
				ID:        "device123",
				Name:      "HomeHub",
				ModelType: persistencemodel.ModelTypeReported,
			},
		},
		persistencemodel.ReadOperation{
			ID: persistencemodel.PersistenceID{
				ID:        "device123",
				Name:      "Car",
				ModelType: persistencemodel.ModelTypeDesired,
			},
		},
	)

	assert.Len(t, readResults, 2, "There should be two results for the read operation")
	assert.NoError(t, readResults[0].Error, "Read operation should not return an error for 'Reported'")
	assert.Equal(t, map[string]any{"temperature": 23.0}, readResults[0].Model, "'Reported' model should be updated")
	assert.NoError(t, readResults[1].Error, "Read operation should not return an error for 'Desired'")
	assert.Equal(t, map[string]any{"temperature": 20.0}, readResults[1].Model, "'Desired' model should remain unchanged")
}

func TestListNoMatchingID(t *testing.T) {
	persistence := mempersistence.New()

	persistence.Write(context.TODO(), persistencemodel.WriteOptions{},
		persistencemodel.WriteOperation{
			ID: persistencemodel.PersistenceID{
				ID:        "device123",
				Name:      "HomeHub",
				ModelType: persistencemodel.ModelTypeReported,
			},
			Model: map[string]any{
				"temperature": 22.5,
			},
		},
		persistencemodel.WriteOperation{
			ID: persistencemodel.PersistenceID{
				ID:        "device124",
				Name:      "Car",
				ModelType: persistencemodel.ModelTypeDesired,
			},
			Model: map[string]any{
				"temperature": 18.0,
			},
		},
	)

	listResults, err := persistence.List(context.TODO(), persistencemodel.ListOptions{
		ID: "device999", // ID that does not exist
	})

	assert.NoError(t, err)
	assert.Empty(t, listResults, "Results should be empty for a non-matching ID")
}
