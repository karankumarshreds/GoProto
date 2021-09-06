# Protocol Buffers and gRPC in Go

We will cover:

- Protocol Buffers
- Protocol Buffer Language
- Compiling a protocol buffer with `protoc`
- GRPC

---

## Protocol Buffers

Protocol buffers is a data exchange format similar to JSON. Protobuf vs JSON Both are for serialization, however, the key difference is that Protobuf is binary data –interchange format whereas JSON stores data in human-readable text format.

Protocol Buffers are strongly typed. As per definition by Google: _You can update your datastructure without breaking deployed programs that are compiled against the old datastructure format_. Interesting right? It is, we will see how that works.

Protocol Buffers allow us to define the data contract between multiple systems. Once a proto-buff file has beenn defined, we can compile it to a `target programming language`. The output of the compilation will be `classes` and `functions` of the target programming language. In Go, protocol buffers can be transported over different transports, such as HTTP/2
and Advanced Message Queuing Protocol (AMQP).

![](2021-09-05-16-00-12.png)

---

## Protocol Buffer Language

A protocol buffer is a file, which when compiled, will generate a file accessible to the targetted programming language. In go, it will be a `.go` file which will be a `structs mapping`.

Let us write a simple message in protobuf:

```proto
syntax 'proto3'

message UserInterace {
  int      index         = 1;
  string   firstName     = 2;
  string   lastName      = 3;
}
```

Here we just define a message type called `UserInterface`. It we were to write the same using JSON:

```json
{
  "index": 0,
  "firstName": "John",
  "lastName": "Doe"
}
```

The field names are changed to comply with the JSON style guide, but the essence and structure are the same.

**But, what are the sequential numbers(1,2,3) in the protobuf file?**

Those are the ordering tags used to `serialize` and `deserialize` proto-bufs between two systems. It tells the system to write the data in that particular order with the specified types. So when this proto-buf is compiled for targetted language `Go`, it will be a struct with the **empty default values**.

Different types that are used in protobuf:

- Scalar
- Enumerations
- Nested

### Scalar Values

Exmaples of scalar type values: int, int32, int64, string, bool etc. These types are converted to the corresponding language types after the compilation.

Since we are using Go in our case, the equivalent types in `Go` for these scalar types are:

![](2021-09-05-16-34-30.png)

### Enumerations (enums)

Lets take an example of a `proto` file using enums:

```proto
syntax 'proto3'
message Schedule {
  enum Days {
    SUNDAY = 0;
    MONDAY = 1;
    FRIDAY = 2;
  }
}
```

And in case we want to assign the same values to multiple enum members:

```proto
syntax 'proto3'
message Schedule {
  enum Days {
    option allow_aias = true; 👈
    UNKNOWN  = 0;
    ACTIVE   = 1;
    INACTIVE = 1;
  }
}
```

**repeated** field is equavalent to array/list:

```proto
message CarInfo {
  string type = 1;
  repeated string cars = 2;
}
```

This means, the `cars` could be `[]string` like `["bmw", "toyota", "honda"]`.

### Nested Fields

We can also use message as a type for another message.

Example:

```proto
message User{
  string firstName  = 1;
  string lastName  = 2;
  repeated Comment comments= 3;
}

message Comment{
  int id = 1;
}
```

---

## Compiling Protocol Buffer with protoc

To transfer data between the systems, we will make use of the compiled files from the `proto` proto files. We will make use of the structs (gotten from the compiled files) to create the binary data.

Steps we will follow:

1. Install the `protoc` comman-line tool and the `proto` library
2. Write `.proto` file
3. Compile the file for `Go` target language
4. Import the structs from the generated file and create data using those
5. Serialize the data into binary format and send it to the other system
6. On the remote machine, de-serialize the data and consume it

We will install he protoc command line from `https://github.com/protocolbuffers/protobuf/releases`

I am on windows, so I will install download the latest stable package and set the environment variable from `Edit the system environment variable setting` in the **control panel**.

Download and unzip (for Windows) : https://github.com/protocolbuffers/protobuf/releases/download/v3.17.3/protoc-3.17.3-win64.zip

- I will extract it in the `C:\Program Files\protoc` folder (by creating the protoc folder). It must now have `bin` folder.
- Now search for `Edit the system environment variable setting` on start menu, double click on `Path` in the list.
- Now click on `New` and then `Browse...`
- Browse until `C:\Program Files\protoc\bin` and click `ok`.
- Now open a fresh terminal and type `protoc --version` to confirm the installation. The output must be:

```bash
$ protoc --version
libprotoc 3.17.3
```

### **For MAC and Linux**

**On mac**

$ brew install protobuf

**On Ubuntu or Linux, we can copy protoc to the /usr/bin folder:**

Make sure you grab the latest version

`curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.17.3/protobuf-all-3.17.3.zip`

Unzip

`unzip protoc-3.11.3-linux-x86_64.zip -d protoc3`

Move only protoc\* to /usr/bin/

`sudo mv protoc3/bin/protoc /usr/bin/protoc`

After installing and setting up the command line tool, make sure you are able to access it from your terminal:

```bash
$ protoc --version
libprotoc 3.17.3
```

Now in our project directory, create a new folder `protofiles` and then in it let us create a proto file representing/modeling a person's information:

`person.proto`

```proto
syntax = "proto3"
package protofiles;
option go_package = "./";

message Person {
  string     name     =    1;
  int32      id       =    2;
  string     email    =    3;
  repeated   PhoneNumber  phones = 4;

  enum PhoneType {
    MOBILE  =  0;
    HOME    =  1;
    WORK    =  2;
  }

  message PhoneNumber {
    string      number  = 1;
    PhoneType   type    = 2;
  }
}

message AddressBook {
  repeated   Person  people = 1;
}
```

So we just created two main messages called Person and AddressBook. The AddressBook contains the list of Persons.
A person has a name, id, email and list of PhoneNumbers.

The second line `package protofiles` is the package name for go to compile.

To compile our person.proto proto-buf file, `cd` to the protofiles directory and run:

1. `go get -u github.com/golang/protobuf/protoc-gen-go`
   _protoc-gen-go is a plugin for the Google protocol buffer compiler to generate Go code._

2. `protoc --go_out=. *.proto`
   _this will create a Go target file in the present working directory from where the command is run (the dot) and make use of all the proto files to do so_

This will generate a `person.pb.go` file. If you open this file you will see that it contains the auto generated code for us.
This will have multiple `getter` and `setter` methods for the structs to get and set the values.

Once done, **make sure you push the repository on github to import the package in the main.go file**.

Now let us write code to create `Person` struct from the auto generated file/package(`person.pb.go`) and serialize it into a `buffer string` using the `proto` package.

> main.go

```Go
package main

import (
	"fmt"
	pb "github.com/karankumarshreds/GoProto/protofiles"
	"google.golang.org/protobuf/proto"
)

func main() {

	// using the profo created struct
	p := &pb.Person{
		Id: 1234,
		Name: "John Doe",
		Email: "test@test.com",
		Phones: []*pb.Person_PhoneNumber{
			{Number: "555-444", Type: pb.Person_HOME},
		},
	}

	// Serializing the struct and assigning it to body
	body, _ := proto.Marshal(p)

	// De-serializing the body and saving it to p1 for testing
	p1 := &pb.Person{}
	_ = proto.Unmarshal(body, p1)

	fmt.Println("Original struct loaded from proto file:", p)
	fmt.Println("Marshalled proto data: ", body)
	fmt.Println("Unmarshalled struct: ", p1)

}
```

Let us run the code now using `go run main.go`:

```
$ go run main.go
Original struct loaded from proto file: name:"John Doe" id:1234
email:"test@test.com" phones:{number:"555-444" type:HOME}


Marshalled proto data:  [10 8 74 111 104 110 32 68 111 101 16
210 9 26 13 116 101 115 116 64 116 101 115 116 46 99 111 109 34
11 10 7 53 53 53 45 52 52 52 16 1]

Unmarshalled struct:  name:"John Doe" id:1234 email:"test@test.com"
phones:{number:"555-444" type:HOME}
```

The second output is the binary bytes which has been serialized into by the `proto` library. This serialized binary data needs a medium to move between the two or more systems. This is where `gRPC` kicks in. Using gRPC (Google Remote Procedure Call), a server and client(not frontend) can talk which each other in the protocol buffer format.

---

## GRPC

GRPC is a transport mechanism that sends and receives messages (protocol buffers) between two systems. The two parties are referred to as `server` and `client`.

**The main advantage of gRPC is that it can be understood by multiple programming languages (via their respective grpc packages mostly) making the communication easy between different tech stacks**

We need to install the grpc Go library and a protoc-gen plugin before writing the
services. Install them using the following commands:

`go get google.golang.org/grpc`

`go get -u github.com/golang/protobuf/protoc-gen-go` (already did in the previous section)

In this section we will create a money transaction service which will communicate over GRPC:

1. Create a proto-buf with the definistion of service and messages
2. Compile the protocol buffer file
3. Use the generated file package to create a gRPC server
4. Create a gRPC client to talk to the server

We will follow the SAME steps to create protobuf and compile package as we did in the last example.

Let us create a `transaction.proto` file nad put it in the protofiles directory.

```proto
syntax = "proto3";
package protofiles;
option go_package = "./";

message TransactionRequest {
  string    from   =   1;
  string    to     =   2;
  float     amount =   3;
}

message TransactionResponse {
  bool     confirmation = 1;
}

service MoneyTransaction {
  rpc MakeTransaction(TransactionRequest) returns (TransactionResponse) {}
}
```

The new keyword `service` defines the GRPC service. Here we are defining another type (type of function) for our RPC system, which will take in a type of TransactionRequest and return a TransactionResponse. Once compiled using `protoc`, the compiled file will contain an interface to invoke this function.

Read more here: https://developers.google.com/protocol-buffers/docs/proto3#services

Now let us `cd` into the parent directory to the protofiles directory and run the command to compile the `transaction.proto`:

```
$ ls

protofiles/

$ protoc -I protofiles/ protofiles/transaction.proto --go_out=plugins=grpc:protofiles

```

**NOTE** This time we have used a `grpc` plugin to compile the package.

Now `ls protofiles/` to confirm if the `go` package has been compiled:

```
-a----        06-09-2021     16:13           7672 transaction.pb.go
-a----        06-09-2021     16:13            357 transaction.proto
```

Now let us setup the GRPC server. The code explanation is done in the code comments:

Create a folder and file `server/main.go`

```Go
package main

import (
	"fmt"
	"net"
	"log"
	pb "github.com/karankumarshreds/GoProto/protofiles"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type server struct{}

func main() {

	// NewServer creates a gRPC server which has no service registered and has not started
	// to accept requests yet.
	s := grpc.NewServer()
	lis, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// We are making use of the function that compiled proto made for us to register
	// our GRPC server so that the clients can make use of the functions tide to our
	// server remotely via the GRPC server (like MakeTransaction function)

	// The first argument is the grpc server instance
	// The second argument is the service who's methods we want to expose (in our case)
	// we have put it in this program only
	pb.RegisterMoneyTransactionServer(s, &server{})

	// Serve accepts incoming connections on the listener lis, creating a new ServerTransport
	// and service goroutine for each. The service goroutines read gRPC requests and then
	// call the registered handlers to reply to them.
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}

// [ctx] is used by the goroutines to interact with GRPC
// [in] is the type of TransactionRequest
/*
	This function signature matches the service that we mentioned in the protobuf
*/
func (s *server) MakeTransaction(ctx context.Context, in *pb.TransactionRequest) (*pb.TransactionResponse, error) {
	// Business logic will come here
	fmt.Println("Got amount ", in.Amount)
	fmt.Println("Got from ", in.From)
	fmt.Println("For ", in.To)
	// Returning a response of type Transaction Response
	return &pb.TransactionResponse{Confirmation: true}, nil
}
```

This is it for the server side. The idea is that if the client makes a request to invoke the `MakeTransaction` function remotely, the server should go ahead and execute it (here it will only print the data received).

Now let us write down code for the client side as well:

Folder and file : `client/main.go`

```go
package main

import (
	"log"
	pb "github.com/karankumarshreds/GoProto/protofiles"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// grpc server address
const address = "localhost:8000"

func main() {
	// Set up connection with the grpc server
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error while making connection, %v", err)
	}

	// Create a client instance
	c := pb.NewMoneyTransactionClient(conn)

	// Lets invoke the remote function from client on the server
	c.MakeTransaction(
		context.Background(),
		&pb.TransactionRequest{
			From: "John",
			To: "Alice",
			Amount: float32(120.15),
		},
	)
}
```

This client code will make connection to the grpc server we created through the GRPC package and make use of the protofile to run the `MakeTransaction` function remotely.

Now if we run both together. First we will run the server code:

`cd server`
`go run main.go`

This will start the server code, now open **a new terminal** and start the client code as well. As soon as we do that, we should see the logs on the `server` side terminal.

`cd client`
`go run main.go`

That's it, the client now made a connection with the server and invoked the `MakeTransaction` function remotely.
The server must have the following logs:

```
$ go run main.go
Got amount  120.15
Got from  John
For  Alice
```

Congratulations, you have successfully created the server and client side grpc connection. A gRPC client can request a gRPC server to perform a computationheavy/secure operation. The client can be a mobile device too.
