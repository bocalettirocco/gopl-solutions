package main

import "fmt"

type AuthError struct {
	msg string
}

type LabelError struct {
	msg string
}

type RepoError struct {
	msg string
}

type RequestError struct {
	msg string
}

type TitleError struct {
	msg string
}

type UsageError struct {
	msg string
}

func (e AuthError) Error() string {
	return fmt.Sprintf("ERROR - TOKEN: %s\n", e.msg)
}

func (e UsageError) Error() string {
	return fmt.Sprintf("ERROR - USAGE: %s\n", e.msg)
}

func (e RepoError) Error() string {
	return fmt.Sprintf("ERROR - INPUT: %s\n", e.msg)
}

func (e RequestError) Error() string {
	return fmt.Sprintf("ERROR - REQUEST: %s\n", e.msg)
}

func (e LabelError) Error() string {
	return fmt.Sprintf("ERROR - INPUT: %s\n", e.msg)
}

func (e TitleError) Error() string {
	return fmt.Sprintf("ERROR - INPUT: %s\n", e.msg)
}
