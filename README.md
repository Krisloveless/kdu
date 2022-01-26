# kdu

*A concurrent du program* ![localhost](https://user-images.githubusercontent.com/15829091/151237736-339bca80-8a4e-4246-a70b-2b5cfc61885c.gif)


# Current known issue
**file size != file size on disk** 🚬


Todos: 

### use syscall instead of file size

- linux: stat is for directory, statfs for is file
- darwin: golang.org/x/sys unix
- windows: golang.org/x/sys GetFileInformationByHandleEx 

### testcases: setup github action
