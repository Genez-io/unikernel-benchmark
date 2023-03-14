#include "../include/container.h"

void container_init(container *c) {
    c->data = (void **)malloc(1024 * sizeof(void *));
    c->size = 0;
    c->capacity = 1024;
}

void container_init_with_capacity(container *c, int capacity) {
    c->data = (void *)malloc(capacity * sizeof(void *));
    c->size = 0;
    c->capacity = capacity;
}

void container_push_back(container *c, void *value) {
    if (c->size == c->capacity) {
        c->capacity *= 2;
        c->data = (void **)realloc(c->data, c->capacity * sizeof(void *));
    }
    c->data[c->size++] = value;
}

void container_free(container *c) {
    for (int i = 0; i < c->size; i++)
        free(c->data[i]);
    free(c->data);
}