<h1 align="center">Welcome to GoChat 👋</h1>
<p>
</p>

> This program is a console chat utility written in Golang that uses Protobuff and Gocui libraries. In other words it's an update of this project "https://github.com/hadi-ilies/SimpleGolangChatServer.git".
The main goal is to familiarize with those libraries.

## Screenshots
  ![Screenshots](https://i.imgur.com/ghSX2QJ.png)
  
## commands
  * ctrl + c :
        Quit the chat
  * ctrl + v :
        Switch input view to chat view and vice versa
  * ctrl + l :
        Clear input view
  * key arrow up :
        Scroll up the chat view 
  * key arrow down :
        Scroll down the chat view
  * key Enter :
        send message to the chat

## Usage

```sh
USAGE
        go run main.go [options]
OPTIONS
        -S, --server <[ipaddr], port>
                start server using the ip and port pass in argv
        -C, --client <ipaddr, port>
                start client on the ip and port pass in argv
        -h, --help
                Display the program usage
```

## Run server

```sh
  go run main.go --server <[ipaddr], port>
```

## Run Client

```sh
  go run main.go --client <ipaddr, port>
```

## Author

👤 **hadi-ilies**

* Website: https://hadibereksi.fr/
* Github: [@hadi-ilies](https://github.com/hadi-ilies)
* LinkedIn: [@https:\/\/www.linkedin.com\/in\/hadibereksi\]

## Show your support

Give a ⭐️ if this project helped you!
