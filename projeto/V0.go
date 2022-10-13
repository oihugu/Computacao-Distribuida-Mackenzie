package main
import "fmt"
import "time"

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

func main() {
  start := time.Now()
  const lim int = 9
  var acc float64 = 0
    
  for i := 1; i <= lim; i++ {
    acc += (1/factorial(i))
  }

  fmt.Printf("O programa tomou %v para rodar.\n", time.Since(start))
  fmt.Println("Aprox:", acc)
}
