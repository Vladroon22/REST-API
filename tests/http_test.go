package tests

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestGetAddress(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8000/", strings.NewReader("{}"))
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Test Denied - %v\n", resp.StatusCode)
	} else {
		fmt.Printf("Test Passed: %v\n", resp.StatusCode)
	}
}

func TestGetAuth(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8000/auth/", strings.NewReader("{}"))
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Test Denied - %v\n", resp.StatusCode)
	} else {
		fmt.Printf("Test Passed: %v\n", resp.StatusCode)
	}
}

func TestGetAuthIn(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:8000/auth/sign-in", strings.NewReader("{}"))
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Test Denied - %v\n", resp.StatusCode)
	} else {
		fmt.Printf("Test Passed: %v\n", resp.StatusCode)
	}
}

func TestGetAuthUp(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:8000/auth/sign-up", strings.NewReader("{}"))
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Test Denied - %v\n", resp.StatusCode)
	} else {
		fmt.Printf("Test Passed: %v\n", resp.StatusCode)
	}
}

func TestPutUsers(t *testing.T) {
	req, err := http.NewRequest(http.MethodPut, "http://127.0.0.1:8000/users/{1}", strings.NewReader("{}"))
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Test Denied - %v\n", resp.StatusCode)
	} else {
		fmt.Printf("Test Passed: %v\n", resp.StatusCode)
	}
}

func TestPatchUserName(t *testing.T) {
	req, err := http.NewRequest(http.MethodPatch, "http://127.0.0.1:8000/users/name/{1}", strings.NewReader("{}"))
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Test Denied - %v\n", resp.StatusCode)
	} else {
		fmt.Printf("Test Passed: %v\n", resp.StatusCode)
	}
}

func TestPatchUserEmail(t *testing.T) {
	req, err := http.NewRequest(http.MethodPatch, "http://127.0.0.1:8000/users/email/{1}", strings.NewReader("{}"))
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Test Denied - %v\n", resp.StatusCode)
	} else {
		fmt.Printf("Test Passed: %v\n", resp.StatusCode)
	}
}

func TestPatchUserPass(t *testing.T) {
	req, err := http.NewRequest(http.MethodPatch, "http://127.0.0.1:8000/users/pass/{1}", strings.NewReader("{}"))
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Test Denied - %v\n", resp.StatusCode)
	} else {
		fmt.Printf("Test Passed: %v\n", resp.StatusCode)
	}
}

func TestDeleteUsers(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, "http://127.0.0.1:8000/users/{1}", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 || resp.StatusCode == 500 { // 500 - nothing to delete
		fmt.Printf("Test Passed: %v\n", resp.StatusCode)
	} else {
		t.Errorf("Test Denied - %v\n", resp.StatusCode)
	}
}
