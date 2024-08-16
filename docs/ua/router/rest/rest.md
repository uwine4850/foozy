## rest package
Даний пакет потрібен для більш зручного використання REST API.

## type DTO struct
Data Transfer Object - cтруктура для генерації об'єктів, які будуть передаватися у REST API.
Для правильної роботи вам потрібно переконатися, що дозволені повідомлення збігаються з переданими.
Важливо розуміти, що дозволені повідомлення мають використовуватися ВСІ, тобто, якщо повідомлення дозволено, воно завжди має використовуватися під час створення.
Будь-які залежності також мають бути включені в розв’язані повідомлення, і вони мають бути в тому самому файлі, що й батьківський об’єкт.

Отже, для того щоб згенерувати REST інтерфейси typescript потрібно дотриматись наступних кроків:
* Додати дозволені повідомлення. Вони передаються у вигляді AllowMessage. А точніше, це просто назва пакету та структури. Навіть залежності потрібно додати сюди.
* Використати усі дозволені повідомлення у методі `Messages`.
* Запустити метод `Generate`.

Базовий приклад викистання:
```
type TT struct {
	rest.InmplementDTOMessage
	Id string
}

type Message struct {
	rest.InmplementDTOMessage
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
Дозволені повідомлення.

__Messages__
```
Messages(messages map[string]*[]irest.IMessage)
```
Повідомлення для генерації.

__Generate__
```
Generate()
```
Генерація інтерфейсів typescript із переданих повідомлень.

## type InmplementDTOMessage struct
Cтруктура, яка буде вбудована в повідомлення.
Після вбудовування фреймворк реалізує інтерфейс irest.IMessage.

## type AllowMessage struct
Використовується для передачі пакетних даних і назви повідомлення у вигляді рядка.