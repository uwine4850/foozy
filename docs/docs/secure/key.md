## Key
`Key` structure that generates and stores three types of keys:

* hashKey — a key that is used for HMAC and can be dynamic.
* blockKey — a key that is used for encoding and can be dynamic.
* staticKey — a key that cannot change.
The old keys haskKey and blockKey are also stored here.

#### Key.GenerateBytesKeys
Generates keys.<br>
hashKey and blockKey will be updated. staticKey will only be generated once, cannot be regenerated.
```golang
func (k *Key) GenerateBytesKeys(length int) {
	k.oldHashKey = k.hashKey
	k.oldBlockKey = k.blockKey
	k.hashKey = string(k.generateKeys(length))
	k.blockKey = string(k.generateKeys(length))
	if k.staticKey == "" {
		k.staticKey = string(k.generateKeys(length))
	}
	k.date = time.Now()
}
```

#### Key.Generate32BytesKeys
Generates keys with a length of 32 bytes.
```golang
func (k *Key) Generate32BytesKeys() {
	k.GenerateBytesKeys(32)
}
```

#### Key.HashKey
Getting HashKey.
```golang
func (k *Key) HashKey() string {
	return k.hashKey
}
```

#### Key.OldHashKey
Getting the old OldHashKey.
```golang
func (k *Key) OldHashKey() string {
	return k.oldHashKey
}
```

#### Key.BlockKey
Getting BlockKey.
```golang
func (k *Key) BlockKey() string {
	return k.blockKey
}
```

#### Key.OldBlockKey
Getting the old OldBlockKey.
```golang
func (k *Key) OldBlockKey() string {
	return k.oldBlockKey
}
```

#### Key.StaticKey
Getting StaticKey.
```golang
func (k *Key) StaticKey() string {
	return k.staticKey
}
```

#### Key.Date
Getting the date of the last key update.
```golang
func (k *Key) Date() time.Time {
	return k.date
}
```