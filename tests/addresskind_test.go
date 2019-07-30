package tests

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestAddressKind(t *testing.T) {

	t.Run("marshal", testMarshal)
	t.Run("unmarshal", testUnmarshal)
}

func testMarshal(t *testing.T) {

	for i, test := range tests {
		t.Logf("test %d", i)

		b, err := json.MarshalIndent(test.obj, "", "  ")
		if err != nil {
			t.Error(err)
		}
		if string(b) != test.str {
			t.Errorf("test %d\nexpected:\n%s\nactual:\n%s\n", i, test.str, string(b))
		}
	}
}

func testUnmarshal(t *testing.T) {

	for i, test := range tests {
		t.Logf("test %d", i)

		var info Info
		err := json.Unmarshal([]byte(test.str), &info)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(info, test.obj) {
			t.Errorf("test %d\nexpected:\n%v\nactual:\n%v\n", i, test.obj, info)
		}
	}
}

var tests = []struct {
	obj Info
	str string
}{
	{Info{
		Name:        "Homer",
		AddressKind: AddressKindHome,
		Address: Home{
			HouseNumber: "123",
			Street:      "Home Street",
			City:        "Seattle",
			State:       "WA",
			PostalCode:  "98101",
		},
		Country: "USA",
	},
		`{
  "name": "Homer",
  "addressKind": "Home",
  "address": {
    "houseNumber": "123",
    "street": "Home Street",
    "city": "Seattle",
    "state": "WA",
    "postalCode": "98101"
  },
  "country": "USA"
}`},
	{Info{
		Name:        "Buro",
		AddressKind: AddressKindOffice,
		Address: Office{
			OfficeNumber: "987",
			Street:       "Office Park",
			City:         "Bellevue",
			State:        "WA",
			PostalCode:   "98004",
		},
		Country: "USA",
	},
		`{
  "name": "Buro",
  "addressKind": "Office",
  "address": {
    "officeNumber": "987",
    "street": "Office Park",
    "city": "Bellevue",
    "state": "WA",
    "postalCode": "98004"
  },
  "country": "USA"
}`},
	{Info{
		Name:        "Going Postal",
		AddressKind: AddressKindPostalBox,
		Address: PostalBox{
			POBoxNumber: "30125",
			City:        "Seattle",
			State:       "WA",
			PostalCode:  "98113-0125",
		},
		Country: "USA",
	},
		`{
  "name": "Going Postal",
  "addressKind": "PostalBox",
  "address": {
    "poBoxNumber": "30125",
    "city": "Seattle",
    "state": "WA",
    "postalCode": "98113-0125"
  },
  "country": "USA"
}`},
}

// Info -- common infomation of business object
type Info struct {
	Name        string      `json:"name"`
	AddressKind AddressKind `json:"addressKind,omitempty"`
	Address     Address     `json:"address"`
	Country     string      `json:"country"`
}

// UnmarshalJSON - -Info structure
func (i *Info) UnmarshalJSON(b []byte) error {

	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	err = json.Unmarshal(*objMap["name"], &i.Name)
	if err != nil {
		return err
	}

	err = json.Unmarshal(*objMap["addressKind"], &i.AddressKind)
	if err != nil {
		return err
	}

	err = json.Unmarshal(*objMap["country"], &i.Country)
	if err != nil {
		return err
	}

	switch i.AddressKind {
	case AddressKindHome:

		var home Home
		err = json.Unmarshal(*objMap["address"], &home)
		if err != nil {
			return err
		}
		i.Address = home

	case AddressKindOffice:

		var office Office
		err = json.Unmarshal(*objMap["address"], &office)
		if err != nil {
			return err
		}
		i.Address = office

	case AddressKindPostalBox:

		var postalBox PostalBox
		err = json.Unmarshal(*objMap["address"], &postalBox)
		if err != nil {
			return err
		}
		i.Address = postalBox

	}

	return nil
}

// Address -- Oneof (Home|Office|PostalBox)
type Address interface {
	isAddressType()
}

// Home -- Private home location address
type Home struct {
	HouseNumber string `json:"houseNumber"`
	Street      string `json:"street"`
	City        string `json:"city"`
	State       string `json:"state"`
	PostalCode  string `json:"postalCode"`
}

func (h Home) isAddressType() {}

// Office -- Legal office location address
type Office struct {
	OfficeNumber string `json:"officeNumber"`
	Street       string `json:"street"`
	City         string `json:"city"`
	State        string `json:"state"`
	PostalCode   string `json:"postalCode"`
}

func (o Office) isAddressType() {}

// PostalBox -- PO Box Address
type PostalBox struct {
	POBoxNumber string `json:"poBoxNumber"`
	City        string `json:"city"`
	State       string `json:"state"`
	PostalCode  string `json:"postalCode"`
}

func (p PostalBox) isAddressType() {}
