build: lib.o container.o
	ar rcs benchmark_lib.a $^

lib.o: lib.c

container.o: container.c

clean:
	rm -f lib.o container.o benchmark_lib.a