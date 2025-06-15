package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"sort"
	"sync"
	"time"
)

func main() {
	url := flag.String("url", "", "URL do serviço a ser testado")
	requests := flag.Int("requests", 1, "Número total de requests")
	concurrency := flag.Int("concurrency", 1, "Número de chamadas simultâneas")
	flag.Parse()

	fmt.Printf("Parâmetros recebidos: \nurl: %s \nrequests: %d \nconcurrency: %d\n\n", *url, *requests, *concurrency)

	if *url == "" || *requests <= 0 || *concurrency <= 0 {
		fmt.Println("Parametros inválidos ou ausentes")
		fmt.Println("Parâmetros obrigatórios --url <url> --requests <número> --concurrency <número>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	start := time.Now()
	fmt.Printf("Iniciando teste de requisições: %v\n", start.Format(time.RFC3339))

	var wg sync.WaitGroup
	ch := make(chan struct{}, *concurrency)
	var successCount, errorCount int
	var mu sync.Mutex
	statusCount := make(map[int]int)

	for i := 0; i < *requests; i++ {
		wg.Add(1)
		ch <- struct{}{}

		go func() {
			defer wg.Done()
			defer func() { <-ch }()

			resp, err := http.Get(*url)
			if err != nil {
				mu.Lock()
				errorCount++
				mu.Unlock()
				return
			}

			resp.Body.Close()

			mu.Lock()
			successCount++
			statusCount[resp.StatusCode]++
			mu.Unlock()
		}()
	}

	wg.Wait()
	duration := time.Since(start)
	fmt.Printf("Teste finalizado: %s\n\n", time.Now().Format(time.RFC3339))
	fmt.Println("Relatório:")
	fmt.Println("Tempo total de duração:", duration)
	fmt.Printf("Total de requisições: %d\n", *requests)
	fmt.Printf("Requisições completas: %d\n", successCount)
	fmt.Printf("Requisições com erro: %d\n", errorCount)
	fmt.Println("Detalhamento por código de status:")

	var statusOrder []int
	for status := range statusCount {
		statusOrder = append(statusOrder, status)
	}

	sort.Ints(statusOrder)

	for _, statusCode := range statusOrder {
		fmt.Printf("Status %d - %s - %d ocorrências\n", statusCode, http.StatusText(statusCode), statusCount[statusCode])
	}

}
