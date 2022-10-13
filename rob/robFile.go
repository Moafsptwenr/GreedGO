package rob

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var paths = make(chan map[string]string, 200)
var results = make(chan map[string]string, 200)

var regexes = []*regexp.Regexp{
	regexp.MustCompile(`(?i)user`),
	regexp.MustCompile(`(?i)password`),
	regexp.MustCompile(`(?i)kdb`),
	regexp.MustCompile(`(?i)login`),
}

func walkFn(path string, f os.FileInfo, err error) error {
	var v1 = make(map[string]string)
	for _, r := range regexes {
		if r.MatchString(path) {
			fi, err := os.Stat(path)
			if err != nil {
				log.Fatal(err)
			}
			timeLayoutStr := "2006-01-02 15:04"
			v1["path"] = path
			v1["size"] = strconv.Itoa(int(fi.Size()))
			v1["time"] = fi.ModTime().Format(timeLayoutStr)
			paths <- v1
		}
	}

	return nil
}

func filter(times string, path, result chan map[string]string) {
	if times == "all" {
		for val := range path {
			result <- val
		}
		close(results)
	} else if strings.Contains(times, "~") {
		timeslice := strings.Split(string(times), "~")
		time1, err := time.Parse("2006-01-02 15:04:05", timeslice[0])
		if err != nil {
			log.Fatal(err)
		}
		time2, err := time.Parse("2006-01-02 15:04:05", timeslice[1])
		if err != nil {
			log.Fatal(err)
		}
		for val := range path {
			time3, err := time.Parse("2006-01-02 15:04", val["time"])
			if err != nil {
				log.Fatal(err)
			}
			if time3.Before(time2) && time3.After(time1) {
				result <- val
			}
		}
		close(results)
	} else if !strings.Contains(times, "~") && strings.Contains(times, "-") {
		time1, err := time.Parse("2006-01-02 15:04:05", times)
		if err != nil {
			log.Fatal()
		}
		// timeLayoutStr := "2006-01-02 15:04"
		time2 := time.Now()
		for val := range path {
			time3, err := time.Parse("2006-01-02 15:04", val["time"])
			if err != nil {
				log.Fatal(err)
			}
			if time3.After(time1) && time3.Before(time2) {
				result <- val
			}
		}
		close(results)
	} else {
		fmt.Println("[*]还没有此功能")
		close(results)
	}

}

func RobFileWindows() {
	var rel1 []map[string]string
	read := bufio.NewReader(os.Stdin)
	read1 := bufio.NewReader(os.Stdin)
	fmt.Println("[*]格式为:全部输入all,范围例如:2015-05-06 10:20:20~2020-07-05 20:20:10,某个时间到现在直接按格式输入一个时间就可以,如:2020-01-27 20:10:10")
	fmt.Println("[*]输入筛选的时间范围:")
	times, _, err := read1.ReadLine()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(times)

	fmt.Println("[*]传入一个路径:")
	path, _, err := read.ReadLine()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := filepath.Walk(string(path), walkFn); err != nil {
			log.Fatal(err)
		}
		close(paths)
	}()

	go filter(string(times), paths, results)

	for {
		val, ok := <-results
		if ok {
			rel1 = append(rel1, val)
		} else {
			fmt.Println("[*]complete")
			break
		}
	}

	for _, val1 := range rel1 {
		fmt.Println("****************************************")
		fmt.Println("[*]path :", val1["path"])
		fmt.Println("[*]size :", val1["size"])
		fmt.Println("[*]time :", val1["time"])
	}
	fmt.Println("****************************************")
}
