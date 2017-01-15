/*
Copyright 2017 Jean-Christophe Sirot.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"time"
	"os"
	"strconv"
)

func main() {
	count := GetCount();
	sleep := GetSleep();
	
	for i := 1; i <= count; i++ {
		fmt.Printf("%d...\n", i)
		time.Sleep(time.Duration(sleep) * time.Second)
	}
	fmt.Printf("DONE\n")
}

func GetCount() int {
	count, err := strconv.Atoi(os.Getenv("JOB_COUNT"))
	if (err != nil) {
		return 10
	}
	return count
}

func GetSleep() int {
	sleep, err := strconv.Atoi(os.Getenv("JOB_WAIT"))
	if (err != nil) {
		return 10
	}
	return sleep
}
