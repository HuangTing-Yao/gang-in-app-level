package main

// map[jobName]jobReadyNumber
var jobMember = make(map[string]int)

func main() {
	webapp := NewWebApp()
	webapp.StartWebApp()
}
