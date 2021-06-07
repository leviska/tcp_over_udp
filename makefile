loss: export PUMBA := netem --duration 1h loss -p 50 network_hw_server_1
loss:
	docker-compose up

rate: export PUMBA := netem --duration 1h rate -r 40kbit network_hw_server_1
rate:
	docker-compose up

corrupt: export PUMBA := netem --duration 1h corrupt -p 50 network_hw_server_1
corrupt:
	docker-compose up
