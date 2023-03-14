#include <stdlib.h>

typedef struct {
    void **data;
    int size;
    int capacity;
} container;

void container_init(container *c);
void container_init_with_capacity(container *c, int capacity);
void container_push_back(container *c, void *value);
void container_free(container *c);