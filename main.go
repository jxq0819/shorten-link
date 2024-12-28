package main

func main() {
	a := App{}
	a.Initialise(getEnv())
	a.Run(":8000")
}
