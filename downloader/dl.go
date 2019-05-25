package main

import (
	"strings"
    "strconv"
	"fmt"
	"log"
	"os"

	"github.com/reznov53/law-cots2/mq"
	"github.com/reznov53/law-cots2/download"
)

func deleteEmpty(s []string) []string {
    var r []string
    for _, str := range s {
        if str != "" {
            r = append(r, str)
        }
    }
    return r
}

func joint(i string, j string) string {
	var str strings.Builder
	str.WriteString(i)
	str.WriteString(j)
	return str.String()
}

func dl(arr []string, ch *mq.Channel, id string) {
	path := joint(joint("files/", id), "/")
	os.MkdirAll(path, 0755)
	for i, v := range arr {
		go func(i int, v string, ch *mq.Channel, routeKey string) {
			splits := strings.Split(v, "/")
			err := download.File(joint(path, splits[len(splits) - 1]), v, ch, routeKey)
			if err != nil {
				log.Println(err)
				ch.PostMessage("Failed to download", routeKey)
			}
		}(i, v, ch, joint("dlstatus", fmt.Sprint(strconv.Itoa(i))))
	}
}