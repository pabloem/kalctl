//go:build integration
// +build integration

package test

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/pabloem/kalctl/auth"
	"github.com/pabloem/kalctl/commands"
	"github.com/pabloem/kalctl/output"
)

type outputTestProcessor struct {
	outputs []string
}

func (d *outputTestProcessor) Title(title string) string {
	d.outputs = append(d.outputs, title)
	return ""
}

func (d *outputTestProcessor) Description(desc string) string {
	d.outputs = append(d.outputs, desc)
	return ""
}

func (d *outputTestProcessor) Attribute(name string) string {
	d.outputs = append(d.outputs, name)
	return ""
}

func (d *outputTestProcessor) AttributeDescription(desc string) string {
	d.outputs = append(d.outputs, desc)
	return ""
}

func (d *outputTestProcessor) CommandResult(result string) string {
	d.outputs = append(d.outputs, result)
	return ""
}

type CommandItTest struct {
	command  string
	expected map[string]string
}

var COMMANDS_TEST_DATA = []CommandItTest{
	{
		command: "kalctl events list",
		expected: map[string]string{
			"length":  "100",
			"listkey": "events",
		},
	},
	{
		command: "kalctl markets list",
		expected: map[string]string{
			"length":  "100",
			"listkey": "markets",
		},
	},
	{
		command: "kalctl trades list",
		expected: map[string]string{
			"length":  "100",
			"listkey": "trades",
		},
	},
	{
		command: "kalctl exchange get-schedule",
		expected: map[string]string{
			"string_length_min": "100",
			"keys":              "schedule",
		},
	},
	{
		command: "kalctl exchange get-announcements",
		expected: map[string]string{
			"string_length_min": "10",
			"keys":              "exchange_active,trading_active",
		},
	},
}

var outputProcessor = &outputTestProcessor{}

func TestAllCommandsWithoutAuth(t *testing.T) {
	output.DefaultOutputFormatter = outputProcessor
	temp_dir := t.TempDir()
	// We switch the home directory to make sure we don't have any auth token
	t.Setenv("HOME", temp_dir)
	for _, test := range COMMANDS_TEST_DATA {
		fmt.Println("Running command: ", test.command)
		t.Run(test.command, func(t *testing.T) {
			err := commands.RunCommand(strings.Split(test.command, " "))
			if err == nil {
				t.Error("Expected an authentication error, but got nil")
			}
			if !strings.Contains(err.Error(), "unable to get auth token. run 'kalctl auth login' to authenticate") {
				t.Error("Expected an authentication error, but got ", err)
			}
		})
	}
}

func TestAllCommandsWithAuth(t *testing.T) {
	output.DefaultOutputFormatter = outputProcessor
	temp_dir := t.TempDir()
	oldHome := os.Getenv("HOME")
	// We switch the home directory to make sure we don't have any auth token at the
	// beginning, and we get one by running the login command
	t.Setenv("HOME", temp_dir)
	token, err := os.ReadFile(filepath.Join(oldHome, auth.DEFAULT_TOKEN_FILE))
	if err != nil {
		log.Fatal("No authentication. Run 'kalctl auth login' and re-run the test", err)
	}
	os.Mkdir(filepath.Join(temp_dir, ".kalctl"), 0700)
	err = os.WriteFile(filepath.Join(temp_dir, auth.DEFAULT_TOKEN_FILE), token, 0600)
	if err != nil {
		log.Fatal("No authentication. Run 'kalctl auth login' and re-run the test", err)
	}
	for _, test := range COMMANDS_TEST_DATA {
		fmt.Println("Running command: ", test.command)
		t.Run(test.command, func(t *testing.T) {
			err = commands.RunCommand(strings.Split(test.command, " "))
			if err != nil {
				t.Error("Expected no error, but got ", err)
			}

			lastOutput := outputProcessor.outputs[len(outputProcessor.outputs)-1]

			expectedLengthStr, checkLength := test.expected["length"]
			listKey, _ := test.expected["listkey"]
			if checkLength {
				checkError(checkJsonLength(lastOutput, expectedLengthStr, listKey), t)
			}
			expectedStrLengthMin, checkStrLengthMin := test.expected["string_length_min"]
			if checkStrLengthMin {
				checkError(checkStringMinLength(lastOutput, expectedStrLengthMin), t)
			}

			keys, checkKeys := test.expected["keys"]
			if checkKeys {
				checkError(checkContainsKeys(lastOutput, keys), t)
			}
		})
	}
}

func checkContainsKeys(jsonStr string, keysCsv string) error {
	var objmap map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &objmap); err != nil {
		log.Fatal(err)
	}
	keys := strings.Split(keysCsv, ",")
	for _, key := range keys {
		if _, ok := objmap[key]; !ok {
			return fmt.Errorf("Expected key %s not found in %s", key, jsonStr)
		}
	}
	return nil
}

func checkStringMinLength(str string, expectedLengthStr string) error {
	expectedLength, err := strconv.Atoi(expectedLengthStr)
	if err != nil {
		return fmt.Errorf("Expected length is not a number: %s", expectedLengthStr)
	}
	if len(str) < expectedLength {
		return fmt.Errorf("Expected length at least %d, but got %d", expectedLength, len(str))
	}
	return nil
}

func checkJsonLength(jsonStr string, expectedLengthStr string, listKey string) error {
	var objmap map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &objmap); err != nil {
		log.Fatal(err)
	}
	expectedLength, err := strconv.Atoi(expectedLengthStr)
	if err != nil {
		return fmt.Errorf("Expected length is not a number: %s", expectedLengthStr)
	}
	elms, ok := objmap[listKey]
	if !ok {
		return fmt.Errorf("Expected key %s not found in %s", listKey, jsonStr)
	}
	elmList, ok := elms.([]interface{})
	if ok && len(elmList) != expectedLength {
		return fmt.Errorf("Expected length %d, but got %d", expectedLength, len(objmap))
	}
	return nil
}

func checkError(err error, t *testing.T) {
	if err != nil {
		t.Error(err)
	}
}
