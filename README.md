<img width=100% src="https://capsule-render.vercel.app/api?type=waving&color=9932CC&height=120&section=header"/>
<div align="center">  
<h1> Sistema de Recarga Distribuida de Veiculos Elétricos Baseado em Blockchain </h1>
 </div>

 <div align="center">  
  <img width=20% src="http://img.shields.io/static/v1?label=STATUS&message=FINALIZADO&color=9932CC&style=for-the-badge"/>
</div>

<p align="center"> O projeto desenvolveu um sistema distribuído com blockchain para gerenciar estações de recarga de veículos elétricos, garantindo registro imutável de transações e sincronização entre nós via API REST, validando a viabilidade da tecnologia em aplicações descentralizadas. </p>

##
A crescente adoção de veículos elétricos exige soluções eficientes para o gerenciamento de estações de recarga.
Nesse contexto, tecnologias como blockchain e sistemas distribuídos oferecem vantagens como segurança, transparência e descentralização no controle de transações.
Este projeto apresenta a concepção e implementação de um sistema distribuído baseado em blockchain para simular a operação de estações de recarga de veículos elétricos. 
A aplicação foi desenvolvida em GoLang, aproveitando seus recursos nativos de concorrência, e dockerizada para permitir a simulação de uma rede peer-to-peer com múltiplos nós.
O sistema utiliza um livro-razão distribuído (blockchain) para registrar transações de reserva, recarga e pagamento de forma imutável. 
A sincronização entre os nós é garantida por uma API REST e um mecanismo de consenso simplificado, assegurando a consistência dos dados na rede simulada.

## Como Executar

    1. Para buildar as imagens do projeto use:
        docker-compose build  
    2. Para criar os containers sem iniciar:
        docker-compose create 
    3. Para verificar os containers criados, use o comando: 
        docker ps -a 
    4. Para iniciar os servidores
        docker-compose start server1 server2 server3
    Obs.: Para facilitar o processo, podem ser usados os scripts
    5. Para rodar o cliente:
        docker-compose start client
        docker exec -it client sh
        ./client


## Equipe

<table>
  <tr>
    <td align="center"><img style="" src="https://avatars.githubusercontent.com/u/144626169?v=4" width="100px;" alt=""/><br /><sub><b> Helena Filemon </b></sub></a><br />👨‍💻</a></td>
    <td align="center"><img style="" src="https://avatars.githubusercontent.com/u/143294885?v=4" width="100px;" alt=""/><br /><sub><b> Sara Souza </b></sub></a><br />👨‍💻</a></td>
  </tr>
</table>
