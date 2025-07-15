## System Programming/Performance/Go: Static Allocation

- *tiny passage*
- *part of the though: What features of Go makes it fit for extreme system engineering and what features make it not a good fit?*

A common technique to make programs faster in system engineering
is to use static memory allocation than dynamic memory allocation.
Others common techniques: 
- zero/less dynamic allocation (static, less string formating)
- zero/less GC (reusing)
- zero/less Copying (resuing/pointer)
- zero/less any/reflection (accurate typing)


### Use Global Arrays
```
// page of a database in memory can easiliy be represented
// as a huge array

var page [4096]byte  // 4KB page buffer, correct choice
var globalSlice = make([]byte, 1024*1024) // bad choice
```
- Use Global Array: In this case, the array is static, laid out at compile time, and persists for the program’s lifetime. 
- Don't Use Slice: The ***slice descriptor*** is a global variable, allocated statically in the data segment. But the underlying array (1 MB) is allocated dynamically at runtime on the heap by the call to make.
So we cannot use slice to manage static allocation in which only the descriptor part is static.

### Check It
In the following snippet, the global data addr should be consistent.
But on my MacOS the virtual addresses are randomized (seems to be fine on linux without PIE/ASLR in place), but still, we can use "go tool nm|grep global" to find the global variables.
See [extreme/staticalloc](https://github.com/maki3cat/gogymnastics/tree/main/extreme/staticalloc) for the code and test script.  
To try it yourself, run `bash test.sh` in that directory.


Symbol Types found by `go tool nm` will show that our page is B or D depending on it is initialized or not.

| Symbol Type | Memory Region/Executbale Layout | Allocation Type    | Description                                                |
| ----------- | ------------------------------- | ------------------ | ---------------------------------------------------------- |
| D           | Data Segment                    | Static Allocation  | Initialized global variables                               |
| B           | BSS Segment                     | Static Allocation  | Uninitialized global variables (`var x int`)               |
| R           | Read-Only Data                  | Static Allocation  | Constants, string literals, etc.                           |
| T, t        | Code/Text Segment               | Static Allocation  | Functions and executable code                              |
| (none)      | Stack (runtime only)            | Stack Allocation   | Function-local variables that are not captured             |
| (none)      | Heap (runtime only)             | Dynamic Allocation | Created via `make`, `new`, or variables that escape        |
| U           | External Symbol                 | —                  | Symbol resolved at link time; not allocated in this binary |