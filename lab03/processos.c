#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

int main (int argc, char *argv[]) {
    pid_t pid;
    int i, n;
    pid_t childpid = 0;
    int fibo_num = 9;
    int result = 1;
    int vec[fibo_num];
    int slices[fibo_num/2];
    n = 2;

    for (int i = 1; i < fibo_num+1; i++){
        vec[i - 1] = i;
    }

    for (i = 1; i < n; i++){
        pid = fork();
        if (pid == 0){
            break;
        }
    }

    for (int j = (fibo_num/n) * (i-1); j < (fibo_num/n) * i; j++){
        result *= vec[j];
    }

    slices[i] = result;

    if (pid != 0){
        result = 1;
        for (int k = 0; k < fibo_num/2;k++){
            fprintf(stderr, "%d\n", slices[k]);
        }
        fprintf(stderr, "Resultado: %d", slices[1]);
    }

    return 0;
}