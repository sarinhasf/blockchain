# se for rodar no linux, autorizar o script
# chmod +x ./inicia.sh

docker-compose build
docker-compose create

# Mostrar todos os containers
docker ps -a

# Iniciando servidores
docker-compose start server1 server2 server3

# Inicia cliente
#docker-compose start client
#docker exec -it client sh
#para rodar o codigo dentro do terminal do cliente, rode ./client
#cat blockchain.json