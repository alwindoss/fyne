// +build ignore

package main

import "fmt"
import "io/ioutil"
import "os"
import "path"
import "runtime"
import "strings"

import "github.com/fyne-io/fyne"

func bundleFile(name string, filepath string, f *os.File) {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println("Unable to load file " + filepath)
		return
	}
	res := fyne.NewResource(path.Base(filepath), bytes)

	_, err = f.WriteString("var " + name + " = " + res.ToGo() + "\n")
	if err != nil {
		fmt.Println("Unable to write to bundled file")
	}
}

func bundleFont(name string, f *os.File) {
	_, dirname, _, _ := runtime.Caller(0)
	path := path.Join(path.Dir(dirname), "font", "NotoSans-"+name+".ttf")

	bundleFile(strings.ToLower(name), path, f)
}

func bundleIcon(name, theme string, f *os.File) {
	_, dirname, _, _ := runtime.Caller(0)
	path := path.Join(path.Dir(dirname), "icons", name+"-24px-"+theme+".svg")

	formatted := strings.ToLower(name) + strings.Title(strings.ToLower(theme))
	bundleFile(formatted, path, f)
}

func openFile(filename string) *os.File {
	os.Remove(filename)
	_, dirname, _, _ := runtime.Caller(0)
	f, err := os.Create(path.Join(path.Dir(dirname), filename))
	if err != nil {
		fmt.Println("Unable to open file " + filename)
		return nil
	}

	_, err = f.WriteString("package theme\n\nimport \"github.com/fyne-io/fyne\"\n\n")
	if err != nil {
		fmt.Println("Unable to write file " + filename)
		return nil
	}

	return f
}

func main() {
	f := openFile("bundled-fonts.go")
	if f == nil {
		return
	}
	bundleFont("Regular", f)
	bundleFont("Bold", f)
	bundleFont("Italic", f)
	bundleFont("BoldItalic", f)
	f.Close()

	f = openFile("bundled-icons.go")
	if f == nil {
		return
	}

	themes := []string{"dark", "light"}
	for _, theme := range themes {
		bundleIcon("cancel", theme, f)
		bundleIcon("check", theme, f)
	}
	f.Close()
}