package tag

import "testing"

type test struct {
	tag          Tag
	trackingTags map[string]interface{}
	result       bool
}

func TestReachedFreezer(t *testing.T) {

	destination := "FREEZER"

	testCases := []test{
		{
			tag:          Tag{Epc: "TEST0", LocationHistory: []LocationHistory{{Location: destination}}},
			trackingTags: map[string]interface{}{"TEST0": nil},
			result:       true,
		},
		{
			tag:          Tag{Epc: "TEST1", LocationHistory: []LocationHistory{{Location: destination}}},
			trackingTags: map[string]interface{}{"TEST0": nil, "TEST2": nil},
			result:       false,
		},
	}

	for _, value := range testCases {

		if ReachedFreezer(value.tag, destination, value.trackingTags) != value.result {
			t.Error("Returned value doesn't match with expected result.")
		}
	}

}
