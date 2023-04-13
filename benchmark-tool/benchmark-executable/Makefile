build: main.o benchmark_lib
	gcc -o benchmark_executable main.o ../benchmark-framework/src/lib.o ../benchmark-framework/src/container.o -lm

main.o: main.c

benchmark_lib:
	$(MAKE) -C ../benchmark-framework/src build

clean:
	rm main.o benchmark_executable
	$(MAKE) -C ../benchmark-framework/src build clean