package basal

import (
	"reflect"
	"runtime"
	"unsafe"
)

type functab struct {
	entry   uintptr
	funcoff uintptr
}

type bitvector struct {
	n        int32 // # of bits
	bytedata *uint8
}

type moduledata struct {
	pclntable    []byte
	ftab         []functab
	filetab      []uint32
	findfunctab  uintptr
	minpc, maxpc uintptr

	text, etext           uintptr
	noptrdata, enoptrdata uintptr
	data, edata           uintptr
	bss, ebss             uintptr
	noptrbss, enoptrbss   uintptr
	end, gcdata, gcbss    uintptr

	typelinks []int32 // offsets from types

	modulename   string
	modulehashes []interface{}

	gcdatamask, gcbssmask bitvector

	next *moduledata
}

//go:linkname firstmoduledata runtime.firstmoduledata
var firstmoduledata moduledata

func FindFuncWithName(name string) (uintptr, error) {
	for moduleData := &firstmoduledata; moduleData != nil; moduleData = moduleData.next {
		for _, ftab := range moduleData.ftab {
			funcPtr := *(*uintptr)(unsafe.Pointer(&moduleData.pclntable[ftab.funcoff]))
			myFunc := runtime.FuncForPC(funcPtr)

			if myFunc == nil {
				continue
			}
			//fmt.Println(myFunc.Name())
			if myFunc.Name() == name {
				return myFunc.Entry(), nil
			}
		}
	}

	return 0, NewError("invalid function " + name)
}

func GetFunc(outFuncPtr interface{}, name string) (err error) {
	defer exception(func(e error) {
		err = NewError("exception: %v", e)
	})
	var codePtr uintptr
	codePtr, err = FindFuncWithName(name)
	if err == nil {
		err = CreateFuncForCodePtr(outFuncPtr, codePtr)
	}
	return
}

type function struct {
	codePtr uintptr
}

func CreateFuncForCodePtr(outFuncPtr interface{}, codePtr uintptr) (err error) {
	defer exception(func(e error) {
		err = e
	})
	outFuncVal := reflect.ValueOf(outFuncPtr).Elem()
	newFuncVal := reflect.MakeFunc(outFuncVal.Type(), nil)
	funcValuePtr := reflect.ValueOf(newFuncVal).FieldByName("ptr").Pointer()
	funcPtr := (*function)(unsafe.Pointer(funcValuePtr))
	funcPtr.codePtr = codePtr
	outFuncVal.Set(newFuncVal)
	return
}

type _func struct {
	entry   uintptr // start pc
	nameoff int32   // function name

	args   int32  // in/out args size
	funcID uint32 // set for certain special runtime functions

	pcsp      int32
	pcfile    int32
	pcln      int32
	npcdata   int32
	nfuncdata int32
}
