#include "../include/lib.h"
#include <stdio.h>
#include <string.h>
#include <time.h>
#include <sys/socket.h>
#include <netinet/in.h>

static long get_memory_usage() {
    FILE *f = fopen("/proc/self/status", "r");
    char line[128];
    while (fgets(line, 128, f) != NULL) {
        if (strncmp(line, "VmRSS:", 6) == 0) {
            fclose(f);
            return atol(line + 7);
        }
    }
    fclose(f);
    return -1;
}

static void open_server(benchmark_data *data) {
    int sock, len;
    struct sockaddr_in saddr, cli;
    
    if ((sock = socket(AF_INET, SOCK_STREAM, 0)) < 0) {
        benchmark_error("Socket creation failed!", data);
    }

    saddr.sin_family = AF_INET;
    saddr.sin_port = htons(25565);
    saddr.sin_addr.s_addr = htonl(INADDR_ANY);

    if (bind(sock, (struct sockaddr *) &saddr, sizeof(saddr))) {
        benchmark_error("Socket bind failed!", data);
    }

    if (listen(sock, 1)) {
        benchmark_error("Socket listen failed!", data);
    }

    len = sizeof(cli);

    printf("Waiting for client to connect...\n");
    if((data->comm_socket = accept(sock, (struct sockaddr *) &cli, &len)) < 0) {
        benchmark_error("Socket accept failed!", data);
    }
    send(data->comm_socket, "booted!", 7, 0);
    printf("Client connected!\n");
}

void benchmark_init(benchmark_data *data) {
    open_server(data);

    container_init(&data->memory_datapoints);
    data->display_errors = 1;
}

void benchmark_capture_memory_datapoint(benchmark_data *data) {
    long used_mem = get_memory_usage();
    if (used_mem == -1) {
        benchmark_error("Failed to get memory usage", data);
        return;
    }

    memory_datapoint *mem_datapoint = malloc(sizeof(memory_datapoint));
    mem_datapoint->used_mem = used_mem;
    mem_datapoint->time = time(NULL);

    container_push_back(&data->memory_datapoints, (void *)mem_datapoint);
}

void benchmark_capture_cpu_usage_datapoint(benchmark_data *data) {}

void benchmark_end(benchmark_data *data) {
    container_free(&data->memory_datapoints);
}

void benchmark_set_display_errors(benchmark_data *data, int display_errors) {
    data->display_errors = display_errors;
}

void benchmark_error(char *msg, benchmark_data *data) {
    if (!data->display_errors)
        return;

    fprintf(stderr, "[benchmark-framework] %s\n", msg);
}