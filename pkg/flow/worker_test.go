package flow

import (
	"fmt"
	"regexp"
	"testing"
)

func TestA(t *testing.T) {
	re := regexp.MustCompile(`(cherry picked from commit ([\w]+))`)
	params := re.FindAllString("()123(cherry picked from commit 8fe8f231cf539e3346a4fd31d9c275bf168f6cc8)123(((23401230123 123", -1)

	for _, param := range params {
		fmt.Println(param)
	}
}
