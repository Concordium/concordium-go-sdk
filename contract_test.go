package concordium

import (
	"testing"
)

var (
	testModuleRef      = ModuleRef("1d40f9366f6fcb4586ac8e09ed391b5832cfd752fb63ee7bd38da0f3e77c4204")
	testModuleRefBytes = []byte{
		29, 64, 249, 54, 111, 111, 203, 69, 134, 172, 142, 9, 237, 57, 27, 88,
		50, 207, 215, 82, 251, 99, 238, 123, 211, 141, 160, 243, 231, 124, 66, 4,
	}

	testContractName = ContractName("foo")

	testInitName        = NewInitName(testContractName)
	testInitNameBeBytes = []byte{0, 8, 105, 110, 105, 116, 95, 102, 111, 111}
	testInitNameLeBytes = []byte{8, 0, 105, 110, 105, 116, 95, 102, 111, 111}

	testReceiveName        = NewReceiveName(testContractName, "bar")
	testReceiveNameBeBytes = []byte{0, 7, 102, 111, 111, 46, 98, 97, 114}
	testReceiveNameLeBytes = []byte{7, 0, 102, 111, 111, 46, 98, 97, 114}
)

func TestModuleRef_Serialize(t *testing.T) {
	testSerialize(t, &testModuleRef, testModuleRefBytes)
}

func TestModuleRef_Deserialize(t *testing.T) {
	m := ModuleRef("")
	testDeserialize(t, &m, &testModuleRef, testModuleRefBytes)
}

func TestNewInitName(t *testing.T) {
	v := InitName("init_foo")
	a := NewInitName(testContractName)
	if a != v {
		t.Errorf("NewInitName() got = %v, w %v", a, v)
	}
}

func TestInitName_Serialize(t *testing.T) {
	testSerialize(t, &testInitName, testInitNameBeBytes)
}

func TestInitName_SerializeModel(t *testing.T) {
	testSerializeModel(t, &testInitName, testInitNameLeBytes)
}

func TestInitName_DeserializeModel(t *testing.T) {
	n := InitName("")
	testDeserializeModel(t, &n, &testInitName, testInitNameLeBytes)
}

func TestNewReceiveName(t *testing.T) {
	v := ReceiveName("foo.bar")
	a := NewReceiveName(testContractName, "bar")
	if a != v {
		t.Errorf("NewReceiveName() got = %v, w %v", a, v)
	}
}

func TestReceiveName_Serialize(t *testing.T) {
	testSerialize(t, &testReceiveName, testReceiveNameBeBytes)
}

func TestReceiveName_SerializeModel(t *testing.T) {
	testSerializeModel(t, &testReceiveName, testReceiveNameLeBytes)
}

func TestReceiveName_DeserializeModel(t *testing.T) {
	n := ReceiveName("")
	testDeserializeModel(t, &n, &testReceiveName, testReceiveNameLeBytes)
}
