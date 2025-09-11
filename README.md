# CP ultimate money control system
*Система предназначения для отмыва денег*

## Архитектура
![Арх](/docs/Архитектура.svg)

## Сообщения grpc
### Регистрация воркера
- Отправляем
```
auth: str = "ask"
```
- Получаем
```
uuid: str
key: str
```

### Получить здание
- Отправляем
```
worker_uuid: str
worker_key: str
```
- Ответ
```
uuid: str
task: str
payload: str(json)
unix_time: uint64
```

### Отметить выполненным
- Отправляем
```
uuid: str
payload: str
unix_time: uin64
```
- Получаем
```
ok: bool = 1
```