package printer

/*
// try in command

func main() {
	// Load some text for our viewport
	content, err := os.ReadFile("artichoke.md")
	if err != nil {
		fmt.Println("could not load file:", err)
		os.Exit(1)
	}
	p := NewPaper(context.Background()).SetTitle("artichoke.md").
		SetContent(string(content)).
		SetLoading("loading...").
		EnableProgress(true)
	p.Run()
	go func() {
		time.Sleep(time.Second * 5)
		p.Quit()
	}()
	//select {
	//case <- p.Ctx().Done():
	//	fmt.Println("done")
	//}
	p.Wait()
}
*/
