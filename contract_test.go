package concordium

import (
	"testing"
)

var (
	testModuleRefString = "1d40f9366f6fcb4586ac8e09ed391b5832cfd752fb63ee7bd38da0f3e77c4204"
	testModuleRef       = MustNewModuleRef(testModuleRefString)

	testContractName = ContractName("foo")

	testInitName        = NewInitName(testContractName)
	testInitNameBeBytes = []byte{0, 8, 105, 110, 105, 116, 95, 102, 111, 111}
	testInitNameLeBytes = []byte{8, 0, 105, 110, 105, 116, 95, 102, 111, 111}

	testReceiveName        = NewReceiveName(testContractName, "bar")
	testReceiveNameBeBytes = []byte{0, 7, 102, 111, 111, 46, 98, 97, 114}
	testReceiveNameLeBytes = []byte{7, 0, 102, 111, 111, 46, 98, 97, 114}
)

func TestModuleRef_Serialize(t *testing.T) {
	testSerialize(t, &testModuleRef, testModuleRef[:])
}

func TestModuleRef_Deserialize(t *testing.T) {
	a := testModuleRef
	testDeserialize(t, &a, &testModuleRef, testModuleRef[:])
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
