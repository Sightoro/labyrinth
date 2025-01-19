package tests

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestProgramOutput(t *testing.T) {
	files, err := os.ReadDir("./examples")
	if err != nil {
		t.Fatalf("Ошибка при чтении директории: %v", err)
	}

	var fileCount int
	for _, file := range files {
		if !file.IsDir() {
			fileCount++
		}
	}

	for _, file := range files {
		name_test := file.Name()
		if len(name_test) >= 2 {
			if name_test[len(name_test)-2:] == ".a" {
				continue
			}
		}
		name_ans := file.Name() + ".a"

		in, err := os.Open("./examples/" + name_test)
		if err != nil {
			t.Errorf("Ошибка при открытии входного файла: %v", err)
			continue
		}
		defer in.Close()

		exp_out, err := os.ReadFile("./examples/" + name_ans)
		if err != nil {
			t.Errorf("Ошибка при чтении ожидаемого ответа: %v", err)
			continue
		}

		cmd := exec.Command("go", "run", "../main.go")
		cmd.Stdin = in

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err = cmd.Run()
		if err != nil {
			t.Errorf("Ошибка выполнения программы для теста %v: CE\nЛоги ошибок: %v", name_test, stderr.String())
			continue
		}

		programOutput := strings.TrimSpace(stdout.String())
		expectedOutput := strings.TrimSpace(string(exp_out))

		if programOutput != expectedOutput {
			t.Errorf("Неправильный ответ для теста %v: WA\nПолучено:\n%v\nОжидалось:\n%v\nЛоги ошибок:\n%v", name_test, programOutput, expectedOutput, stderr.String())
		} else {
			t.Logf("Тест %v пройден успешно.", name_test)
		}
	}

	if t.Failed() {
		fmt.Println("Некоторые тесты завершились с ошибками. Проверьте вывод выше.")
	} else {
		fmt.Println("Все тесты пройдены успешно.")
	}
}
