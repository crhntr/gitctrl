package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	git "gopkg.in/src-d/go-git.v4"
)

func remoteOriginMustSSH(wd string) {
	if err := filepath.Walk(wd, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		stat, err := os.Stat(filepath.Join(path, ".git"))
		if err != nil || !stat.IsDir() {
			return nil
		}

		repo, err := git.PlainOpen(path)
		if err != nil {
			if err == git.ErrRepositoryNotExists {
				return nil
			}
			return err
		}
		config, err := repo.Config()
		if err != nil {
			return err
		}
		remoteConfig, ok := config.Remotes["origin"]
		if !ok {
			return filepath.SkipDir
		}
		for i := range remoteConfig.URLs {
			if !strings.HasPrefix(remoteConfig.URLs[i], "https://github.com/") {
				continue
			}
			fmt.Print(remoteConfig.URLs[i])
			remoteConfig.URLs[i] = strings.TrimPrefix(remoteConfig.URLs[i], "https://github.com/")
			remoteConfig.URLs[i] = "git@github.com:" + remoteConfig.URLs[i]
			fmt.Println(" --> " + remoteConfig.URLs[i])
		}
		configBits, err := config.Marshal()
		if err != nil {
			return err
		}
		configFilepath := filepath.Join(path, ".git", "config")
		os.Remove(configFilepath)
		configFile, err := os.Create(configFilepath)
		if err != nil || !stat.IsDir() {
			return err
		}
		defer configFile.Close()
		if _, err := configFile.Write(configBits); err != nil {
			return err
		}
		return filepath.SkipDir
	}); err != nil {
		log.Fatal(err)
	}
}
