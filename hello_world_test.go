package go_initial_project_v2

import "testing"

func TestHelloWorld(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The code did a panic")
		}
	}()
	HelloWorld()
}
