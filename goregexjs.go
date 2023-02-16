package main

import (
    "bufio"
    "flag"
    "fmt"
    "io"
    "net/http"
    "os"
    "regexp"
    "sync"
)

func main() {
    // cria um objeto flag e adiciona as opções
    regexFile := flag.String("regex-file", "", "caminho para o arquivo contendo as expressões regulares")
    showChars := flag.Int("show-chars", 10, "quantidade de caracteres a serem mostrados da palavra encontrada (padrão: 10)")
    numThreads := flag.Int("threads", 10, "número de threads a serem usadas (padrão: 10)")
    flag.Parse()

    // lê as URLs do arquivo de entrada padrão
    var urls []string
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        urls = append(urls, scanner.Text())
    }
    if err := scanner.Err(); err != nil {
        panic(err)
    }

    // lê as expressões regulares do arquivo
    regexes, err := readRegexes(*regexFile)
    if err != nil {
        panic(err)
    }

    // cria um canal para passar as URLs para as goroutines
    urlChan := make(chan string)

    // cria um WaitGroup para esperar todas as goroutines terminarem
    var wg sync.WaitGroup

    // cria as goroutines para processar as URLs
    for i := 0; i < *numThreads; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for url := range urlChan {
                processUrl(url, regexes, *showChars)
            }
        }()
    }

    // adiciona as URLs no canal
    for _, url := range urls {
        urlChan <- url
    }

    // fecha o canal para indicar que não há mais URLs a serem processadas
    close(urlChan)

    // espera todas as goroutines terminarem
    wg.Wait()
}

func processUrl(url string, regexes []*regexp.Regexp, showChars int) {
    resp, err := http.Get(url)
    if err != nil {
        fmt.Printf("Erro ao baixar URL %s: %s\n", url, err.Error())
        return
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        fmt.Printf("Erro ao ler resposta da URL %s: %s\n", url, err.Error())
        return
    }

    js := string(body)
    for _, regex := range regexes {
        match := regex.FindStringIndex(js)
        if match != nil {
            start := match[0]
            end := match[1]
            word := js[start:end]
            if len(word) > showChars {
                word = word[:showChars] + "..."
            }
            fmt.Printf("URL: %s\nPalavra: %s\nRegex: %s\n\n", url, word, regex.String())
        }
    }
}

func readRegexes(filename string) ([]*regexp.Regexp, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("Erro ao abrir arquivo %s: %s", filename, err.Error())
    }
    defer file.Close()

    var regexes []*regexp.Regexp
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        regex, err := regexp.Compile(scanner.Text())
        if err != nil {
            return nil, fmt.Errorf("Erro ao compilar expressão regular %s: %s", scanner.Text(), err.Error())
        }
        regexes = append(regexes, regex)
    }
    if err := scanner.Err(); err != nil {
        return nil, fmt.Errorf("Erro ao ler arquivo %s: %s", filename, err.Error())
    }
    
    return regexes, nil
}