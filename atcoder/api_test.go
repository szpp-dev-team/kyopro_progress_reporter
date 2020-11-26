package atcoder

import (
	"fmt"
	"log"
	"testing"
)

func TestGetSubmissionResult(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		sr, err := GetSubmissionResult("earlgray283")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(len(*sr))
	})

}

func TestGetUniqueACs(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		sr, err := GetUniqueAC("earlgray283")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(len(*sr))
	})

}
