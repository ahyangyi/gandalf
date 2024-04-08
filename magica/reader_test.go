package magica

import (
	"bytes"
	"encoding/binary"
	"github.com/ahyangyi/gandalf/geometry"
	"github.com/ahyangyi/gandalf/magica/types"
	"strings"
	"testing"
)

func TestIsHeaderValid(t *testing.T) {
	var testCases = []struct {
		input    string
		expected bool
	}{
		{"VOX ", true},
		{"BLAH", false},
		{"A", false},
		{"ABCDE", false},
		{"VOX AAA", true},
	}

	for _, testCase := range testCases {
		reader := strings.NewReader(testCase.input)
		result := isHeaderValid(reader)
		if result != testCase.expected {
			t.Errorf("Magic string %s expected %v, got %v", testCase.input, testCase.expected, result)
		}

		expectedLength := len(testCase.input) - 4
		if expectedLength < 0 {
			expectedLength = 0
		}

		if reader.Len() != expectedLength {
			t.Errorf("Did not read 4 bytes of string %s, %d remaining of %d", testCase.input, reader.Len(), len(testCase.input))
		}
	}
}

func TestGetSizeFromChunk(t *testing.T) {
	var testCases = []struct {
		input    []byte
		expected types.Size
	}{
		{getSizedByteSlice(12, []byte{1, 0, 0, 0, 2, 0, 0, 0, 3, 0, 0, 0}),
			types.Size{X: 1, Y: 2, Z: 3}},
		{getSizedByteSlice(16, []byte{3, 0, 0, 0, 5, 0, 0, 0, 1, 0, 0, 0, 2, 4, 6, 8}),
			types.Size{X: 3, Y: 5, Z: 1}},
	}

	for _, testCase := range testCases {
		data, _ := getChunkData(bytes.NewBuffer(testCase.input), 0)
		reader := types.GetReader(data)
		result := reader.GetSize()

		if result != testCase.expected {
			t.Errorf("Byte array %v expected %v, got %v", testCase.input, testCase.expected, result)
		}
	}
}

func TestGetPointDataFromChunk(t *testing.T) {
	var testCases = []struct {
		input    []byte
		expected []geometry.PointWithColour
	}{
		{input: getSizedByteSlice(4, []byte{1, 2, 3, 64}),
			expected: []geometry.PointWithColour{{Point: geometry.Point{X: 1, Y: 2, Z: 3}, Colour: 64}}},
		{input: getSizedByteSlice(8, []byte{1, 2, 3, 64, 4, 5, 6, 128}),
			expected: []geometry.PointWithColour{
				{Point: geometry.Point{X: 1, Y: 2, Z: 3}, Colour: 64},
				{Point: geometry.Point{X: 4, Y: 5, Z: 6}, Colour: 128},
			}},
		{getSizedByteSlice(5, []byte{1, 2, 3, 4, 5}),
			[]geometry.PointWithColour{{Point: geometry.Point{X: 1, Y: 2, Z: 3}, Colour: 4}}},
	}

	for _, testCase := range testCases {
		data, _ := getChunkData(bytes.NewBuffer(testCase.input), 0)
		reader := types.GetReader(data)
		result := reader.GetPointData()

		if !arePointWithColourSlicesEqual(result, testCase.expected) {
			t.Errorf("Byte array %v expected %v, got %v", testCase.input, testCase.expected, result)
		}
	}
}

func Test_getPaletteDataFromChunk(t *testing.T) {
	testCases := []struct {
		input    []byte
		expected []byte
	}{
		{getSizedByteSlice(4, []byte{1, 2, 3, 4}), []byte{1, 2, 3, 4}},
	}

	for _, tc := range testCases {
		data, _ := getChunkData(bytes.NewBuffer(tc.input), 0)
		reader := types.GetReader(data)

		if result := reader.GetPalette(); !bytes.Equal(result, tc.expected) {
			t.Errorf("Get palette for %v expected %v, got %v", tc.input, tc.expected, result)
		}
	}
}

func TestGetRawVoxels(t *testing.T) {
	testData := []byte{'V', 'O', 'X', ' '}
	testData = append(testData, []byte{150, 0, 0, 0}...)
	testData = append(testData, []byte{'S', 'I', 'Z', 'E'}...)
	testData = append(testData, getSizedByteSlice(12, []byte{2, 0, 0, 0, 2, 0, 0, 0, 2, 0, 0, 0})...)
	testData = append(testData, []byte{'X', 'Y', 'Z', 'I'}...)
	testData = append(testData, getSizedByteSlice(12, []byte{1, 1, 1, 255, 0, 1, 1, 1, 20, 31, 11, 1})...)
	testData = append(testData, []byte{'U', 'N', 'K', ' '}...)
	testData = append(testData, getSizedByteSlice(2, []byte{1, 2})...)

	reader := bytes.NewReader(testData)
	result, err := GetMagicaVoxelObject(reader, []int{})

	if err != nil {
		t.Errorf("Encountered error %v", err)
	}

	testMagicaVoxelObject(result, t)
}

func testMagicaVoxelObject(object VoxelObject, t *testing.T) {
	if len(object.Voxels) != 2 {
		t.Error("x dimension not correctly sized")
		return
	}
	if len(object.Voxels[0]) != 2 {
		t.Error("y dimension not correctly sized")
		return
	}

	if len(object.Voxels[0][0]) != 2 {
		t.Error("z dimension not correctly sized")
		return
	}

	if object.Voxels[1][1][1] != 255 {
		t.Error("Point at (1,1,1) was not set")
	}

	if object.Voxels[0][1][1] != 1 {
		t.Error("Point at (0,1,1) was not set")
	}

	if object.Voxels[0][0][1] != 0 {
		t.Error("Point at (0,0,0) was not left unset")
	}
}

func arePointWithColourSlicesEqual(a []geometry.PointWithColour, b []geometry.PointWithColour) bool {
	if len(a) != len(b) {
		return false
	}

	for i, p := range a {
		if p != b[i] {
			return false
		}
	}

	return true
}

func getSizedByteSlice(size int64, slice []byte) []byte {
	result := make([]byte, 8)
	binary.LittleEndian.PutUint64(result, uint64(size))
	result = append(result, slice...)
	return result
}

func TestMultipartVoxelObject(t *testing.T) {
	object, err := FromFile("testdata/test_cube_joined.vox")
	if err != nil {
		t.Errorf("Could not read object: %v", err)
	}

	if object.Size.X != 4 {
		t.Errorf("Object size X - expected %d, got %d", 4, object.Size.X)
	}
}
