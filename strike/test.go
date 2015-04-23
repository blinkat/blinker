package strike

import (
	"fmt"
	"time"
)

//w

func Test(text []byte) {
	now := time.Now()
	//ret := strike_css(text)
	//strike_css(text)
	StrikeJs(string(text))
	fmt.Println("use time:", time.Now().Sub(now), "ms")
	//fmt.Println("content:", string(ret))
	//fmt.Println(ret)
}
