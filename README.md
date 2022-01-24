# kdu
A concurrent du program  📺

# Current known issue
file size != file size on disk 🚬
Todo: use syscall instead of file size

linux: stat is for directory, statfs for is file
windows: golang.org/x/sys GetFileInformationByHandleEx 