## codegen
This package is designed to generate new files and packages. To work properly, the 
source file or package must have already been created. That is, the path to a real 
file is passed to the generator, and it, in turn, generates the same file in 
another location.

__Generate__
```
Generate(data map[string]string) error
```
Generates a file at the specified path.<br>
Key — path to the target directory where the new file should be.<br>
Value — path to the file to be generated.