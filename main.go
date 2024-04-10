package main

import (
	"bytes"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	md "github.com/nao1215/markdown"
)

func main() {
	buf := new(bytes.Buffer)
	md.NewMarkdown(buf).
		H1("This is H1").
		PlainText("This is plain text").
		H2f("This is %s with text format", "H2").
		PlainTextf("Text formatting, such as %s and %s, %s styles.",
			md.Bold("bold"), md.Italic("italic"), md.Code("code")).
		H2("Code Block").
		CodeBlocks(md.SyntaxHighlightHaskell,
			`module Main where
import Data.Function ( (&) )
import Data.List (intercalate)

hello :: String -> String
hello s =
	"Hello, " ++ s ++ "."
main :: IO ()
main =
	map hello [ "artichoke", "alcachofa" ] & intercalate "\n" & putStrLn
			`).
		H2("List").
		BulletList("Bullet Item 1", "Bullet Item 2", "Bullet Item 3").
		OrderedList("Ordered Item 1", "Ordered Item 2", "Ordered Item 3").
		H2("CheckBox").
		CheckBox([]md.CheckBoxSet{
			{Checked: false, Text: md.Code("sample code")},
			{Checked: true, Text: md.Link("Go", "https://golang.org")},
			{Checked: false, Text: md.Strikethrough("strikethrough")},
		}).
		H2("Blockquote").
		Blockquote("If you can dream it, you can do it.").
		H3("Horizontal Rule").
		HorizontalRule().
		H2("Table").
		Table(md.TableSet{
			Header: []string{"Name", "Age", "Country"},
			Rows: [][]string{
				{"David", "23", "USA"},
				{"John", "30", "UK"},
				{"Bob", "25", "Canada"},
			},
		}).
		H2("Image").
		PlainTextf(md.Image("sample_image", "./sample.png")).
		Build()
	pk, err := NewPokemonViewModel(buf.String())
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	if _, err := tea.NewProgram(pk).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
