#include "../benchmark-framework/include/lib.h"
#include <stdio.h>
#include <unistd.h>

int main(void) {
    benchmark_data data;
    benchmark_init(&data);

    char *a[4];

    for (int i = 0; i < 4; i++) {
        benchmark_capture_memory_datapoint(&data);
        a[i] = malloc(1024 * 1024 * 10);
        for (int j = 0; j < 1024 * 1024 * 10; j += 1)
            a[i][j] = 'a';
        sleep(1);
    }

    for (int i = 0; i < 4; i++) {
        free(a[i]);
        benchmark_capture_memory_datapoint(&data);
        sleep(1);
    }

    for (int i = 0; i < data.memory_datapoints.size; i++) {
        memory_datapoint *mem_datapoint =
            (memory_datapoint *)data.memory_datapoints.data[i];

        printf("Used memory: %ld, time: %ld\n", mem_datapoint->used_mem,
               mem_datapoint->time);
    }

    benchmark_end(&data);

    return 0;
}