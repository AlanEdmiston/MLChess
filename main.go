package main

type player struct {
	colour   string
	strategy string
}

func main() {
	writeRandomConfigs("./policies/", 2)
	tournament("./policies/")
}
