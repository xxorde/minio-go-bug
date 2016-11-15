This little piece of code demonstrates a bug in minio-go.

The bug appears in ReaderSize() but affects the function PutObject() and maybe others.

A hotfix can be found here: github.com/xxorde/minio-go

Usage
=====
1. Change keys in main.go
2. go run main.go; cat minio-path/test/test

If you get the output of "uptime" the software worked.
If you do not get this output something went wrong.

Hotfix
======
```
diff --git a/api-put-object.go b/api-put-object.go
index f7dd2da..7267088 100644
--- a/api-put-object.go
+++ b/api-put-object.go
@@ -111,6 +111,13 @@ func getReaderSize(reader io.Reader) (size int64, err error) {
                                return
                        }
                        size = st.Size()
+                       // FileInfo.Size() returns:
+                       //      length in bytes for regular files; system-dependent for others
+                       // For other types like pipes it can return 0 instead.
+                       // This is an ugly fix to make pipes work on linux systems.
+                       if size == 0 {
+                               size = -1
+                       }
                case *Object:
                        var st ObjectInfo
                        st, err = v.Stat()
```
