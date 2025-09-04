## database

### DatabasePool
stores a pool of database connections.

To work properly, you need to connect to the database using the `Open` method 
and pass the open connection to the `ConnectionPool` method. There can be several connections, 
but usually one connection pool is enough.

After all the settings you need to call the `Lock` method to lock all the changes. After calling 
this method you can no longer add connections to the pool.
Thus, the connection pool remains static and predictable from the beginning to the end of the application.
```golang
type DatabasePool struct {
	// Type map[string]interfaces.IReadDatabase.
	connnectionPool sync.Map
	locked          atomic.Bool
}
```

#### DatabasePool.ConnectionPool
Get a named connection pool.
```golang
func (dp *DatabasePool) ConnectionPool(name string) (interfaces.DatabaseInteraction, error) {
	if !dp.locked.Load() {
		return nil, &ErrDatabasePoolNotLocked{}
	}
	if val, ok := dp.connnectionPool.Load(name); ok {
		return val.(interfaces.DatabaseInteraction), nil
	} else {
		return nil, &ErrConnectionNotExists{Name: name}
	}
}
```

#### DatabasePool.AddConnection
Adding a new open named database connection pool.
```golang
func (dp *DatabasePool) AddConnection(name string, dbInteraction interfaces.DatabaseInteraction) error {
	if dp.locked.Load() {
		return &ErrDatabasePoolIsLocked{}
	}
	if _, ok := dp.connnectionPool.Load(name); ok {
		return &ErrConnectionAlreadyExists{Name: name}
	} else {
		dp.connnectionPool.Store(name, dbInteraction)
		return nil
	}
}
```

#### DatabasePool.Lock
Blocks further changes to the pool.
```golang
func (dp *DatabasePool) Lock() {
	dp.locked.Store(true)
}
```