package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type ayah struct {
	surahNumber int
	ayahNumber  int

	surahName string
	ayahText  string
}

func getRandomAyah() (*ayah, error) {
	ayahFile, err := os.Open("./translation.txt")
	if err != nil {
		return nil, err
	}

	randomAyahLine, err := getRandomAyahLine(ayahFile)
	if err != nil {
		return nil, err
	}

	if err := ayahFile.Close(); err != nil {
		return nil, err
	}

	ayahParts := strings.SplitN(randomAyahLine, "|", 3)
	surahNumber, err := strconv.ParseInt(ayahParts[0], 10, 0)
	if err != nil {
		return nil, err
	}
	ayahNumber, err := strconv.ParseInt(ayahParts[1], 10, 0)
	if err != nil {
		return nil, err
	}

	return &ayah{
		surahNumber: int(surahNumber),
		ayahNumber:  int(ayahNumber),

		surahName: getSurahNames()[surahNumber-1],
		ayahText:  ayahParts[2],
	}, nil
}

func getRandomAyahLine(r io.Reader) (string, error) {
	rand.Seed(time.Now().Unix())
	randomLineNumber := rand.Intn(TOTAL_AYAHS)

	lineReader := bufio.NewScanner(r)
	for i := 0; i < randomLineNumber; i++ {
		if !lineReader.Scan() {
			break
		}
	}

	if lineReader.Err() != nil {
		return "", lineReader.Err()
	}

	return lineReader.Text(), nil
}

func (a *ayah) getFooter() string {
	return fmt.Sprintf("%d:%d %s", a.surahNumber, a.ayahNumber, a.surahName)
}
