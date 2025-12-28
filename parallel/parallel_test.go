package parallel

import (
	"log"
	"testing"
	"time"
)

func TestForeach(t *testing.T) {
	err := Foreach([]string{"a", "b", "c", "d", "e", "g", "h", "a", "b", "c", "d", "e", "g", "h", "a", "b", "c", "d", "e", "g", "h", "a", "b", "c", "d", "e", "g", "h", "a", "b", "c", "d", "e", "g", "h", "a", "b", "c", "d", "e", "g", "h", "a", "b", "c", "d", "e", "g", "h", "a", "b", "c", "d", "e", "g", "h", "a", "b", "c", "d", "e", "g", "h", "a", "b", "c", "d", "e", "g", "h", "a", "b", "c", "d", "e", "g", "h", "a"}, func(s string) error {
		time.Sleep(1 * time.Second)
		log.Println(s)
		return nil
	})
	if err != nil {
		log.Println(err.Error())
	}
}

func TestMap(t *testing.T) {
	list, err := Map([]string{"1", "e", "r"}, func(t string) (string, error) {
		return t + "aaaaa", nil
	})
	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println(list)
	}

}
