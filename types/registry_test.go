package types_test

import (
	"reflect"
	"strconv"
	"sync"
	"testing"

	"github.com/mariotoffia/godeviceshadow/model"
	"github.com/mariotoffia/godeviceshadow/types"

	"github.com/stretchr/testify/assert"
)

// SomeStruct is a sample struct used in our tests:
type SomeStruct struct{}

// AnotherStruct is a second sample struct used in our tests:
type AnotherStruct struct{}

// ConcurrentStruct is a sample struct used for concurrent registration tests.
type ConcurrentStruct struct{}

func TestRegisterWithExplicitName(t *testing.T) {
	// 1. Create a TypeRegistryImpl instance.
	registry := types.NewRegistry()

	// 2. Register the type with a specific (non-empty) name.
	explicitName := "MyExplicitName"
	registry.Register(explicitName, &SomeStruct{})

	// 3. Invoke Get to retrieve the entry.
	te, ok := registry.Get(explicitName)

	// 4. Assert we found the type entry.
	assert.True(t, ok, "expected to retrieve an entry for name %q", explicitName)
	assert.NotNil(t, te, "TypeEntry should not be nil")

	// 5. Assert that the retrieved name is exactly what was registered.
	assert.Equal(t, explicitName, te.Name, "expected TypeEntry name to match the explicit name")

	// 6. Check that the underlying reflect.Type is correct (i.e., SomeStruct).
	// Because we registered a pointer to SomeStruct, but the code strips pointers (Elem),
	// expect it to be SomeStruct (not *SomeStruct).
	expectedType := reflect.TypeOf(SomeStruct{})
	assert.Equal(t, expectedType, te.Model, "expected reflect.Type to match SomeStruct")
}

func TestRegisterWithEmptyNameAndRetrieve(t *testing.T) {
	registry := types.NewRegistry()

	// 1. Register with an empty name.
	registry.Register("", &SomeStruct{})

	// 2. Determine the expected name derived from package path + type name.
	//    Note: Register calls `reflect.TypeOf(t).Elem()` if it's a pointer,
	//    so we'll derive the type from SomeStruct (non-pointer).
	expectedName := reflect.TypeOf(SomeStruct{}).PkgPath() + "." + reflect.TypeOf(SomeStruct{}).Name()

	// 3. Retrieve by the derived name.
	te, ok := registry.Get(expectedName)

	// 4. Assert we found an entry.
	assert.True(t, ok, "expected to find entry for derived name %q", expectedName)
	assert.NotNil(t, te, "TypeEntry should not be nil")

	// 5. Confirm the registered name is the derived name.
	assert.Equal(t, expectedName, te.Name, "expected TypeEntry name to be packagePath.typeName")

	// 6. Confirm the underlying type is SomeStruct.
	expectedType := reflect.TypeOf(SomeStruct{})
	assert.Equal(t, expectedType, te.Model, "expected reflect.Type to match SomeStruct")
}

func TestRegisterMultipleTypesAndRetrieveEach(t *testing.T) {
	// 1. Create a TypeRegistryImpl instance.
	registry := types.NewRegistry()

	// 2. Register two different types with unique names.
	firstName := "FirstName"
	secondName := "SecondName"

	registry.Register(firstName, &SomeStruct{})
	registry.Register(secondName, &AnotherStruct{})

	// 3. Retrieve each type by its respective name.
	teFirst, okFirst := registry.Get(firstName)
	teSecond, okSecond := registry.Get(secondName)

	// 4. Assert that both types were successfully retrieved.
	assert.True(t, okFirst, "expected to retrieve an entry for name %q", firstName)
	assert.True(t, okSecond, "expected to retrieve an entry for name %q", secondName)

	// 5. Assert that the first TypeEntry has the correct name and type.
	assert.Equal(t, firstName, teFirst.Name, "expected TypeEntry name to match %q", firstName)
	expectedTypeFirst := reflect.TypeOf(SomeStruct{})
	assert.Equal(t, expectedTypeFirst, teFirst.Model, "expected reflect.Type to match SomeStruct for %q", firstName)

	// 6. Assert that the second TypeEntry has the correct name and type.
	assert.Equal(t, secondName, teSecond.Name, "expected TypeEntry name to match %q", secondName)
	expectedTypeSecond := reflect.TypeOf(AnotherStruct{})
	assert.Equal(t, expectedTypeSecond, teSecond.Model, "expected reflect.Type to match AnotherStruct for %q", secondName)

	// 7. Ensure that retrieving one name does not return the other type.
	teNonExistent, okNonExistent := registry.Get("NonExistentName")
	assert.False(t, okNonExistent, "did not expect to retrieve an entry for a non-registered name")
	assert.Equal(t, model.TypeEntry{}, teNonExistent, "expected zero-value TypeEntry for non-registered name")
}

// TestConcurrentRegistrations verifies that multiple concurrent registrations
// do not cause data races and that all types are registered correctly.
func TestConcurrentRegistrations(t *testing.T) {
	registry := types.NewRegistry()

	var wg sync.WaitGroup
	numRegistrations := 1000000
	wg.Add(numRegistrations)

	// Use a mutex to protect access to the names slice.
	var namesMutex sync.Mutex
	names := make([]string, 0, numRegistrations)

	for i := 0; i < numRegistrations; i++ {
		go func(i int) {
			defer wg.Done()
			name := "ConcurrentName_" + strconv.Itoa(i)
			registry.Register(name, &ConcurrentStruct{})

			// Safely append the name to the names slice for later verification.
			namesMutex.Lock()
			names = append(names, name)
			namesMutex.Unlock()
		}(i)
	}

	// Wait for all goroutines to finish.
	wg.Wait()

	// Verify that all registered names are present in the registry.
	for _, name := range names {
		te, ok := registry.Get(name)
		assert.True(t, ok, "expected to retrieve an entry for name %q", name)
		assert.NotNil(t, te, "TypeEntry should not be nil for name %q", name)
		assert.Equal(t, name, te.Name, "expected TypeEntry name to match %q", name)

		// Confirm the underlying type is ConcurrentStruct.
		expectedType := reflect.TypeOf(ConcurrentStruct{})
		assert.Equal(t, expectedType, te.Model, "expected reflect.Type to match ConcurrentStruct for name %q", name)
	}

	// Optionally, check the total number of registrations matches expected.
	assert.Equal(t, numRegistrations, len(names), "expected number of registered names to match")
}

// TestGetNonExistentName verifies that retrieving a non-registered name
// returns false and an empty TypeEntry.
func TestGetNonExistentName(t *testing.T) {
	// 1. Create a TypeRegistryImpl instance.
	registry := types.NewRegistry()

	// 2. Attempt to retrieve a TypeEntry using a name that hasn't been registered.
	nonExistentName := "IAmNotRegistered"
	te, ok := registry.Get(nonExistentName)

	// 3. Assert that the retrieval was unsuccessful.
	assert.False(t, ok, "expected retrieval of non-registered name %q to return ok=false", nonExistentName)

	// 4. Confirm that the returned TypeEntry is the zero value.
	assert.Equal(t, model.TypeEntry{}, te, "expected zero-value TypeEntry when retrieving non-registered name %q", nonExistentName)

	// Additional Checks (Optional):
	// Although not strictly necessary, you can verify individual fields.
	assert.Equal(t, reflect.Type(nil), te.Model, "expected TypeEntry.Model to be nil for non-registered name")
	assert.Equal(t, "", te.Name, "expected TypeEntry.Name to be empty for non-registered name")
	assert.Nil(t, te.Meta, "expected TypeEntry.Meta to be nil for non-registered name")
}

// TestRegisterDuplicateName verifies that registering the same name multiple times
// using the Register method results in the latest registration overwriting the previous one.
func TestRegisterDuplicateName(t *testing.T) {
	// Define two different structs for testing.
	type AnotherStruct struct{}

	// 1. Create a TypeRegistryImpl instance.
	registry := types.NewRegistry()

	// 2. Register the first type with the name "DupName".
	firstName := "DupName"
	registry.Register(firstName, &SomeStruct{})

	// 3. Register the second type with the same name "DupName".
	registry.Register(firstName, &AnotherStruct{})

	// 4. Retrieve the TypeEntry by the name "DupName".
	te, ok := registry.Get(firstName)

	// 5. Assert that the retrieval was successful.
	assert.True(t, ok, "expected to retrieve an entry for name %q", firstName)

	// 6. Assert that the retrieved TypeEntry has the name "DupName".
	assert.Equal(t, firstName, te.Name, "expected TypeEntry name to match %q", firstName)

	// 7. Assert that the Model is of type AnotherStruct, indicating that the second registration overwrote the first.
	expectedType := reflect.TypeOf(AnotherStruct{})
	assert.Equal(t, expectedType, te.Model, "expected reflect.Type to match AnotherStruct after duplicate registration")
}

// TestRegisterIDLookupWithValidID verifies that registering a type with a valid ID
// using RegisterIDLookup stores the entry correctly and that ResolveByID can retrieve it using the ID alone.
func TestRegisterIDLookupWithValidID(t *testing.T) {
	// Define the ID and the type to be registered.
	id := "SomeID"

	// 1. Create a TypeRegistryImpl instance.
	registry := types.NewRegistry()

	// 2. Register the type with the specified ID using RegisterIDLookup.
	err := registry.RegisterIDLookup(id, &SomeStruct{})
	assert.NoError(t, err, "RegisterIDLookup should not return an error for a valid ID")

	// 3. Use ResolveByID with the registered ID and an empty name to retrieve the TypeEntry.
	te, ok := registry.ResolveByID(id, "")

	// 4. Assert that the retrieval was successful.
	assert.True(t, ok, "ResolveByID should return ok=true for a registered ID")

	// 5. Assert that the retrieved TypeEntry has the default generated name.
	// Since RegisterIDLookup does not provide a name, the Name should be "{pkgPath}.SomeStruct".
	expectedName := reflect.TypeOf(SomeStruct{}).PkgPath() + "." + reflect.TypeOf(SomeStruct{}).Name()
	assert.Equal(t, expectedName, te.Name, "TypeEntry.Name should be the default generated name based on package path and type name")

	// 6. Confirm that the underlying type is SomeStruct.
	expectedType := reflect.TypeOf(SomeStruct{})
	assert.Equal(t, expectedType, te.Model, "TypeEntry.Model should match SomeStruct for the registered ID")
}

// TestRegisterIDLookupWithEmptyID verifies that RegisterIDLookup returns an error
// when an empty ID is provided.
func TestRegisterIDLookupWithEmptyID(t *testing.T) {
	// 1. Create a TypeRegistryImpl instance.
	registry := types.NewRegistry()

	// 2. Attempt to register a type with an empty ID.
	emptyID := ""
	err := registry.RegisterIDLookup(emptyID, &SomeStruct{})

	// 3. Assert that an error is returned.
	assert.Error(t, err, "RegisterIDLookup should return an error when ID is empty")

	// 4. Assert that the error message matches the expected message.
	expectedErrorMessage := "id must be a non empty string"
	assert.EqualError(t, err, expectedErrorMessage, "expected error message to match")
}

// TestRegisterIDLookupDuplicateID verifies that attempting to register a duplicate ID using
// RegisterIDLookup returns an error, and that the original registration remains intact.
func TestRegisterIDLookupDuplicateID(t *testing.T) {
	// 1. Create a TypeRegistryImpl instance.
	registry := types.NewRegistry()

	// 2. Register the first type (SomeStruct) with ID "SomeID".
	firstID := "SomeID"
	err := registry.RegisterIDLookup(firstID, &SomeStruct{})
	assert.NoError(t, err, "First RegisterIDLookup should succeed without error")

	// 3. Attempt to register the second type (AnotherStruct) with the same ID "SomeID".
	err = registry.RegisterIDLookup(firstID, &AnotherStruct{})
	assert.Error(t, err, "Second RegisterIDLookup with duplicate ID should return an error")
	assert.EqualError(t, err, "id already registered: SomeID", "Error message should indicate duplicate ID registration")

	// 4. Resolve by ID "SomeID" to ensure the original type (SomeStruct) is still registered.
	te, ok := registry.ResolveByID(firstID, "")
	assert.True(t, ok, "ResolveByID should find the original TypeEntry for ID 'SomeID'")
	assert.NotNil(t, te, "TypeEntry should not be nil for registered ID 'SomeID'")

	// 5. Assert that the retrieved TypeEntry corresponds to SomeStruct, not AnotherStruct.
	expectedType := reflect.TypeOf(SomeStruct{})
	assert.Equal(t, expectedType, te.Model, "TypeEntry.Model should match SomeStruct after duplicate registration attempt")
	assert.Equal(t, reflect.TypeOf(SomeStruct{}).PkgPath()+"."+reflect.TypeOf(SomeStruct{}).Name(), te.Name, "TypeEntry.Name should match the default generated name for SomeStruct")

	// 6. Optionally, verify that the 'ids' map contains only the first registration.
	//    Since 'ids' is unexported, direct access is not possible. Instead, rely on ResolveByID.
	//    Attempting to resolve with AnotherStruct's type should fail.
	anotherType := reflect.TypeOf(AnotherStruct{})
	assert.NotEqual(t, anotherType, te.Model, "TypeEntry.Model should not match AnotherStruct")
}

// TestRegisterNameLookupWithValidName verifies that RegisterNameLookup correctly registers a type with a specific name,
// and that ResolveByID can retrieve it using the name with an empty ID.
func TestRegisterNameLookupWithValidName(t *testing.T) {
	// Define the name and the type to be registered.
	name := "SomeName"

	// 1. Create a TypeRegistryImpl instance.
	registry := types.NewRegistry()

	// 2. Register the type with the specified name using RegisterNameLookup.
	err := registry.RegisterNameLookup(name, &SomeStruct{})
	assert.NoError(t, err, "RegisterNameLookup should not return an error for a valid name")

	// 3. Use ResolveByID with an empty ID and the registered name to retrieve the TypeEntry.
	te, ok := registry.ResolveByID("", name)

	// 4. Assert that the retrieval was successful.
	assert.True(t, ok, "ResolveByID should return ok=true for a registered name")
	assert.NotNil(t, te, "TypeEntry should not be nil for a registered name")

	// 5. Assert that the retrieved TypeEntry has the correct name.
	assert.Equal(t, name, te.Name, "TypeEntry.Name should match the registered name")

	// 6. Confirm that the underlying type is SomeStruct.
	expectedType := reflect.TypeOf(SomeStruct{})
	assert.Equal(t, expectedType, te.Model, "TypeEntry.Model should match SomeStruct for the registered name")

	// 7. Additionally, verify that the type can be retrieved using Get.
	teGet, okGet := registry.Get(name)
	assert.True(t, okGet, "Get should return ok=true for a registered name")
	assert.Equal(t, te, teGet, "Get should return the same TypeEntry as ResolveByID for the registered name")

	// 8. Optionally, ensure that resolving with both ID and name still works correctly.
	// Since we didn't register with ID+name, resolving with a name should still retrieve the entry.
	teResolve, okResolve := registry.ResolveByID("", name)
	assert.True(t, okResolve, "ResolveByID should return ok=true when resolving with empty ID and valid name")
	assert.Equal(t, te, teResolve, "ResolveByID should return the correct TypeEntry when resolving with empty ID and valid name")

	// 9. Confirm that resolving with a non-empty ID but valid name still retrieves the entry.
	// This tests that the name lookup works regardless of the ID parameter when the name is provided.
	someID := "AnyID"
	teResolveWithID, okResolveWithID := registry.ResolveByID(someID, name)
	assert.True(t, okResolveWithID, "ResolveByID should return ok=true when resolving with any ID and valid name")
	assert.Equal(t, te, teResolveWithID, "ResolveByID should return the correct TypeEntry when resolving with any ID and valid name")

	// 10. Finally, verify that attempting to resolve with an incorrect name does not retrieve the entry.
	incorrectName := "IncorrectName"
	teIncorrect, okIncorrect := registry.ResolveByID("", incorrectName)
	assert.False(t, okIncorrect, "ResolveByID should return ok=false for an unregistered name")
	assert.Equal(t, model.TypeEntry{}, teIncorrect, "ResolveByID should return zero-value TypeEntry for an unregistered name")
}

// TestRegisterNameLookupWithEmptyName verifies that RegisterNameLookup returns an error
// when an empty name is provided.
func TestRegisterNameLookupWithEmptyName(t *testing.T) {
	// 1. Create a TypeRegistryImpl instance.
	registry := types.NewRegistry()

	// 2. Attempt to register a type with an empty name.
	emptyName := ""
	err := registry.RegisterNameLookup(emptyName, &SomeStruct{})

	// 3. Assert that an error is returned.
	assert.Error(t, err, "RegisterNameLookup should return an error when name is empty")

	// 4. Assert that the error message matches the expected message.
	expectedErrorMessage := "name must be a non empty string"
	assert.EqualError(t, err, expectedErrorMessage, "expected error message to match")

	// 5. Optionally, verify that the type was not registered by attempting to retrieve it.
	// Since the registration failed, resolving by name should not find the type.
	te, ok := registry.ResolveByID("", emptyName)
	assert.False(t, ok, "ResolveByID should return ok=false for an empty name registration attempt")
	assert.Equal(t, model.TypeEntry{}, te, "ResolveByID should return zero-value TypeEntry for an empty name")
}

// TestRegisterIDandNameLookupWithValidIDAndName verifies that registering a type with both a valid ID and name
// using RegisterIDandNameLookup stores the entry correctly and that ResolveByID can retrieve it only when both
// ID and Name are provided together.
func TestRegisterIDandNameLookupWithValidIDAndName(t *testing.T) {
	// Define the ID, Name, and the type to be registered.
	id := "SomeID"
	name := "SomeName"

	// 1. Create a TypeRegistryImpl instance.
	registry := types.NewRegistry()

	// 2. Register the type with both ID and Name using RegisterIDandNameLookup.
	err := registry.RegisterIDandNameLookup(id, name, &SomeStruct{})
	assert.NoError(t, err, "RegisterIDandNameLookup should not return an error for valid ID and Name")

	// 3. Use ResolveByID with the registered ID and Name to retrieve the TypeEntry.
	te, ok := registry.ResolveByID(id, name)

	// 4. Assert that the retrieval was successful.
	assert.True(t, ok, "ResolveByID should return ok=true for a registered ID and Name combination")
	assert.NotNil(t, te, "TypeEntry should not be nil for a registered ID and Name combination")

	// 5. Assert that the retrieved TypeEntry has the correct Name.
	assert.Equal(t, name, te.Name, "TypeEntry.Name should match the registered name")

	// 6. Confirm that the underlying type is SomeStruct.
	expectedType := reflect.TypeOf(SomeStruct{})
	assert.Equal(t, expectedType, te.Model, "TypeEntry.Model should match SomeStruct for the registered ID and Name")

	// 7. Attempt to resolve with only ID.
	teOnlyID, okOnlyID := registry.ResolveByID(id, "")
	assert.False(t, okOnlyID, "ResolveByID should return ok=false when resolving with only ID for a combined entry")
	assert.Equal(t, model.TypeEntry{}, teOnlyID, "ResolveByID should return zero-value TypeEntry when resolving with only ID for a combined entry")

	// 8. Attempt to resolve with only Name.
	teOnlyName, okOnlyName := registry.ResolveByID("", name)
	assert.False(t, okOnlyName, "ResolveByID should return ok=false when resolving with only Name for a combined entry")
	assert.Equal(t, model.TypeEntry{}, teOnlyName, "ResolveByID should return zero-value TypeEntry when resolving with only Name for a combined entry")

	// 9. Ensure that the type cannot be retrieved unless both ID and Name are provided together.
	// Attempt with incorrect ID and correct Name.
	incorrectID := "IncorrectID"
	teIncorrectID, okIncorrectID := registry.ResolveByID(incorrectID, name)
	assert.False(t, okIncorrectID, "ResolveByID should return ok=false when resolving with incorrect ID and correct Name")
	assert.Equal(t, model.TypeEntry{}, teIncorrectID, "ResolveByID should return zero-value TypeEntry for incorrect ID and correct Name combination")

	// Attempt with correct ID and incorrect Name.
	incorrectName := "IncorrectName"
	teIncorrectName, okIncorrectName := registry.ResolveByID(id, incorrectName)
	assert.False(t, okIncorrectName, "ResolveByID should return ok=false when resolving with correct ID and incorrect Name")
	assert.Equal(t, model.TypeEntry{}, teIncorrectName, "ResolveByID should return zero-value TypeEntry for correct ID and incorrect Name combination")

	// Attempt with both incorrect ID and Name.
	teBothIncorrect, okBothIncorrect := registry.ResolveByID("WrongID", "WrongName")
	assert.False(t, okBothIncorrect, "ResolveByID should return ok=false when resolving with incorrect ID and Name")
	assert.Equal(t, model.TypeEntry{}, teBothIncorrect, "ResolveByID should return zero-value TypeEntry for incorrect ID and Name combination")

	// 10. Ensure that other registrations do not interfere.
	// Register another type with a different ID and Name.
	anotherID := "AnotherID"
	anotherName := "AnotherName"
	err = registry.RegisterIDandNameLookup(anotherID, anotherName, &AnotherStruct{})
	assert.NoError(t, err, "RegisterIDandNameLookup should not return an error for another valid ID and Name")

	// Resolve the new registration correctly.
	teAnother, okAnother := registry.ResolveByID(anotherID, anotherName)
	assert.True(t, okAnother, "ResolveByID should return ok=true for the second registered ID and Name combination")
	assert.Equal(t, reflect.TypeOf(AnotherStruct{}), teAnother.Model, "TypeEntry.Model should match AnotherStruct for the second registration")

	// Ensure that resolving the first ID and Name still works correctly.
	teFirst, okFirst := registry.ResolveByID(id, name)
	assert.True(t, okFirst, "ResolveByID should still return ok=true for the first registered ID and Name combination")
	assert.Equal(t, te, teFirst, "ResolveByID should return the correct TypeEntry for the first registration")
}

// TestRegisterIDandNameLookupWithEmptyIDOrName verifies that RegisterIDandNameLookup returns an error
// when either the ID or the Name is an empty string.
func TestRegisterIDandNameLookupWithEmptyIDOrName(t *testing.T) {
	// 1. Create a TypeRegistryImpl instance.
	registry := types.NewRegistry()

	// 2. Attempt to register a type with an empty ID and a valid Name.
	emptyID := ""
	validName := "SomeName"
	err := registry.RegisterIDandNameLookup(emptyID, validName, &SomeStruct{})
	assert.Error(t, err, "RegisterIDandNameLookup should return an error when ID is empty")
	assert.EqualError(t, err, "both id and name must be non empty strings", "Error message should indicate that both ID and Name must be non-empty")

	// 3. Attempt to register a type with a valid ID and an empty Name.
	validID := "SomeID"
	emptyName := ""
	err = registry.RegisterIDandNameLookup(validID, emptyName, &SomeStruct{})
	assert.Error(t, err, "RegisterIDandNameLookup should return an error when Name is empty")
	assert.EqualError(t, err, "both id and name must be non empty strings", "Error message should indicate that both ID and Name must be non-empty")

	// 4. Attempt to resolve using the empty ID and valid Name.
	// Since the registration with empty ID should have failed, this should not retrieve any entry.
	teEmptyID, okEmptyID := registry.ResolveByID(emptyID, validName)
	assert.False(t, okEmptyID, "ResolveByID should return ok=false when resolving with empty ID and valid Name due to failed registration")
	assert.Equal(t, model.TypeEntry{}, teEmptyID, "ResolveByID should return zero-value TypeEntry when resolving with empty ID and valid Name")

	// 5. Attempt to resolve using the valid ID and empty Name.
	// Similarly, this should not retrieve any entry.
	teEmptyName, okEmptyName := registry.ResolveByID(validID, emptyName)
	assert.False(t, okEmptyName, "ResolveByID should return ok=false when resolving with valid ID and empty Name due to failed registration")
	assert.Equal(t, model.TypeEntry{}, teEmptyName, "ResolveByID should return zero-value TypeEntry when resolving with valid ID and empty Name")

	// 6. Ensure that no entries were added to the 'id_names' map due to failed registrations.
	// Attempt to resolve with both empty ID and empty Name.
	teBothEmpty, okBothEmpty := registry.ResolveByID(emptyID, emptyName)
	assert.False(t, okBothEmpty, "ResolveByID should return ok=false when resolving with both empty ID and Name")
	assert.Equal(t, model.TypeEntry{}, teBothEmpty, "ResolveByID should return zero-value TypeEntry when resolving with both empty ID and Name")

	// 7. Additionally, verify that a valid registration still works as expected.
	// Register a type with both valid ID and Name.
	err = registry.RegisterIDandNameLookup(validID, validName, &SomeStruct{})
	assert.NoError(t, err, "RegisterIDandNameLookup should succeed with valid ID and Name")

	// Resolve with both valid ID and Name.
	teValid, okValid := registry.ResolveByID(validID, validName)
	assert.True(t, okValid, "ResolveByID should return ok=true for valid ID and Name combination")
	assert.NotNil(t, teValid, "TypeEntry should not be nil for valid ID and Name combination")
	assert.Equal(t, validName, teValid.Name, "TypeEntry.Name should match the registered Name")
	assert.Equal(t, reflect.TypeOf(SomeStruct{}), teValid.Model, "TypeEntry.Model should match SomeStruct for the registered ID and Name")
}

// TestRegisterIDandNameLookupDuplicateCombination verifies that attempting to register a duplicate
// ID and Name combination using RegisterIDandNameLookup returns an error and that the original
// registration remains intact.
func TestRegisterIDandNameLookupDuplicateCombination(t *testing.T) {
	// Define the ID, Name, and the types to be registered.
	id := "SomeID"
	name := "SomeName"

	// 1. Create a TypeRegistryImpl instance.
	registry := types.NewRegistry()

	// 2. Register the first type (SomeStruct) with the specified ID and Name.
	err := registry.RegisterIDandNameLookup(id, name, &SomeStruct{})
	assert.NoError(t, err, "First RegisterIDandNameLookup should succeed without error")

	// 3. Attempt to register the second type (AnotherStruct) with the same ID and Name.
	err = registry.RegisterIDandNameLookup(id, name, &AnotherStruct{})
	assert.Error(t, err, "Second RegisterIDandNameLookup with duplicate ID and Name should return an error")

	// 4. Assert that the error message matches the expected message.
	assert.EqualError(t, err, "id+name already registered: SomeID#SomeName")

	// 5. Resolve by ID and Name to ensure the original type (SomeStruct) is still registered.
	te, ok := registry.ResolveByID(id, name)
	assert.True(t, ok, "ResolveByID should find the original TypeEntry for ID 'SomeID' and Name 'SomeName'")
	assert.NotNil(t, te, "TypeEntry should not be nil for registered ID and Name combination")

	// 6. Assert that the retrieved TypeEntry corresponds to SomeStruct, not AnotherStruct.
	expectedType := reflect.TypeOf(SomeStruct{})
	assert.Equal(t, expectedType, te.Model, "TypeEntry.Model should match SomeStruct after duplicate registration attempt")
	assert.Equal(t, name, te.Name, "TypeEntry.Name should match the registered Name 'SomeName'")

	// 7. Ensure that the duplicate registration did not overwrite the original entry.
	// Attempting to resolve with the same ID and Name should still return SomeStruct.
	teDuplicate, okDuplicate := registry.ResolveByID(id, name)
	assert.True(t, okDuplicate, "ResolveByID should still find the original TypeEntry after duplicate registration attempt")
	assert.Equal(t, te, teDuplicate, "ResolveByID should return the same TypeEntry as before the duplicate registration")
}

// TestRegisterWithMetadata verifies that registering a type with metadata correctly stores the metadata
// within the TypeEntry and that it can be accurately retrieved.
func TestRegisterWithMetadata(t *testing.T) {
	// 1. Create a TypeRegistryImpl instance.
	registry := types.NewRegistry()

	// 2. Define the metadata to be associated with the type.
	metadata := map[string]string{
		"Description": "This is a sample struct for testing metadata storage.",
		"Version":     "1",
		"Tags":        "test,metadata,struct",
	}

	// 3. Register the type with a specific name using RegisterNameLookup, including the metadata.
	name := "SomeNameWithMeta"
	err := registry.RegisterNameLookup(name, &SomeStruct{}, metadata)
	assert.NoError(t, err, "RegisterNameLookup should not return an error when registering with metadata")

	// 4. Retrieve the TypeEntry using Get.
	te, ok := registry.Get(name)
	assert.True(t, ok, "Get should return ok=true for the registered name")
	assert.NotNil(t, te, "TypeEntry should not be nil for the registered name")

	// 5. Assert that the metadata in TypeEntry matches the provided metadata.
	assert.NotNil(t, te.Meta, "TypeEntry.Meta should not be nil when metadata is provided during registration")
	assert.Equal(t, metadata, te.Meta, "TypeEntry.Meta should match the provided metadata")

	// 6. Additionally, verify the metadata's contents for deeper validation.
	retrievedMeta := te.Meta
	assert.Equal(t, metadata["Description"], retrievedMeta["Description"], "Metadata.Description should match")
	assert.Equal(t, metadata["Version"], retrievedMeta["Version"], "Metadata.Version should match")
	assert.Equal(t, metadata["Tags"], retrievedMeta["Tags"], "Metadata.Tags should match")

	// 7. Optionally, verify that metadata is also correctly associated when using other registration methods.

	// Example: Register with ID and Name, including metadata.
	id := "SomeIDWithMeta"
	err = registry.RegisterIDandNameLookup(id, name, &SomeStruct{}, metadata)
	assert.NoError(t, err, "RegisterIDandNameLookup should not return an error when registering with metadata")

	// Retrieve using ResolveByID with both ID and Name.
	teIDName, okIDName := registry.ResolveByID(id, name)
	assert.True(t, okIDName, "ResolveByID should return ok=true for the registered ID and Name combination")
	assert.NotNil(t, teIDName, "TypeEntry should not be nil for the registered ID and Name combination")
	assert.Equal(t, metadata, teIDName.Meta, "TypeEntry.Meta should match the provided metadata for ID and Name registration")

	// Assert that the metadata fields are correctly set.
	retrievedMetaIDName := teIDName.Meta
	assert.Equal(t, metadata["Description"], retrievedMetaIDName["Description"], "Metadata.Description should match for ID and Name registration")
	assert.Equal(t, metadata["Version"], retrievedMetaIDName["Version"], "Metadata.Version should match for ID and Name registration")
	assert.Equal(t, metadata["Tags"], retrievedMetaIDName["Tags"], "Metadata.Tags should match for ID and Name registration")

	// 8. Verify that resolving without metadata (e.g., registering without metadata) behaves correctly.
	// Register another type without metadata.
	nameNoMeta := "SomeNameWithoutMeta"
	err = registry.RegisterNameLookup(nameNoMeta, &SomeStruct{})
	assert.NoError(t, err, "RegisterNameLookup should not return an error when registering without metadata")

	// Retrieve the TypeEntry.
	teNoMeta, okNoMeta := registry.Get(nameNoMeta)
	assert.True(t, okNoMeta, "Get should return ok=true for the registered name without metadata")
	assert.NotNil(t, teNoMeta, "TypeEntry should not be nil for the registered name without metadata")
	assert.Nil(t, teNoMeta.Meta, "TypeEntry.Meta should be nil when no metadata was provided during registration")

	// 9. Ensure that metadata does not interfere with type resolution.
	// Attempt to resolve the type with correct identifiers.
	teResolved, okResolved := registry.ResolveByID(id, name)
	assert.True(t, okResolved, "ResolveByID should successfully retrieve the TypeEntry with both ID and Name")
	assert.Equal(t, teIDName, teResolved, "ResolveByID should return the correct TypeEntry with both ID and Name")

	// Attempt to resolve with incorrect metadata.
	// Since metadata is part of the TypeEntry, attempting to change it externally should not affect the registry.
	// This is more of a conceptual assertion, as Go's type system prevents altering internal state directly.
	// However, ensuring that the metadata remains unchanged can be considered.
	assert.Equal(t, metadata, teResolved.Meta, "Metadata should remain unchanged and match the originally registered metadata")
}

// Define multiple structs for testing.
type SomeStructA struct{}
type SomeStructB struct{}
type SomeStructC struct{}

// ResolverA handles resolution for specific IDs.
type ResolverA struct{}

// Resolve method implementation for ResolverA.
func (r *ResolverA) ResolveByID(id string, name string) (model.TypeEntry, bool) {
	switch id {
	case "A1":
		return model.TypeEntry{
			Name:  "NameA1",
			Model: reflect.TypeOf(SomeStructA{}),
			Meta: map[string]string{
				"Description": "ResolverA - Type SomeStructA",
				"Version":     "1",
			},
		}, true
	case "A2":
		return model.TypeEntry{
			Name:  "NameA2",
			Model: reflect.TypeOf(SomeStructB{}),
			Meta: map[string]string{
				"Description": "ResolverA - Type SomeStructB",
				"Version":     "2",
			},
		}, true
	default:
		return model.TypeEntry{}, false
	}
}

// ResolverB handles resolution for specific IDs.
type ResolverB struct{}

// Resolve method implementation for ResolverB.
func (r *ResolverB) ResolveByID(id string, name string) (model.TypeEntry, bool) {
	switch id {
	case "B1":
		return model.TypeEntry{
			Name:  "NameB1",
			Model: reflect.TypeOf(SomeStructC{}),
			Meta: map[string]string{
				"Description": "ResolverB - Type SomeStructC",
				"Version":     "1",
			},
		}, true
	case "B2":
		// Simulate a resolver that cannot handle ID "B2"
		return model.TypeEntry{}, false
	default:
		return model.TypeEntry{}, false
	}
}

// TestRegisterMultipleResolversAndResolveByID verifies that multiple resolvers can be registered
// and that ResolveByID correctly retrieves TypeEntries using these resolvers.
func TestRegisterMultipleResolversAndResolveByID(t *testing.T) {
	// 1. Create a TypeRegistryImpl instance.
	registry := types.NewRegistry()

	// 2. Create resolver instances.
	resolverA := &ResolverA{}
	resolverB := &ResolverB{}

	// 3. Register resolvers with the registry.
	registry.RegisterResolver(resolverA)
	registry.RegisterResolver(resolverB)

	// 4. Define test cases with various ID and Name combinations.
	testCases := []struct {
		ID          string
		Name        string
		ExpectedOK  bool
		ExpectedTE  model.TypeEntry
		Description string
	}{
		{
			ID:         "A1",
			Name:       "NameA1",
			ExpectedOK: true,
			ExpectedTE: model.TypeEntry{
				Name:  "NameA1",
				Model: reflect.TypeOf(SomeStructA{}),
				Meta: map[string]string{
					"Description": "ResolverA - Type SomeStructA",
					"Version":     "1",
				},
			},
			Description: "ResolverA should resolve ID A1 to SomeStructA",
		},
		{
			ID:         "A2",
			Name:       "NameA2",
			ExpectedOK: true,
			ExpectedTE: model.TypeEntry{
				Name:  "NameA2",
				Model: reflect.TypeOf(SomeStructB{}),
				Meta: map[string]string{
					"Description": "ResolverA - Type SomeStructB",
					"Version":     "2",
				},
			},
			Description: "ResolverA should resolve ID A2 to SomeStructB",
		},
		{
			ID:         "B1",
			Name:       "NameB1",
			ExpectedOK: true,
			ExpectedTE: model.TypeEntry{
				Name:  "NameB1",
				Model: reflect.TypeOf(SomeStructC{}),
				Meta: map[string]string{
					"Description": "ResolverB - Type SomeStructC",
					"Version":     "1",
				},
			},
			Description: "ResolverB should resolve ID B1 to SomeStructC",
		},
		{
			ID:          "B2",
			Name:        "NameB2",
			ExpectedOK:  false,
			ExpectedTE:  model.TypeEntry{},
			Description: "ResolverB should not resolve ID B2",
		},
		{
			ID:          "C1",
			Name:        "NameC1",
			ExpectedOK:  false,
			ExpectedTE:  model.TypeEntry{},
			Description: "No resolver should resolve ID C1",
		},
	}

	// 5. Perform lookups using ResolveByID.
	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			te, ok := registry.ResolveByID(tc.ID, tc.Name)
			if tc.ExpectedOK {
				assert.True(t, ok, "Expected ResolveByID to succeed for %s", tc.Description)
				assert.Equal(t, tc.ExpectedTE.Name, te.Name, "TypeEntry.Name mismatch for %s", tc.Description)
				assert.Equal(t, tc.ExpectedTE.Model, te.Model, "TypeEntry.Model mismatch for %s", tc.Description)
				assert.Equal(t, tc.ExpectedTE.Meta, te.Meta, "TypeEntry.Meta mismatch for %s", tc.Description)
			} else {
				assert.False(t, ok, "Expected ResolveByID to fail for %s", tc.Description)
				assert.Equal(t, tc.ExpectedTE, te, "Expected zero-value TypeEntry for %s", tc.Description)
			}
		})
	}
}
