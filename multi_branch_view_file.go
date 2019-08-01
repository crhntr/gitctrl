package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/julienschmidt/httprouter"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func multibranchviewfile(repoPath *string, server bool) {
	fmt.Println(*repoPath)
	repo, err := git.PlainOpen(*repoPath)
	if err != nil {
		log.Fatal(err)
	}
	branchItr, err := repo.Branches()
	if err != nil {
		log.Fatal(err)
	}
	var branches []*plumbing.Reference
	err = branchItr.ForEach(func(branch *plumbing.Reference) error {
		branches = append(branches, branch)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	workspace := make(map[string]map[string]object.TreeEntry)

	var branchNames []string
	addToUniqueStrings := nonEmptyUniqueStringListBuilder()

	for _, branch := range branches {
		branchName := strings.TrimPrefix(branch.Name().String(), "refs/heads/")
		branchNames = append(branchNames, branchName)

		func() {
			obj, err := repo.Object(plumbing.CommitObject, branch.Hash())
			if err != nil {
				log.Fatal(err)
			}
			commit := obj.(*object.Commit)
			tree, err := commit.Tree()
			if err != nil {
				log.Fatal(err)
			}
			walker := object.NewTreeWalker(tree, true, nil)
			defer walker.Close()

			for {
				fileName, entry, err := walker.Next()
				if err != nil {
					break
				}
				addToUniqueStrings(fileName)
				if _, ok := workspace[fileName]; !ok {
					workspace[fileName] = make(map[string]object.TreeEntry)
				}
				workspace[fileName][branchName] = entry
			}
		}()
	}

	fileNames := addToUniqueStrings()()

	fileNames = filterStrings(fileNames, func(index int, name string) bool {
		return !(index < len(fileNames)-1 && strings.HasPrefix(fileNames[index+1], name+"/"))
	})

	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/", func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
		res.Header().Set("content-type", "text/html")

		html(res, func() {
			head(res, "gitctrl | multi branch view", nil)
			body(res, func() {
				table(res, func() {
					tr(res, func() {
						td(res, nil)
						for _, name := range branchNames {
							td(res, func() {
								fmt.Fprintf(res, name)
							})
						}
					})

					for _, fileName := range fileNames {
						branches := workspace[fileName]

						tr(res, func() {
							td(res, func() {
								a(res, "/"+fileName, func() {
									fmt.Fprintf(res, fileName)
								})
							})

							for _, branchName := range branchNames {
								entry, ok := branches[branchName]
								if !ok {
									td(res, nil)
									return
								}
								td(res, func() {
									fmt.Fprintf(res, entry.Hash.String()[:6])
								})
							}
						})
					}
				})
			})
		})
	})

	for _, fname := range fileNames {
		func(fileName string) {
			router.HandlerFunc(http.MethodGet, "/"+fileName, func(res http.ResponseWriter, req *http.Request) {
				html(res, func() {
					body(res, func() {
						h(res, 1, func() {
							fmt.Fprintf(res, fileName)
						})
						for _, branchName := range branchNames {
							entry, ok := workspace[fileName][branchName]
							if !ok {
								continue
							}
							obj, err := repo.Object(plumbing.BlobObject, entry.Hash)
							if err != nil {
								log.Println(err)
								continue
							}
							blob := obj.(*object.Blob)
							blobReader, _ := blob.Reader()
							if err != nil {
								log.Println(err)
								continue
							}
							defer blobReader.Close()

							h(res, 2, func() {
								fmt.Fprintf(res, branchName)
							})
							pre(res, func() {
								io.Copy(res, blobReader)
							})
						}
					})
				})
			})
		}(fname)
	}

	log.Fatal(http.ListenAndServe(":8080", router))
}

func html(res io.Writer, inner func()) {
	fmt.Fprintf(res, "<html>")
	defer fmt.Fprintf(res, "</html>")
	if inner != nil {
		inner()
	}
}

func head(res io.Writer, title string, inner func()) {
	fmt.Fprintf(res, "<head>")
	fmt.Fprintf(res, "<title>%s</title>", title)
	defer fmt.Fprintf(res, "</head>")
	if inner != nil {
		inner()
	}
}

func body(res io.Writer, inner func()) {
	fmt.Fprintf(res, "<body>")
	defer fmt.Fprintf(res, "</body>")
	if inner != nil {
		inner()
	}
}

func tr(res io.Writer, inner func()) {
	fmt.Fprintf(res, "<tr>")
	defer fmt.Fprintf(res, "</tr>")
	if inner != nil {
		inner()
	}
}

func h(res io.Writer, level int, inner func()) {
	fmt.Fprintf(res, "<h%d>", level)
	defer fmt.Fprintf(res, "</h%d>", level)
	if inner != nil {
		inner()
	}
}

func table(res io.Writer, inner func()) {
	fmt.Fprintf(res, "<table>")
	defer fmt.Fprintf(res, "</table>")
	if inner != nil {
		inner()
	}
}

func td(res io.Writer, inner func()) {
	fmt.Fprintf(res, "<td>")
	defer fmt.Fprintf(res, "</td>")
	if inner != nil {
		inner()
	}
}

func pre(res io.Writer, inner func()) {
	fmt.Fprintf(res, "<pre>")
	defer fmt.Fprintf(res, "</pre>")
	if inner != nil {
		inner()
	}
}

func style(res io.Writer, inner func()) {
	fmt.Fprintf(res, "<style>")
	defer fmt.Fprintf(res, "</style>")
	if inner != nil {
		inner()
	}
}

func a(res io.Writer, href string, inner func()) {
	fmt.Fprintf(res, "<a href=%q>", href)
	defer fmt.Fprintf(res, "</a>")
	if inner != nil {
		inner()
	}
}

func nonEmptyUniqueStringListBuilder() func(...string) func() []string {
	var list []string

	return func(strs ...string) func() []string {
		if len(strs) > 0 {
			list = append(list, strs...)
		}

		return func() []string {
			dedup := make(map[string]struct{})

			for _, str := range list {
				dedup[str] = struct{}{}
			}

			list = list[:0]
			for str, _ := range dedup {
				list = append(list, str)
			}

			sort.Strings(list)
			return list
		}
	}
}

func filterStrings(ss []string, test func(int, string) bool) []string {
	var ret []string
	for i, s := range ss {
		if test(i, s) {
			ret = append(ret, s)
		}
	}
	return ret
}
