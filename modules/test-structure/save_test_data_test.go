package test_structure

import (
	"testing"
	"io/ioutil"
	"github.com/stretchr/testify/assert"
	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

type testData struct {
	Foo string
	Bar bool
	Baz map[string]interface{}
}

func TestSaveAndLoadTestData(t *testing.T) {
	t.Parallel()

	isTestDataPresent := IsTestDataPresent(t, "/file/that/does/not/exist")
	assert.False(t, isTestDataPresent, "Expected no test data would be present because no test data file exists.")

	tmpFile, err := ioutil.TempFile("", "save-and-load-test-data")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	expectedData := testData{
		Foo: "foo",
		Bar: true,
		Baz: map[string]interface{}{"abc": "def", "ghi": 1.0, "klm": false},
	}

	isTestDataPresent = IsTestDataPresent(t, tmpFile.Name())
	assert.False(t, isTestDataPresent, "Expected no test data would be present because file exists but no data has been written yet.")

	SaveTestData(t, tmpFile.Name(), expectedData)

	isTestDataPresent = IsTestDataPresent(t, tmpFile.Name())
	assert.True(t, isTestDataPresent, "Expected test data would be present because file exists and data has been written to file.")

	actualData := testData{}
	LoadTestData(t, tmpFile.Name(), &actualData)
	assert.Equal(t, expectedData, actualData)

	CleanupTestData(t, tmpFile.Name())
	assert.False(t, files.FileExists(tmpFile.Name()))
}

func TestIsEmptyJson(t *testing.T) {
	t.Parallel()

	var jsonValue []byte
	var isEmpty bool

	jsonValue = []byte("null")
	isEmpty = isEmptyJson(t, jsonValue)
	assert.True(t, isEmpty, `The JSON literal "null" should be treated as an empty value.`)

	jsonValue = []byte("false")
	isEmpty = isEmptyJson(t, jsonValue)
	assert.True(t, isEmpty, `The JSON literal "false" should be treated as an empty value.`)

	jsonValue = []byte("true")
	isEmpty = isEmptyJson(t, jsonValue)
	assert.False(t, isEmpty, `The JSON literal "true" should be treated as a non-empty value.`)

	jsonValue = []byte("0")
	isEmpty = isEmptyJson(t, jsonValue)
	assert.True(t, isEmpty, `The JSON literal "0" should be treated as an empty value.`)

	jsonValue = []byte("1")
	isEmpty = isEmptyJson(t, jsonValue)
	assert.False(t, isEmpty, `The JSON literal "1" should be treated as a non-empty value.`)

	jsonValue = []byte("{}")
	isEmpty = isEmptyJson(t, jsonValue)
	assert.True(t, isEmpty, `The JSON value "{}" should be treated as an empty value.`)

	jsonValue = []byte(`{ "key": "val" }`)
	isEmpty = isEmptyJson(t, jsonValue)
	assert.False(t, isEmpty, `The JSON value { "key": "val" } should be treated as a non-empty value.`)

	jsonValue = []byte(`[]`)
	isEmpty = isEmptyJson(t, jsonValue)
	assert.True(t, isEmpty, `The JSON value "[]" should be treated as an empty value.`)

	jsonValue = []byte(`[{ "key": "val" }]`)
	isEmpty = isEmptyJson(t, jsonValue)
	assert.False(t, isEmpty, `The JSON value [{ "key": "val" }] should be treated as a non-empty value.`)
}

func TestSaveAndLoadTerraformOptions(t *testing.T) {
	t.Parallel()

	tmpFolder, err := ioutil.TempDir("", "save-and-load-terratest-options")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	expectedData := &terraform.Options{
		TerraformDir: "/abc/def/ghi",
		Vars: map[string]interface{}{},
	}
	SaveTerraformOptions(t, tmpFolder, expectedData)

	actualData := LoadTerraformOptions(t, tmpFolder)
	assert.Equal(t, expectedData, actualData)
}

func TestSaveAndLoadAmiId(t *testing.T) {
	t.Parallel()

	tmpFolder, err := ioutil.TempDir("", "save-and-load-ami-id")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	expectedData := "ami-abcd1234"
	SaveAmiId(t, tmpFolder, expectedData)

	actualData := LoadAmiId(t, tmpFolder)
	assert.Equal(t, expectedData, actualData)
}

func TestSaveAndLoadNamedStrings(t *testing.T) {
	t.Parallel()

	tmpFolder, err := ioutil.TempDir("", "save-and-load-named-strings")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	name1 := "test-ami"
	expectedData1 := "ami-abcd1234"

	name2 := "test-ami2"
	expectedData2 := "ami-xyz98765"

	SaveString(t, tmpFolder, name1, expectedData1)
	SaveString(t, tmpFolder, name2, expectedData2)

	actualData1 := LoadString(t, tmpFolder, name1)
	actualData2 := LoadString(t, tmpFolder, name2)

	assert.Equal(t, expectedData1, actualData1)
	assert.Equal(t, expectedData2, actualData2)
}

func TestSaveDuplicateTestData(t *testing.T) {
	t.Parallel()

	tmpFolder, err := ioutil.TempDir("", "save-and-load-duplicate-test-data")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	name := "hello-world"
	val1 := "hello world"
	val2 := "buenos dias, mundo"

	SaveString(t, tmpFolder, name, val1)
	SaveString(t, tmpFolder, name, val2)

	actualVal := LoadString(t, tmpFolder, name)

	assert.Equal(t, val2, actualVal, "Actual test data should use overwritten values")
}

func TestSaveAndLoadNamedInts(t *testing.T) {
	t.Parallel()

	tmpFolder, err := ioutil.TempDir("", "save-and-load-named-ints")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	name1 := "test-int1"
	expectedData1 := 23842834

	name2 := "test-int2"
	expectedData2 := 52

	SaveInt(t, tmpFolder, name1, expectedData1)
	SaveInt(t, tmpFolder, name2, expectedData2)

	actualData1 := LoadInt(t, tmpFolder, name1)
	actualData2 := LoadInt(t, tmpFolder, name2)

	assert.Equal(t, expectedData1, actualData1)
	assert.Equal(t, expectedData2, actualData2)
}