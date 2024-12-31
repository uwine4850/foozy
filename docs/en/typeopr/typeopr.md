## typeopr
This package contains tools for operations with different types.

__IsPointer__
```
IsPointer(a any) bool
```
Checks if the passed value is a pointer.

__PtrIsStruct__
```
PtrIsStruct(a any) bool
```
Checks if the pointer points to a structure

__IsEmpty__
```
IsEmpty(value interface{}) bool
```
Checks if the passed value is empty.

__AnyToBytes__
```
AnyToBytes(value interface{}) ([]byte, error)
```
Converts any value to bytes.

## type IPtr interface
This interface is used to implement passing any type as a pointer. This means that 
all values placed in this object must be passed through a pointer.

__New__
```
New(value interface{}) IPtr
```
The function sets a value to an object. The value must be passed through a pointer.

__Ptr__
```
Ptr() interface{}
```
Retrieve a pointer that was passed to the object earlier.<br>
IMPORTANT: it is the pointer to the value that is passed, not the value itself.
___

__IsImplementInterface__
```
IsImplementInterface(objectPtr IPtr, interfaceType interface{}) bool
```
Checks if the object implements the desired interface.<br>
Usage example:
```
object := MyObject{}
IsImplementInterface(typeopr.Ptr{}.New(&object), (*MyInterface)(nil))
```