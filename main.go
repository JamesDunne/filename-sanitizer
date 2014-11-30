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

	var exact_garbage = []string{"DD5.1", "AC3 5.1", "AC3-5.1",
		"www.torentz.3xforum.ro", "ResourceRG by Dusty", " by WingTip",
		"{1337x}", "- CODY", " - Ozlem", "anoXmous_0", "anoXmous_", "-FASM", "_Kuth",
		"(Opt.SWESUBS)", "[big_dad_e™]", "_sujaidr", " - IMAGiNE"}

	var fuzzy_garbage = []string{
		// Format:
		"x264", "h264", "xvid",
		"1080p", "720p", "480p",
		"ac3", "eng", "aac", "aac51", "dvdscr", "MP4", "QEBS5", "STEREO",
		"dvdrip", "DVD-Rip", "brrip", "bdrip", "bluray", "web-dl", "hdtv", "dd5", "3li",
		"divx3lm", "divx", "DTS",
		"DXVA", "hq",

		// Signatures:
		"yify", "-rarbg", "-mgb", "-spk", "-jyk", "bluelady", "-ctrlhd", "-vlis", "-PCA",
		"-publichd", "gopo", "-timpe", "-sc4r", "gaz", "-axxo", "wunseedee", "hive-cm8",
		"anoxmous", "-anarchy", "-MAXSPEED", "-NuMy", "-FxM", "d-z0n3", "-Noir",
		"-ExtraTorrentRG", "-Nile", "-Stealthmaster", "BOKUTOX", "-WiKi", "-AMIABLE",
		"-SiNiSTER", "-SHiRK", "-DVL", "_Kuth", "R5", "-HuMPDaY", "-FLAWL3SS", "-BH",
		"-HHAH", "-iNiQUiTY", "-DAH", "-MXMG", "-haSak", "-Ekolb", "-ESiR", "Blood",
		"-PsiX", "-CC", "RC",
	}

	var fuzzy_skip = []string{
		"cd1", "cd2",
	}

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
					if i+len(garbage) < len(name) {
						name = name[0:i] + name[i+len(garbage):]
					} else {
						name = name[0:i]
					}
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
				}
			}
			for _, skip := range fuzzy_skip {
				skip = strings.ToLower(skip)
				nameLower := strings.ToLower(name)
				if strings.HasSuffix(nameLower, skip) {

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
