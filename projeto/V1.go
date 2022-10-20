package main
import ("fmt"
		"time"
		"runtime")

var acc float64 = 0

func factorial(n int) float64 {   
  var factVal float64 = 1
  if(n < 0){
    fmt.Print("Factorial of negative number doesn't exist.")    
  }else{        
    for i:=1; i<=n; i++ {
      factVal *= float64(i)  // mismatched types int64 and int
    }   
  }    
  return factVal  /* return from function*/
}

func accumulate(i int){
	acc += (1/factorial(i))
}
func main() {
  runtime.GOMAXPROCS(runtime.NumCPU())
  start := time.Now()
  const lim int = 9
    
  for i := 1; i <= lim; i++ {
    go accumulate(i)
  }

  fmt.Printf("O programa tomou %v para rodar.\n", time.Since(start))
  fmt.Println("CPU's Utilizadas: ", runtime.NumCPU())
  fmt.Println("Aprox:", acc)
}

// Projeto V2
//
// 
// O programa tomou 26.417µs para rodar.
// CPU's Utilizadas:  2
// Aprox: 1.7182815255731925
//
//O programa tomou 23.965µs para rodar.
//CPU's Utilizadas:  3
//Aprox: 1.7182815255731925
//
//O programa tomou 19.078µs para rodar.
//CPU's Utilizadas:  32
//Aprox: 1.7182815255731925
