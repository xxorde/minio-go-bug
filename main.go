package main

import (
	"os/exec"
	"time"

	log "github.com/Sirupsen/logrus"

	minio "github.com/minio/minio-go"
	// This repo contains a hotfix
	// minio "github.com/xxorde/minio-go"
)

func main() {
	// Prepare some command to execute and get output via pipe!
	cmd := exec.Command("uptime")
	stdout, err := cmd.StdoutPipe()

	// Settings for minio
	endpoint := "127.0.0.1:9000"
	accessKeyID := "accessKeyID"
	secretAccessKey := "secretAccessKey"
	ssl := false
	bucket := "test"
	location := "us-east-1"

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, ssl)
	if err != nil {
		log.Fatal(err)
	}

	// Creates bucket with name bucket
	err = minioClient.MakeBucket(bucket, location)
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, err := minioClient.BucketExists(bucket)
		if err == nil && exists {
			log.Infof("We already own %s\n", bucket)
		} else {
			log.Fatal(err)
		}
	}

	// Write output of the executed command to a bucket, run in background
	go func() {
		// This should write the output of the uptime command.
		// But it writes nothing, 0 bytes.
		// PutObjectget uses ReaderSize(stdout) to determine how much to write.
		// ReaderSize(stdout) returns 0 bytes, therefor 0 bytes are written.
		_, err := minioClient.PutObject(bucket, "test", stdout, "")
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Start the process (in the background)
	if err := cmd.Start(); err != nil {
		log.Fatal("cmd failed on startup, ", err)
	}

	// Wait for process to finish
	time.Sleep(time.Second * 1)
}

/* Fix to ReaderSize()
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
*/
