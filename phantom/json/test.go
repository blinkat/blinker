package json

import "fmt"

func Test(path string) {
	json, err := ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}
	json = RemoveComment(json).(Json)
	fmt.Println(json.String())
}
