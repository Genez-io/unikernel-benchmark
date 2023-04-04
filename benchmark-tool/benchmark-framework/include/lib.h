#include "container.h"

typedef struct {
    container memory_datapoints;
    int display_errors;
    int comm_socket;
} benchmark_data;

typedef struct {
    long used_mem;
    long time;
} memory_datapoint;

void benchmark_init(benchmark_data *);
void benchmark_capture_memory_datapoint(benchmark_data *);
void benchmark_end(benchmark_data *data);

void benchmark_set_display_errors(benchmark_data *, int);
void benchmark_error();