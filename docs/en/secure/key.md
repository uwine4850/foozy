## Key
The object is responsible for all key operations. For example: generation, storage, provisioning of keys.

__HashKey__
```
HashKey() string
```
Generates a hash key.

__OldHashKey__
```
OldHashKey() string
```
Returns the old login key.
This method is used after the hash key is generated again.
It is important to note that this method returns the previous key, starting with the active one.

__BlockKey__
```
BlockKey() string
```
Generates a block key.

__OldBlockKey__
```
OldBlockKey() string
```
Returns the old login key.
This method is used after the block key is generated again.
It is important to note that this method returns the previous key, starting with the active one.

__StaticKey()__
```
StaticKey() string
```
Creates a static key. This key is generated once when the server is started and does not change.

__Date() time.Time__
```
Date() time.Time
```
Returns the time of the last key generation.

__GenerateBytesKeys__
```
GenerateBytesKeys(length int)
```
Generates block key and hash key

__Generate32BytesKeys__
```
Generate32BytesKeys()
```
Generates 32-byte keys.

__Get32BytesKey__
```
Get32BytesKey() IKey
```
Returns an IKey structure responsible for key operations.