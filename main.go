package main

import (
	"log"
	"os"
	"os/exec"
	"io/ioutil"

	"github.com/fsnotify/fsnotify"
)

func main() {


	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					//执行命令
					execCommand()
					log.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func execCommand(){
	cmd := exec.Command("/bin/sh", "-c", os.Args[2])
	//创建获取命令输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("Error:can not obtain stdout pipe for command:%s\n", err)
	return
	}
	//执行命令
	if err := cmd.Start(); err != nil {
		log.Println("Error:The command is err,", err)
	return
	}
	//读取所有输出
	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Println("ReadAll Stdout:", err.Error())
	return
	}
	if err := cmd.Wait(); err != nil {
		log.Println("wait:", err.Error())
	return
	}
	log.Printf("stdout:\n %s", bytes)
}
