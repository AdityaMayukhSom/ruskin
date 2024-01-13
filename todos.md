### **Todos**

1. Inside `transport/producer.go` need to consider between using 
`gob.NewEncoder()` vs `json.Marshal()` for converting `TransportResponse` 
struct into slice of bytes for sending as response.

    Refer to [this](https://stackoverflow.com/questions/16330490/in-go-how-can-i-convert-a-struct-to-a-byte-array) 
stack overflow question for more information regarding speed and tradeoffs.

2. Consumer API needs to be implemented to launch it as a production application.