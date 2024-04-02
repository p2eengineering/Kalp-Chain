package chaincode

// #cgo pkg-config: python3-embed
// #include <Python.h>
import "C"

import (
	"embed"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)


// custom error 
type CompilerError struct{
	msg string
}


func (cerr *CompilerError) Error() string{
	return  cerr.msg
}


//go:embed lib.linux-x86_64-3.8/*
var embeddedScripts embed.FS

func visit(path string, info os.DirEntry, err error) error {
	if err != nil {
		fmt.Printf("Encountered error: %v\n", err)
		return nil
	}

	if info.IsDir() {
		return nil // Skip directories
	}

	fmt.Println(path)
	return nil
}


// returns true if the code correctly compiles as per lark rules
// [Note] : Picks up embedded scripts from hardcoded lib.linux-x86_64-3.8 -- needs changes as-well
// [Note] : this is a tentaive commit  to get things to work / needs clean up and error handling
func CompilerValidate(filePath string) (res bool){

	C.Py_Initialize()

	tempDir, err := ioutil.TempDir("", "myapp_resources")
	if err != nil {
		fmt.Println("Failed to create temporary directory:", err)
		return
	}

	defer os.RemoveAll(tempDir)
	
	fmt.Println("Temporary directory:", tempDir)
	
	
	absolutePath, err := filepath.Abs(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		// os.Exit(1)
	}
	
	fmt.Println("Absolute path:", absolutePath)

	

	// Extract the embedded script files to the temporary directory
	err1 := fs.WalkDir(embeddedScripts, "lib.linux-x86_64-3.8", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			// Skip directories
			return nil
		}

		// Read the embedded file
		scriptContent, err := embeddedScripts.ReadFile(path)
		if err != nil {
			return err
		}

		// Build the absolute path of the file in the temporary directory
		absPath := filepath.Join(tempDir, path)

		// Create any necessary parent directories
		err = os.MkdirAll(filepath.Dir(absPath), 0755)
		if err != nil {
			return err
		}

		// Write the file to the temporary directory
		err = os.WriteFile(absPath, scriptContent, 0644)
		if err != nil {
			return err
		}

		return nil
	})

	if err1 != nil {
		fmt.Println("Error embedding scripts:", err1)
		return
	}

	// Get the current Python module search path
	sysPath := C.PySys_GetObject(C.CString("path"))
	if sysPath == nil {
		fmt.Println("Error getting sys.path")
		return
	}

	// Convert sys.path to a Go slice of strings
	var pathSlice []string
	numPaths := C.PyList_Size(sysPath)
	for i := C.Py_ssize_t(0); i < numPaths; i++ {
		pathItem := C.PyList_GetItem(sysPath, i)
		pathStr := C.PyUnicode_AsUTF8(pathItem)
		pathSlice = append(pathSlice, C.GoString(pathStr))
	}

	// Change current directory to the temporary directory
	err = os.Chdir(tempDir)
	if err != nil {
		fmt.Println("Failed to change directory:", err)
		return
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return
	}

	concatenated := fmt.Sprintf("%s%s", cwd, "/lib.linux-x86_64-3.8")

	err = os.Chdir(concatenated)
	if err != nil {
		fmt.Println("Failed to change directory:", err)
		return
	}

	cwd1, err1 := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err1)
		return
	}

	fmt.Println("current working directory:", cwd1)

	errw := filepath.WalkDir(cwd1, visit)
	if errw != nil {
		fmt.Printf("WalkDir error: %v\n", errw)
	}

	pathSlice = append(pathSlice, cwd1)

	// Convert the modified path slice back to a Python list
	modifiedPath := C.PyList_New(C.Py_ssize_t(len(pathSlice)))
	for i, path := range pathSlice {
		pathStr := C.CString(path)
		pathItem := C.PyUnicode_FromString(pathStr)
		C.PyList_SetItem(modifiedPath, C.Py_ssize_t(i), pathItem)
	}

	// Set the modified sys.path back to Python
	C.PySys_SetObject(C.CString("path"), modifiedPath)

	// Import and use the Python module
	moduleName := C.CString("entry")
	module := C.PyImport_ImportModule(moduleName)

	// Check if the module was imported successfully
	if module == nil {
		fmt.Println("Error importing Python module")
		return
	}
	// ... Perform further operations with the module ...

	functionName := C.CString("run")
	function := C.PyObject_GetAttrString(module, functionName)

	if function == nil {
		// panic("Error importing function") // not allowed / change
	}

	state := C.PyEval_SaveThread()

	// go func() {
	_gstate := C.PyGILState_Ensure()
	args := C.PyTuple_New(2)
	
	args1 := C.CString(absolutePath)
	arg1 := C.PyUnicode_FromString(args1)

	args2 := C.CString("INFO")
	arg2 := C.PyUnicode_FromString(args2)

	C.PyTuple_SetItem(args, 0, arg1)
	C.PyTuple_SetItem(args, 1, arg2)

	//helloFunc.Call(C.PyTuple_New(0), C.PyDict_New())
	result := C.PyObject_CallObject(function, args)

	// Finalize the Python interpreter
	resultStr := C.PyObject_Str(result)
	// Get the string value from the result object
	strValue := C.PyUnicode_AsUTF8(resultStr)

	// Print the result
	if value := C.GoString(strValue); value == "0"{
		return true
	}

	C.PyGILState_Release(_gstate)

	C.PyEval_RestoreThread(state)
	C.Py_Finalize()

	return res
}