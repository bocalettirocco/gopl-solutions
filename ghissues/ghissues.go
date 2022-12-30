package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "%#v\n", &UsageError{"ghissues command [ARGS...]"})
		os.Exit(1)
	}
	token, err := getToken()

	if err != nil {
		fmt.Fprintf(os.Stderr, "%#v\n", err)
		os.Exit(1)
	}

	cmd := os.Args[1]

	switch cmd {
	case "create":
		err := createIssue(os.Args[2:], token)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%#v\n", err)
		}
	case "list":
		err := listIssues(os.Args[2:], token)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%#v\n", err)
		}
	case "read":
		err := readIssue(os.Args[2:], token)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%#v\n", err)
		}
	case "update":
		err := updateIssue(os.Args[2:], token)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%#v\n", err)
		}
	case "open":
		err := openIssue(os.Args[2:], token)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%#v\n", err)
		}
	case "close":
		err := closeIssue(os.Args[2:], token)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%#v\n", err)
		}
	}
}

func createIssue(args []string, token string) error {

	if len(args) != 4 {
		return UsageError{"ghissues create [ORG] [REPO] [TITLE] [LABEL]"}
	}

	org := args[0]
	repo := args[1]
	title := args[2]
	label := args[3]

	body, err := getInput(nil)

	if err != nil {
		return err
	}

	labelPtr := &Label{label}
	issue := Issue{
		Body:   body,
		Labels: []*Label{labelPtr},
		Title:  title,
	}

	data, err := json.Marshal(issue)

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues", org, repo)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	if err != nil {
		return err
	}
	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		return RequestError{fmt.Sprintf("Failed to create issue with error %d", res.StatusCode)}
	}
	return nil
}

func listIssues(args []string, token string) error {
	if len(args) != 2 {
		return UsageError{"ghissues list [ORG] [REPO]"}
	}
	org := args[0]
	repo := args[1]
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues", org, repo)

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	if err != nil {
		return err
	}

	client := http.DefaultClient
	res, err := client.Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return RequestError{fmt.Sprintf("Failed to list issues with error %d", res.StatusCode)}
	}
	issues := make([]Issue, 0)

	json.NewDecoder(res.Body).Decode(&issues)
	printIssues(issues)

	return nil
}

func getIssue(url string, token string) (*Issue, error) {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	if err != nil {
		return nil, err
	}

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, RequestError{fmt.Sprintf("Failed to get issue with error %d", res.StatusCode)}
	}
	issue := &Issue{}
	json.NewDecoder(res.Body).Decode(issue)
	return issue, nil
}

func patchIssue(url string, token string, issue *Issue) error {

	data, err := json.Marshal(issue)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(data))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	if err != nil {
		return err
	}

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(res.Body)
		log.Default().Printf("%s\n", data)
		return RequestError{fmt.Sprintf("Failed to update issue %d with error %d", issue.Number, res.StatusCode)}
	}

	return nil
}

func readIssue(args []string, token string) error {
	if len(args) != 3 {
		return UsageError{"ghissues read [ORG] [REPO] [NUMBER]"}
	}
	org := args[0]
	repo := args[1]
	num := args[2]
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/%s", org, repo, num)

	issuePtr, err := getIssue(url, token)
	if err != nil {
		return err
	}

	printIssue(*issuePtr)

	return nil
}

func updateIssue(args []string, token string) error {

	if len(args) != 3 {
		return UsageError{"ghissues update [ORG] [REPO] [NUMBER]"}
	}

	org := args[0]
	repo := args[1]
	num := args[2]
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/%s", org, repo, num)

	issuePtr, err := getIssue(url, token)
	if err != nil {
		return err
	}

	input, err := getInput(issuePtr)

	if err != nil {
		return err
	}

	updatedIssue := Issue{}
	err = json.Unmarshal([]byte(input), &updatedIssue)
	if err != nil {
		return err
	}

	err = patchIssue(url, token, &updatedIssue)
	if err != nil {
		return err
	}

	return nil
}

func closeIssue(args []string, token string) error {
	if len(args) != 3 {
		return UsageError{"ghissues close [ORG] [REPO] [NUMBER]"}
	}

	org := args[0]
	repo := args[1]
	num := args[2]
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/%s", org, repo, num)

	issue, err := getIssue(url, token)

	if err != nil {
		return err
	}

	issue.State = "closed"

	err = patchIssue(url, token, issue)
	if err != nil {
		return err
	}

	return nil
}

func openIssue(args []string, token string) error {
	if len(args) != 3 {
		return UsageError{"ghissues open [ORG] [REPO] [NUMBER]"}
	}

	org := args[0]
	repo := args[1]
	num := args[2]
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/%s", org, repo, num)

	issue, err := getIssue(url, token)

	if err != nil {
		return err
	}

	issue.State = "open"

	err = patchIssue(url, token, issue)
	if err != nil {
		return err
	}

	return nil
}

func getInput(issuePtr *Issue) (string, error) {
	filename := fmt.Sprintf("/tmp/ghissues-%d", time.Now().Unix())

	if issuePtr != nil {
		content, err := json.MarshalIndent(*issuePtr, "", "  ")
		if err != nil {
			return "", err
		}
		err = ioutil.WriteFile(filename, content, 0644)
		if err != nil {
			return "", err
		}
	}

	vim := exec.Command("vim", filename)
	vim.Stdin, vim.Stdout, vim.Stderr = os.Stdin, os.Stdout, os.Stderr

	err := vim.Run()
	if err != nil {
		return "", err
	}

	var body []byte
	body, err = ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	os.Remove(filename)

	return string(body), nil
}

func printIssues(issues []Issue) {
	fmt.Fprintln(os.Stdout, "Issue Number\tUsername\tIssue Title")
	for _, issue := range issues {
		fmt.Fprintf(os.Stdout, "%d\t%s\t%s\n", issue.Number, issue.User.Login, issue.Title)
	}
}

func printIssue(issue Issue) {
	fmt.Fprintf(os.Stdout, "Issue Number: %d\n", issue.Number)
	fmt.Fprintf(os.Stdout, "Issue Title: %s\n", issue.Title)
	fmt.Fprintf(os.Stdout, "Issue Description: %s\n", issue.Body)
}

func getToken() (string, error) {
	token := os.Getenv("github_token")
	if token == "" {
		return "", AuthError{"valid github_token not found in environment"}
	}
	return token, nil
}
