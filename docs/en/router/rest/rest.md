## rest package
This package is required for more convenient use of the REST API.

## type DTO struct
Data Transfer Object - structure for generating objects that will be passed to the REST API.
For proper operation, you need to ensure that the allowed messages match the ones transmitted.
It is important to understand that allowed messages must be used by ALL, meaning that if a message is allowed, it must always be used during creation.
Any dependencies must also be included in the resolved messages and must be in the same file as the parent object.

So, in order to generate REST typescript interfaces, you need to follow the following steps:
* Add allowed messages. They are passed in the form of AllowMessage. Or rather, it's just the name of the package and structure. Even dependencies need to be added here.
* Use all allowed messages in `Messages` method.
* Run the `Generate` method.

A basic example of a pull-out:
```
type TT struct {
	rest.ImplementDTOMessage
	Id string
}

type Message struct {
	rest.ImplementDTOMessage
	Id   int
	Tee  TT
	Tee1 []TT
	Tee2 map[int][]map[TT][]string
}

func main() {
	d := rest.NewDTO()
	d.AllowedMessages([]rest.AllowMessage{
		{
			Package: "main",
			Name:    "Message",
		},
		{
			Package: "main",
			Name:    "TT",
		},
	})
	d.Messages(map[string]*[]irest.IMessage{"tee.ts": {Message{}, TT{}}})
	if err := d.Generate(); err != nil {
		panic(err)
	}
}
```

__AllowedMessages__
```
AllowedMessages(messages []AllowMessage)
```
Allowed messages.

__Messages__
```
Messages(messages map[string]*[]irest.IMessage)
```
Message for generation.

__Generate__
```
Generate()
```
Generation of typescript interfaces from passed messages.

## type ImplementDTOMessage struct
The structure that will be embedded in the message.
After embedding, the framework implements the irest.IMessage interface.

## type AllowMessage struct
It is used to transfer packet data and the name of the message in the form of a string.