package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

type Collection struct {
	Records []Record `xml:"record"`
}

type Record struct {
	Datafields []Datafield `xml:"datafield"`
}

type Datafield []struct {
	Tag       string     `xml:"tag,attr"`
	Subfields []Subfield `xml:"subfield"`
}

type Subfield struct {
	Code     string `xml:"code,attr"`
	Subfield string `xml:",cdata"`
}

type Search struct {
	Content []Content
}

type Content struct {
	Field string
	Data  []Data
}

type Data struct {
	Subfield string
	Text     []string
}

func main() {

	fileRead, errOpen := os.Open("RNOFA.xml")
	if errOpen != nil {
		panic("Couldn't open the file")
	}

	defer fileRead.Close()

	byteValue, _ := ioutil.ReadAll(fileRead)

	var collection Collection

	xml.Unmarshal(byteValue, &collection)

	p := Search{
		[]Content{
			{"200",
				[]Data{
					{
						"a",
						[]string{"delfim", "Peregrinação"},
					},
				},
			},

			{"325",

				[]Data{
					{
						"a",
						[]string{"Bárbara Guimarães", "Fernando Alvim"},
					},
				},
			},
		},
	}

	var zips, remove_zips = []string{}, []string{}
	var zip = ""

	for a := 0; a < len(collection.Records); a++ {

		for i := 0; i < len(p.Content); i++ {

			for g := 0; g < len(p.Content[i].Data); g++ {

				for h := 0; h < len(p.Content[i].Data[g].Text); h++ {

					for j := 0; j < len(collection.Records[a].Datafields); j++ {

						for k := 0; k < len(collection.Records[a].Datafields[j]); k++ {

							if collection.Records[a].Datafields[j][k].Tag == "339" {

								for x := 0; x < len(collection.Records[a].Datafields[j][k].Subfields); x++ {

									if collection.Records[a].Datafields[j][k].Subfields[x].Code == "e" {

										zip = collection.Records[a].Datafields[j][k].Subfields[x].Subfield + ".zip"
										zips = append(zips, zip)

									}
								}
							}
							if strings.Contains(strings.ToLower(collection.Records[a].Datafields[j][k].Tag), strings.ToLower(p.Content[i].Field)) && strings.Contains(strings.ToLower(collection.Records[a].Datafields[j][k].Subfields[0].Code), strings.ToLower(p.Content[i].Data[g].Subfield)) {

								if strings.Contains(strings.ToLower(collection.Records[a].Datafields[j][k].Subfields[0].Subfield), strings.ToLower(p.Content[i].Data[g].Text[h])) {

									remove_zips = append(remove_zips, zip)

								}
							}
						}
					}
				}
			}
		}
	}

	// keep only the zips to transfer
	z := diff(zips, remove_zips)
	// remove the duplicated zips
	final_zips := removeDuplicate(z)

	for _, zip := range final_zips {
		transfer_file(zip)
	}

	fileRead.Close()

}

// https://stackoverflow.com/questions/53194031/how-to-delete-duplicate-elements-between-slices-on-golang
func diff(a []string, b []string) []string {
	var shortest, longest *[]string
	if len(a) < len(b) {
		shortest = &a
		longest = &b
	} else {
		shortest = &b
		longest = &a
	}
	// Turn the shortest slice into a map
	var m map[string]bool
	m = make(map[string]bool, len(*shortest))
	for _, s := range *shortest {
		m[s] = false
	}
	// Append values from the longest slice that don't exist in the map
	var diff []string
	for _, s := range *longest {
		if _, ok := m[s]; !ok {
			diff = append(diff, s)
			continue
		}
		m[s] = true
	}
	// Append values from map that were not in the longest slice
	for s, ok := range m {
		if ok {
			continue
		}
		diff = append(diff, s)
	}
	// Sort the resulting slice
	sort.Strings(diff)
	return diff
}

// https://stackoverflow.com/questions/66643946/how-to-remove-duplicates-strings-or-int-from-slice-in-go
func removeDuplicate[T string | int](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func transfer_file(zip string) {

	sourceFolder := "C:\\Users\\mbaptista\\OneDrive - Biblioteca Nacional de Portugal\\Documentos\\go\\src\\rnofa-porto\\origem\\"
	destinationFolder := "C:\\Users\\mbaptista\\OneDrive - Biblioteca Nacional de Portugal\\Documentos\\go\\src\\rnofa-porto\\destino\\"

	fmt.Println("Transfering file: ", zip)

	srcFile, err := os.Open(sourceFolder + zip)
	check(err)
	defer srcFile.Close()

	destFile, err := os.Create(destinationFolder + zip) // creates if file doesn't exist

	check(err)
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile) // check first var for number of bytes copied
	check(err)

	err = destFile.Sync()
	check(err)
}

func check(err error) {
	if err != nil {
		fmt.Println("Error : %s", err.Error())
		os.Exit(1)
	}
}
