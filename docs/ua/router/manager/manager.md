## Package manager
Даний пакет реалізує алгоритм роботи менеджера. Менеджер потрібен для деяких
налаштувань проекту та для передачі даних від серевера у обробник http запиту.
На момент написання менеджер реалізує чотири інтерфейси які виконують наступні 
звадання:<br>
* interfaces.IManagerWebsocket — реалізує ініціалізацію та базову взаємодію із 
вебсокетами.
* interfaces.IManagerConfig — зберігає в собі деякі налаштування проекту.
* interfaces.IManagerOneTimeData — зберігає одноразові данні для кожного запиту.
* interfaces.IRender — використовується для відображення html сторінки.

Важливо зазначити, що IManagerOneTimeData та IRender відрізняються від інших,
 адже вони унікальні для кожного запиту.

Реалізації всіх цих інтерфейсів знаходяться у структурі ``Manager``, яка виглядає
 наступним чином:
```
type Manager struct {
    managerWebsocket interfaces.IManagerWebsocket
    managerConf      interfaces.IManagerConfig
    managerData      interfaces.IManagerOneTimeData
    render           interfaces.IRender
}
```
У даній структурі занаходяться прості "get" та "set" методи для вствновлення та 
отримання окремих менеджерів.

Детальніше про кожен з менеджерів можна прочитати по посиланням:

* [IManagerWebsocket](https://github.com/uwine4850/foozy/blob/master/docs/ua/router/manager/manager_ws.md)
* [IManagerConfig](https://github.com/uwine4850/foozy/blob/master/docs/ua/router/manager/manager_conf.md)
* [IManagerOneTimeData](https://github.com/uwine4850/foozy/blob/master/docs/ua/router/manager/manager_otd.md)
* [IRender](https://github.com/uwine4850/foozy/blob/master/docs/ua/router/tmlengine/page_render.md)