# dbKeyValue

Реализация LSM tree хранилища key-value.

Чтобы запустить, необходимо прописать переменную окружения:

```bash
CONFIG_PATH=./config/config.yaml
```

И запустить `main()`. Либо запустить `main()` с флагом `--config`.

Описание основных сущностей ниже.

## LSMT

Имеет интерфейс `Get(key string)`, `Put(key string, value []byte)`, `Delete(key string)`.

Содержит текущую `MemTable`; список `MemTable`, подлежащих записи на диск; список `SSTable`.

Следит за размером `MemTable`, когда он превышает пороговый записывает в списокподлежащих записи на диск.

Раз в единицу времени запускает процесс записи всех `MemTable` списка на диск.

Доступ к данным `Get(key string)` происходит так в последовательности: 

- Проверить, есть ли ключ в текущей `MemTable`;
- В списке `MemTable`, подлежащих выгрузке на диск;
- В списке `SSTable`

## MemTable

Структура в оперативной памяти, которая поддерживает методы: 
`Get(key string)`, `Put(key string, value []byte)`, `Delete(key string)`.

Хранение данных и доступ к ним обеспечивается при помощи `Bucket`.

При достижении максимально допустимого по конфигурации пользователя размера хранящихся в ней данных
выгружает методом `Write` содержимое `Bucket` на диск.

### Bucket

Содержит структуру данных `btree`, размер которой задается пользователем в конфигурации.

Удаляет/добавляет данные по переданному ключу в `btree`. Всегда хранит только последнюю версию ключа.

Поддерживает метод `Scan()` - возвращает все хранимые структурой данные в виде массива структур `Node`.
Массив упорядочен лексикографически по ключу.

#### Node

Элемент `btree`.

## SSTable

Структура в оперативной памяти, которая содержит:

- `filepath` - путь до директории с выгруженными на диск файлами с данными `*.sstable`;

- `index` - `red-black-tree` структура данных, разреженный индекс построенный по ключам данных;

- `bloomFilter` - фильтр Блума, структура данных, позволяющая за `O(1)` ответить, содержит ли `SSTable` данные по запрашиваемому ключу.

Поддерживает операцию `Get(key string)`. 

Поддерживает `Delete(key string)` - выполняет перестроение фильтра Блума, но не удаляет ключ из соответствующего файла на диске.

### Entry

Единица данных `SSTable` в оперативной памяти.

Поддерживает методы `Marshall()`, `Unmarshall(in []byte)`.
Преобразует значения `key` и `value` в массив байт `data` вида: `keyLength`(8 байт), `valueLength`(8 байт), 
`bKey_payload`(размера `keyLength`), `bValue_payload`(размера `valueLength`). 

## Flush

Выполняет `Flush()` данных соответствующей `MemTable` на диск.

## TODO

- В данный момент база данных запускается с чистого листа - не подгружает ранее выгруженные `SSTable` с диска, не восстанавливает данные `MemTable` и списка `MemTable` после отключения;
- Не восстановления данных, находящихся в структурах, лежащих в оперативной памяти;
- Нет процесса `Compaction` для файлов `SSTable` на диске;