run: 
	rm arp_list.txt results.txt || true && \
	go run main.go