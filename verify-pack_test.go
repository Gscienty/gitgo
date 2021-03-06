package gitgo

import (
	"io"
	"log"
	"os"
	"path"
	"testing"
)

func Test_VerifyPack(t *testing.T) {

	// only used for git verify-pack -v
	const expected = `fe89ee30bbcdfdf376beae530cc53f967012f31c commit 267 184 12
3ead3116d0378089f5ce61086354aac43e736b01 commit 243 170 196
1d833eb5b6c5369c0cb7a4a3e20ded237490145f commit 262 180 366
a7f92c920ce85f07a33f948aa4fa2548b270024f commit 250 172 546
97eed02ebe122df8fdd853c1215d8775f3d9f1a1 commit 190 132 718
d22fc8a57073fdecae2001d00aff921440d3aabd tree   121 115 850
df891299372c34b57e41cfc50a0113e2afac3210 tree   25 37 965 1 d22fc8a57073fdecae2001d00aff921440d3aabd
af6e4fe91a8f9a0f3c03cbec9e1d2aac47345d67 blob   18 23 1002
6b32b1ac731898894c403f6b621bdda167ab8d7c blob   1645 700 1025
7147f43ae01c9f04a78d6e80544ed84def06e958 blob   1824 697 1725
05d3cc770bd3524cc25d47e083d8942ad25033f0 blob   16 28 2422 1 7147f43ae01c9f04a78d6e80544ed84def06e958
c3b8133617bbdb72e237b0f163fade7fbf1f0c18 blob   381 317 2450 2 05d3cc770bd3524cc25d47e083d8942ad25033f0
8264d7bcc297e15c452a7aef3a2e40934762b7e3 tree   25 38 2767 1 d22fc8a57073fdecae2001d00aff921440d3aabd
254671773e8cd91e07e36546c9a2d9c27e8dfeec tree   121 115 2805
ba74813270ff557c4a5d1be0562a141bbee4d3e6 blob   16 28 2920 1 6b32b1ac731898894c403f6b621bdda167ab8d7c
b45377f6daf59a4cec9e8de64f5df1533a7994cd blob   10 21 2948 1 7147f43ae01c9f04a78d6e80544ed84def06e958
9de6c72106b169990a83ce7090c7cad84b6b506b tree   38 49 2969
non delta: 11 objects
chain length = 1: 5 objects
chain length = 2: 1 object
.git/objects/pack/pack-d310969c4ba0ebfe725685fa577a1eec5ecb15b2.pack: ok

`

	packFile, err := os.Open(path.Join(RepoDir.Name(), "objects/pack/pack-d310969c4ba0ebfe725685fa577a1eec5ecb15b2.pack"))
	if err != nil {
		t.Error(err)
		return
	}
	defer packFile.Close()

	idxFile, err := os.Open(path.Join(RepoDir.Name(), "objects/pack/pack-d310969c4ba0ebfe725685fa577a1eec5ecb15b2.idx"))
	if err != nil {
		t.Error(err)
		return
	}
	defer idxFile.Close()

	type packObjectMock struct {
		Name           SHA
		_type          packObjectType
		Type           string
		Size           int
		SizeInPackfile int
		Offset         int
		Depth          int
		BaseObjectName SHA
	}
	objs := map[string]packObjectMock{
		"fe89ee30bbcdfdf376beae530cc53f967012f31c": packObjectMock{Name: SHA("fe89ee30bbcdfdf376beae530cc53f967012f31c"), _type: 1, Type: "commit", Size: 267, SizeInPackfile: 184, Offset: 12},
		"3ead3116d0378089f5ce61086354aac43e736b01": packObjectMock{Name: SHA("3ead3116d0378089f5ce61086354aac43e736b01"), _type: 1, Type: "commit", Size: 243, SizeInPackfile: 170, Offset: 196},
		"1d833eb5b6c5369c0cb7a4a3e20ded237490145f": packObjectMock{Name: SHA("1d833eb5b6c5369c0cb7a4a3e20ded237490145f"), _type: 1, Type: "commit", Size: 262, SizeInPackfile: 180, Offset: 366},
		"a7f92c920ce85f07a33f948aa4fa2548b270024f": packObjectMock{Name: SHA("a7f92c920ce85f07a33f948aa4fa2548b270024f"), _type: 1, Type: "commit", Size: 250, SizeInPackfile: 172, Offset: 546},
		"97eed02ebe122df8fdd853c1215d8775f3d9f1a1": packObjectMock{Name: SHA("97eed02ebe122df8fdd853c1215d8775f3d9f1a1"), _type: 1, Type: "commit", Size: 190, SizeInPackfile: 132, Offset: 718},
		"d22fc8a57073fdecae2001d00aff921440d3aabd": packObjectMock{Name: SHA("d22fc8a57073fdecae2001d00aff921440d3aabd"), _type: 2, Type: "tree", Size: 121, SizeInPackfile: 115, Offset: 850},
		"df891299372c34b57e41cfc50a0113e2afac3210": packObjectMock{Name: SHA("df891299372c34b57e41cfc50a0113e2afac3210"), _type: 2, Type: "tree", Size: 25, SizeInPackfile: 37, Offset: 965, Depth: 1, BaseObjectName: "d22fc8a57073fdecae2001d00aff921440d3aabd"},
		"af6e4fe91a8f9a0f3c03cbec9e1d2aac47345d67": packObjectMock{Name: SHA("af6e4fe91a8f9a0f3c03cbec9e1d2aac47345d67"), _type: 3, Type: "blob", Size: 18, SizeInPackfile: 23, Offset: 1002},
		"6b32b1ac731898894c403f6b621bdda167ab8d7c": packObjectMock{Name: SHA("6b32b1ac731898894c403f6b621bdda167ab8d7c"), _type: 3, Type: "blob", Size: 1645, SizeInPackfile: 700, Offset: 1025},
		"7147f43ae01c9f04a78d6e80544ed84def06e958": packObjectMock{Name: SHA("7147f43ae01c9f04a78d6e80544ed84def06e958"), _type: 3, Type: "blob", Size: 1824, SizeInPackfile: 697, Offset: 1725},
		"05d3cc770bd3524cc25d47e083d8942ad25033f0": packObjectMock{Name: SHA("05d3cc770bd3524cc25d47e083d8942ad25033f0"), _type: 3, Type: "blob", Size: 16, SizeInPackfile: 28, Offset: 2422, Depth: 1, BaseObjectName: "7147f43ae01c9f04a78d6e80544ed84def06e958"},
		"c3b8133617bbdb72e237b0f163fade7fbf1f0c18": packObjectMock{Name: SHA("c3b8133617bbdb72e237b0f163fade7fbf1f0c18"), _type: 3, Type: "blob", Size: 381, SizeInPackfile: 317, Offset: 2450, Depth: 2, BaseObjectName: "05d3cc770bd3524cc25d47e083d8942ad25033f0"},
		"8264d7bcc297e15c452a7aef3a2e40934762b7e3": packObjectMock{Name: SHA("8264d7bcc297e15c452a7aef3a2e40934762b7e3"), _type: 2, Type: "tree", Size: 25, SizeInPackfile: 38, Offset: 2767, Depth: 1, BaseObjectName: "d22fc8a57073fdecae2001d00aff921440d3aabd"},
		"254671773e8cd91e07e36546c9a2d9c27e8dfeec": packObjectMock{Name: SHA("254671773e8cd91e07e36546c9a2d9c27e8dfeec"), _type: 2, Type: "tree", Size: 121, SizeInPackfile: 115, Offset: 2805},
		"ba74813270ff557c4a5d1be0562a141bbee4d3e6": packObjectMock{Name: SHA("ba74813270ff557c4a5d1be0562a141bbee4d3e6"), _type: 3, Type: "blob", Size: 16, SizeInPackfile: 28, Offset: 2920, Depth: 1, BaseObjectName: "6b32b1ac731898894c403f6b621bdda167ab8d7c"},
		"b45377f6daf59a4cec9e8de64f5df1533a7994cd": packObjectMock{Name: SHA("b45377f6daf59a4cec9e8de64f5df1533a7994cd"), _type: 3, Type: "blob", Size: 10, SizeInPackfile: 21, Offset: 2948, Depth: 1, BaseObjectName: "7147f43ae01c9f04a78d6e80544ed84def06e958"},
		"9de6c72106b169990a83ce7090c7cad84b6b506b": packObjectMock{Name: SHA("9de6c72106b169990a83ce7090c7cad84b6b506b"), _type: 2, Type: "tree", Size: 38, SizeInPackfile: 49, Offset: 2969},
	}

	objects, err := VerifyPack(packFile, idxFile)
	if err != nil {
		t.Error(err)
	}

	if len(objects) != len(objs) {
		t.Errorf("Read incorrect number of objects: %d, want %d", len(objects), len(objs))
	}

	for _, object := range objects {
		expectedObj, ok := objs[string(object.Name)]
		if !ok {
			t.Errorf("encountered incorrect hash %s", object.Name)
		}
		if object.err != nil {
			log.Printf("%+v", object)
			log.Printf("%+v", object.err)
			t.Errorf("error reading object %s: %s", object.Name, object.err)
		}
		if expectedObj._type.String() != object._type.String() && object._type.String() != OBJ_OFS_DELTA.String() {
			t.Errorf("expected type %s and received type %s", expectedObj._type, object._type.String())
		}

		if expectedObj.Name != object.Name {
			t.Errorf("Expected Name %s and received %s", expectedObj.Name, object.Name)
		}

		if expectedObj._type.String() != object.PatchedType().String() {
			t.Errorf("Expected _type.String() %s and received %s", expectedObj._type.String(), object.PatchedType().String())
		}

		if expectedObj.Type != object.Type() {
			t.Errorf("Expected Type() method to return %s and received %s (%s)", expectedObj.Type, object.Type(), object.Name)
		}

		if expectedObj.Size != object.Size {
			t.Errorf("Expected Size %d and received %d", expectedObj.Size, object.Size)
		}

		if expectedObj.Offset != object.Offset {
			t.Errorf("Expected Offset %d and received %d", expectedObj.Offset, object.Offset)
		}
	}

	/*
		if !reflect.DeepEqual(expected, result) {
			t.Errorf("Expected and result don't match:\n%+v\n%+v", expected, result)
		}
	*/
}

func BenchmarkVerifyPack(b *testing.B) {
	packFile, err := os.Open(path.Join(RepoDir.Name(), "objects/pack/pack-d310969c4ba0ebfe725685fa577a1eec5ecb15b2.pack"))
	if err != nil {
		b.Error(err)
		return
	}
	defer packFile.Close()

	idxFile, err := os.Open(path.Join(RepoDir.Name(), "objects/pack/pack-d310969c4ba0ebfe725685fa577a1eec5ecb15b2.idx"))
	if err != nil {
		return
	}

	defer idxFile.Close()
	for i := 0; i < b.N; i++ {
		_, err := VerifyPack(packFile, idxFile)
		if err != nil {
			b.Errorf("error in iteration %d: %s", i, err)
		}
		packFile.Seek(0, io.SeekStart)
		idxFile.Seek(0, io.SeekStart)
	}
}
