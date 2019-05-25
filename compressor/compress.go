package main

import (
	"path/filepath"
	"os"
	"fmt"
	"archive/tar"
	"strings"
	"io"
	"log"
	"bytes"
	"strconv"
	"io/ioutil"

	"github.com/reznov53/law-cots2/mq"
)

func tarit(source, target string) error {
	filename := filepath.Base(source)
	target = filepath.Join(target, fmt.Sprintf("%s.tar", filename))
	tarfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer tarfile.Close()

	tarball := tar.NewWriter(tarfile)
	defer tarball.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	return filepath.Walk(source, 
	func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}

		if err := tarball.WriteHeader(header); err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(tarball, file)
		return err
	})
}

func initCh(url string, vhost string, exchangeName string, exchangeType string, queueName string) (*mq.Channel, error) {
	ch, err := mq.InitMQ(url, vhost)
	if err != nil {
		return ch, err
	}

	err = ch.ExcDeclare(exchangeName, exchangeType)
	if err != nil {
		return ch, err
	}

	return ch, nil
}

func main() {
	// url := "amqp://" + os.Getenv("UNAME") + ":" + os.Getenv("PW") + "@" + os.Getenv("URL") + ":" + os.Getenv("PORT") + "/"
	url := "amqp://1406568753:167664@152.118.148.103:5672/"
	// vhost := os.Getenv("VHOST")
	vhost := "1406568753"
	// exchangeName := os.Getenv("EXCNAME")
	exchangeName := "1406568753-compress"
	exchangeType := "fanout"
	exchangeName1 := "1406568753-frontdl"
	exchangeType1 := "fanout"
	
	ch, err := initCh(url, vhost, exchangeName, exchangeType, "compresspass")
	if err != nil {
		panic(err)
	}

	ch1, err := initCh(url, vhost, exchangeName1, exchangeType1, "dlpass")
	if err != nil {
		panic(err)
	}

	err = ch.QueueDeclare("compresspass")
	if err != nil {
		panic(err)
	}

	err = ch1.QueueDeclare("dlpass")
	if err != nil {
		panic(err)
	}

	msgs, err := ch.Ch.Consume(
		"compresspass", // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			id := string(d.Body)
			source := "/files/" + id
			target := "/files"
			err := tarit(source, target)
			if err != nil {
				continue
			}
			file, err := os.Open(source + ".tar")
			if err != nil {
				continue
			}
			buf := bytes.NewBuffer(nil)
			if _, err := io.Copy(buf, file); err != nil {
				continue
			}

			chunks := Split(buf.Bytes())

			var cfiles [10][]byte
			var err1 error

			compression := true

			for i, v := range chunks {
				// log.Println(i)
				cfiles[i], err1 = Compress(v)
				if err1 != nil {
					ch1.PostMessage("Failed to compress", "dlpass")
					compression = false
					break
				}
				
				percentage := (i+1) * 10
				// log.Println(string(percentage))
				ch1.PostMessage(strconv.Itoa(percentage) + "% Compressed", "dlpass")
				// time.Sleep(1 * time.Second)
			}

			if compression {
				cfile := Combine(cfiles)
				filename := source + ".tar.gz"
		
				err = ioutil.WriteFile(filename, cfile, 0755)
	
				currURL := "http://152.118.148.103:21705/download/" + id
	
				ch1.PostMessage(currURL, "dlpass")
			}
		}
	}()

	log.Printf("[*] Waiting for messages. To exit press CTRL+C")
	<-forever
}