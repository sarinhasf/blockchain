<img width=100% src="https://capsule-render.vercel.app/api?type=waving&color=9932CC&height=120&section=header"/>
<div align="center">  
<h1> Sistema de Recarga Distribuida de Veiculos Elétricos Baseado em Blockchain </h1>
 </div>

 <div align="center">  
  <img width=20% src="http://img.shields.io/static/v1?label=STATUS&message=FINALIZADO&color=9932CC&style=for-the-badge"/>
</div>

<p align="center"> O projeto desenvolveu um sistema distribuído com blockchain para gerenciar estações de recarga de veículos elétricos, garantindo registro imutável de transações e sincronização entre nós via API REST, validando a viabilidade da tecnologia em aplicações descentralizadas. </p>

## Sumário

- [Introdução](#introdução)
- [Fundamentos Teóricos](#fundamentos-teóricos)
- [Arquitetura do Sistema](#arquitetura-do-sistema)
- [Como Executar](#como-executar)
- [Conclusão](#conclusão)
- [Referências](#referências)



## Introdução
O aumento do uso de Veículos Elétricos (VEs) no Brasil exige maior segurança e transparência nas transações de recarga. Este relatório apresenta um sistema distribuído, construído com GoLang e Docker, que utiliza blockchain para registrar e validar operações de reserva, recarga e pagamento, mitigando fraudes e aprimorando a confiança entre usuários e empresas.



## Fundamentos Teóricos
- **Linguagem GoLang**
Go é uma linguagem de programação compilada, concorrente e com tipagem estática. Ela é conhecida por sua simplicidade, performance e robustez, sendo ideal para o desenvolvimento de sistemas distribuídos e de rede. A concorrência em Go é facilitada pelas goroutines e channels, que permitem a execução simultânea de tarefas de forma eficiente, sem a complexidade de threads tradicionais.
- **Docker**
Docker é uma plataforma de código aberto que permite aos desenvolvedores automatizar a implantação, escala e gerenciamento de aplicações dentro de contêineres. Contêineres são unidades de software padronizadas que empacotam código e todas as suas dependências, garantindo que a aplicação funcione de forma consistente em diferentes ambientes. O docker-compose.yml define e executa aplicações Docker multi-contêiner, simplificando a orquestração de serviços como server1, server2, server3 e client.
- **Sistemas Distribuídos**
Sistemas distribuídos são coleções de computadores autônomos interconectados que trabalham juntos como um único sistema coerente para os usuários. Eles oferecem vantagens como escalabilidade, tolerância a falhas e maior disponibilidade. No contexto deste projeto, os múltiplos servidores (server1, server2, server3) operam como nós de um sistema distribuído, comunicando-se e sincronizando informações para manter uma cadeia de blocos consistente.
- **API REST**
API REST (Representational State Transfer) é um estilo arquitetural para sistemas distribuídos que utiliza o protocolo HTTP. APIs RESTful são sem estado, o que significa que cada requisição do cliente para o servidor contém todas as informações necessárias para entender a requisição. Os endpoints como /blockchain, /add-block, /status e /mensagem são exemplos de como a API REST é utilizada para a comunicação entre os nós e o cliente neste sistema.
- **Livro-Razão (ledger) Distribuído**
Um livro-razão distribuído (Distributed Ledger Technology - DLT) é um banco de dados descentralizado replicado e compartilhado entre vários participantes em uma rede. Ao contrário dos bancos de dados centralizados, não há uma autoridade central. A blockchain é um tipo de DLT, onde os registros são organizados em blocos e encadeados criptograficamente, garantindo imutabilidade e auditabilidade.
- **Blockchain**
Blockchain é uma tecnologia de livro-razão distribuído que organiza os dados em blocos e os encadeia usando criptografia. Cada bloco contém um hash do bloco anterior, uma marca de tempo, dados da transação e seu próprio hash. Essa estrutura impede a alteração de dados retroativamente, tornando-a ideal para aplicações que exigem alta segurança e transparência, como o registro de transações de veículos elétricos. Funções como calculateHash, createGenesisBlock e isBlockValid são essenciais para a integridade da cadeia.



## Arquitetura do sistema
A arquitetura do sistema é distribuída, composta por múltiplos servidores GoLang que atuam como nós de uma blockchain e um cliente GoLang para interação do usuário, sendo tudo orquestrado pelo Docker Compose. Os servidores hospedam cópias da blockchain, validam e sincronizam blocos, e expõem uma API REST para comunicação inter-nós e com o cliente. O cliente simula interações de reserva, recarga e pagamento de veículos elétricos, enviando essas transações para os servidores, que as registram como novos blocos na blockchain, garantindo imutabilidade e auditabilidade através de um mecanismo de consenso de cadeia mais longa.



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


## Conclusão
Em suma, este projeto validou a eficácia da blockchain e sistemas distribuídos, com GoLang e Docker, na criação de um sistema de recarga de VEs transparente, seguro e auditável. A solução mitiga desafios como fraudes e a falta de integração, garantindo a imutabilidade das transações e fomentando a confiança no ecossistema. Futuras evoluções podem incluir autenticação de usuários, otimização da sincronização e integração com APIs de VEs reais.


## Referências
GOLANG. The Go Programming Language Documentation. Disponível em: https://golang.org/doc/.

DOCKER INC. Docker Documentation. Disponível em: https://docs.docker.com/.



## Equipe

<table>
  <tr>
    <td align="center"><img style="" src="https://avatars.githubusercontent.com/u/144626169?v=4" width="100px;" alt=""/><br /><sub><b> Helena Filemon </b></sub></a><br />👨‍💻</a></td>
    <td align="center"><img style="" src="https://avatars.githubusercontent.com/u/143294885?v=4" width="100px;" alt=""/><br /><sub><b> Sara Souza </b></sub></a><br />👨‍💻</a></td>
  </tr>
</table>
