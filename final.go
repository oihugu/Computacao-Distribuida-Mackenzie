package main

import (
	"fmt"
	"sync"
)

// O objetivo do programa era separar a função de taylor em diferentes parciais, sendo que cada uma poderia ser executada de maneira semi-independente
// Aproveitando das caracteristicas das go routines
// A ordem das funções no programa se encontra do menor nível de execução para o maior

func partial_factorial(n_prev int64, n int64, c chan int64, wg *sync.WaitGroup) {
	// Primeiramente fazemos uma função parcial do fatorial
	// Aproveitando que podemos dividir a multiplicação em "fatias"
	// Dividimos o calculo de 10! por exemplo no calculo de (1*2*3*4*5)*(6*7*8*9*10)
	var partial_product int64
	partial_product = 1
	for i := n_prev; i <= n; i++ {
		if i > 1 {
			partial_product *= int64(i)
		}
	}
	c <- partial_product
	wg.Done()
	// Para realizar o controle desta função utilizamos uma estratégia de canais com work group
	// Canais foram utilizados para realizar a comununicação entre as duas goroutines
	// Aproveitando da caracteristica da função fatoria que não depende da ordem em que as multiplicações são feitas
	// E Workgrups foram utilizados para saber quando todos os números haviam sido calculados
}

func factorial(n int64, prev_map *map[int64]int64, prev_holder *int64, mut *sync.Mutex) int64 {
	// A função fatorial cria os Wait Groups e os Canais, porem a mesma é rodada de maneira sequencial
	// O maior ponto de atenção interno a função é a área de exclusão mútua existente no fim da função
	// Esta área existe para controlar o acesso a um dicionário de fatoriais já gerados
	// A ideia era manter uma referência de quais já foram gerados e que isso servisse como uma base para gerar os próximos
	// Evitando assim, gastos com computação já realizada
	prev_holder_exec := *prev_holder
	factorial_n := (*prev_map)[prev_holder_exec]

	var wg sync.WaitGroup
	var channel_size int64
	if n > prev_holder_exec {
		channel_size = (n - prev_holder_exec) / 5
	} else {
		channel_size = n / 5
	}
	channel_products := make(chan int64, int(channel_size)+1)

	if n > prev_holder_exec {
		wg.Add(int((n-prev_holder_exec)/5) + 1)
		for i := prev_holder_exec + 1; i <= n; i++ {
			if i%5 == 0 {
				go partial_factorial(i-4, i, channel_products, &wg)
			}
		}
	} else {
		wg.Add(int(n / 5))
		for i := int64(1); i <= n; i++ {
			if i%5 == 0 {
				go partial_factorial(i-4, i, channel_products, &wg)
			}
		}
	}

	wg.Wait()
	for v := range channel_products { //Aqui
		factorial_n *= v
	}
	mut.Lock()
	if (n-prev_holder_exec > 10) && (n%10 == 0) {
		(*prev_map)[n/10] = factorial_n
	}
	mut.Unlock()

	return factorial_n
	// Provavelmente o erro que impede o programa de rodar está acontecendo aqui
	// Mais especificamente no for que tenta iterar sobre o canal de produtos
	// Já que o mesmo recebe dois tipos diferentes de dados do if imediatamente superior
}

func partial_taylor(n int64, prev_map *map[int64]int64, prev_holder *int64, mut *sync.Mutex, c chan float64, wg *sync.WaitGroup) {
	// A parcial da função de taylor é relativamente simples
	// Basicamente a mesma cria um acumulador ao qual serão somados um subset específico de resultados da função de taylor
	var acc float64 = 0
	for i := n - 4; i <= n; i++ {
		acc += 1 / float64(factorial(i, prev_map, prev_holder, mut))
	}
	c <- acc
	wg.Done()
}

func calculate_taylor(n int64, prev_map *map[int64]int64, prev_holder *int64, mut *sync.Mutex) float64 {
	// A função principal cria os principais canais de comunicação descritos anteriormente
	// Como por exemplo o acumulador da função de taylor e um Wait Group que nos indica quando todas as intânicas
	// Da parcial da função de taylor foram executadas, além disso também cria a região de memória
	// Mutuamente exclusiva, utilizada para segurar os valores de fatorial já calculados
	// Além disso acumula os valores deixados no canal pela parcial de taylor e retorna o resultado
	var taylor_accumulator float64
	var wg sync.WaitGroup

	taylor_accumulator = 1
	taylor_channel := make(chan float64, int(n%5))

	wg.Add(int(n / 5))

	for i := int64(0); i <= n; i++ {
		if i%5 == 0 {
			go partial_taylor(i, prev_map, prev_holder, mut, taylor_channel, &wg)
		}
	}

	wg.Wait()

	for v := range taylor_channel {
		taylor_accumulator += v
	}

	return taylor_accumulator
}

func main() {
	prev_map := make(map[int64]int64)
	var prev_holder int64
	var prev_mutex sync.Mutex

	fmt.Println(calculate_taylor(100, &prev_map, &prev_holder, &prev_mutex))
}
