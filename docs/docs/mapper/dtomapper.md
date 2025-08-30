## dtomapper

#### DeepCheckDTOSafeMessage
Checks whether transmitted messages and internal messages are safe.
That is, there will be a check of internal structures, in depth to the limit.
It is mandatory to have “dto” tags for each field.

__IMPORTANT__: `messagePtr` accepts a pointer to a structure (preferably) or a pointer to 
a structure interface. In both cases, the object must implement the irest.IMessage interface.
```golang
func DeepCheckDTOSafeMessage(dto *rest.DTO, messagePtr typeopr.IPtr) error {
	if err := rest.IsSafeMessage(messagePtr, dto.GetAllowedMessages()); err != nil {
		return err
	}

	message := reflect.ValueOf(messagePtr.Ptr()).Elem()
	var RV reflect.Value
	if message.Type().Kind() == reflect.Interface {
		RV = message.Elem()
	} else {
		RV = message
	}
	rawObject := LoadSomeRawObjectFromCache(RV, &messageRawCache, namelib.TAGS.DTO)
	for _, f := range *rawObject.Fields() {
		if f.Type.Kind() == reflect.Struct && f.Type != implementDTOMessageType && f.Type != typeIdType {
			v := RV.FieldByName(f.Name)
			i := v.Interface().(irest.Message)
			if err := DeepCheckDTOSafeMessage(dto, typeopr.Ptr{}.New(&i)); err != nil {
				return err
			}
		}
	}
	return nil
}
```

#### JsonToDTOMessage
Converts JSON data into the selected message.
It is important that the message is safe.
```golang
func JsonToDTOMessage[T any](jsonData map[string]interface{}, dto *rest.DTO, output *T) error {
	if err := DeepCheckDTOSafeMessage(dto, typeopr.Ptr{}.New(output)); err != nil {
		return err
	}
	if err := FillDTOMessageFromMap(jsonData, output); err != nil {
		return err
	}
	return nil
}
```

#### SendSafeJsonDTOMessage
Sends only safe messages in JSON format.
```golang
func SendSafeJsonDTOMessage(w http.ResponseWriter, code int, dto *rest.DTO, message typeopr.IPtr) error {
	if err := DeepCheckDTOSafeMessage(dto, message); err != nil {
		return err
	}
	if err := router.SendJson(message.Ptr(), w, code); err != nil {
		return err
	}
	return nil
}
```

#### FillDTOMessageFromMap
Fills in a message from the card.<br>
To work you need to use the `"dto"` tag.
If the DTO message is initially created correctly,
there should be no problem with this function.
```golang
func FillDTOMessageFromMap[T any](jsonMap map[string]interface{}, out *T) error {
	if jsonMap == nil || out == nil {
		return errors.New("nil input to FillMessageFromMap")
	}
	RV := typeopr.GetReflectValue(out)
	if !typeopr.IsImplementInterface(typeopr.Ptr{}.New(out), (*irest.Message)(nil)) {
		return errors.New("output param must implement the irest.IMessage interface")
	}
	rawObject := LoadSomeRawObjectFromCache(RV, &messageRawCache, namelib.TAGS.DTO)
	for name, f := range *rawObject.Fields() {
		inputValue, ok := (jsonMap)[name]
		if !ok {
			continue
		}
		fieldValue := RV.FieldByName(f.Name)
		switch f.Type.Kind() {
		case reflect.Struct:
			v, ok := inputValue.(map[string]interface{})
			if !ok {
				return fmt.Errorf("expected object for field '%s'", name)
			}
			if err := FillDTOMessageFromMap(v, &fieldValue); err != nil {
				return err
			}
		default:
			if err := fillField(&fieldValue, inputValue); err != nil {
				return err
			}
		}
	}
	return nil
}
```