package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func nonExt(path string) string {
	for i := len(path) - 1; i >= 0 && !os.IsPathSeparator(path[i]); i-- {
		if path[i] == '.' {
			return path[:i]
		}
	}
	return path
}

func main() {
	localPath := `D:\Video\Movies`
	var raw_fis []os.FileInfo
	{
		// Open the directory to read its contents:
		df, err := os.Open(localPath)
		if err != nil {
			log.Fatalln(err)
			return
		}
		defer df.Close()

		// Read the directory entries:
		raw_fis, err = df.Readdir(0)
		if err != nil {
			log.Fatalln(err)
			return
		}
	}

	var exact_garbage = []string{"DD5.1", "AC3 5.1",
		"www.torentz.3xforum.ro", "ResourceRG by Dusty",
		"{1337x}", "- CODY", "anoXmous_"}

	var fuzzy_garbage = []string{"x264", "h264", "xvid",
		"1080p", "720p", "480p",
		"ac3", "eng", "aac", "dvdscr",
		"dvdrip", "brrip", "bdrip", "bluray", "web-dl", "hdtv", "dd5", "3li",
		"yify", "-rarbg", "-mgb", "-spk", "-jyk", "bluelady", "-ctrlhd", "-vlis",
		"-publichd", "gopo", "-timpe", "-sc4r", "gaz", "-axxo", "wunseedee", "hive-cm8",
		"hq", "anoxmous", "-anarchy", "-MAXSPEED", "-NuMy", "-FxM", "d-z0n3", "-Noir",
		"-ExtraTorrentRG"}

	for _, fi := range raw_fis {
		if fi.IsDir() {
			continue
		}
		name := fi.Name()

		// Skip non-video extensions:
		ext := strings.ToLower(filepath.Ext(name))
		if ext != ".mp4" && ext != ".avi" && ext != ".mkv" && ext != ".m4v" && ext != ".mpg" && ext != ".divx" {
			continue
		}
		name = nonExt(name)

		// Remove exact garbage substrings from names:
		for {
			unrecognized := true
			for _, garbage := range exact_garbage {
				i := strings.Index(name, garbage)
				if i >= 0 {
					name = name[0:i] + name[i+len(garbage):]
					unrecognized = false
					break
				}
			}
			if unrecognized {
				break
			}
		}

		// Clean up names to use spaces for separators if no spaces are found:
		{
			name = strings.Replace(name, ".", " ", -1)
			name = strings.Replace(name, "_", " ", -1)
			name = strings.Replace(name, "[", " ", -1)
			name = strings.Replace(name, "]", " ", -1)
			name = strings.Replace(name, "(", " ", -1)
			name = strings.Replace(name, ")", " ", -1)
			name = strings.TrimRight(name, " -_.")
			name = strings.Replace(name, "  ", " ", -1)
		}

		// Extract film year:
		film_year := ""
		for i := 0; i <= len(name)-4; i++ {
			year_text := name[i : i+4]
			year, err := strconv.Atoi(year_text)
			if err == nil && year >= 1900 && year <= 2032 {
				film_year = year_text
				part1 := strings.TrimSpace(name[0:i])
				part2 := strings.TrimSpace(name[i+4:])
				if len(part2) > 0 {
					name = part1 + " " + part2
				} else {
					name = part1
				}
				break
			}
		}

		// Rip out useless garbage from the tail of the filename:
		for {
			unrecognized := true
			for _, garbage := range fuzzy_garbage {
				garbage = strings.ToLower(garbage)
				nameLower := strings.ToLower(name)
				if strings.HasSuffix(nameLower, garbage) {
					name = strings.TrimRight(name[0:len(name)-len(garbage)], " -_.")
					unrecognized = false
					break
				} else if strings.HasSuffix(nameLower, "["+garbage+"]") {
					name = strings.TrimRight(name[0:len(name)-len(garbage)], " -_.")
					unrecognized = false
					break
				}
			}
			if unrecognized {
				break
			}
		}

		_ = film_year
		fmt_name := strings.TrimSpace(name)
		if film_year != "" {
			fmt_name += " [" + film_year + "]"
		}
		fmt.Println(fmt_name + ext)
	}
}
