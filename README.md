Создаем таблицу в pg.

```sql
create table redirects (id serial primary key, url character varying, p1 character varying, p2 character varying[], p3 character varying[], p4 character varying[], p5 character varying[], p6 character varying[]);
create index ps_idx on redirects (p1,p2,p3,p4,p5,p6);
grant all on redirects to mac;
```

Наполняем данными.

```sh
go run insert.go
```

Тест проводится на MBP 2012, выставляем настройки, чтобы хватило сокетов для обработки параллельных запросов.

```sh
sudo sysctl -w net.inet.ip.portrange.hifirst=32768
sudo sysctl -w net.inet.tcp.msl=1000
```

Проводим первый идеальный нагрузочный тест 30 секунд 10 параллельных запросов, который возвращает один и тот же ответ.
redirect_perfect.

```sh
siege -c 10 -f urls.txt -t 30S -b

Transactions:              54164 hits
Availability:             100.00 %
Elapsed time:              29.65 secs
Data transferred:           0.98 MB
Response time:              0.01 secs
Transaction rate:        1826.78 trans/sec
Throughput:             0.03 MB/sec
Concurrency:                9.58
Successful transactions:           0
Failed transactions:               0
Longest transaction:            0.50
Shortest transaction:           0.00
```

Добавляем github.com/valyala/fasthttp. Опять без чтения из базы, одинаковый ответ.

```sh
Transactions:              52558 hits
Availability:             100.00 %
Elapsed time:              29.99 secs
Data transferred:           1.31 MB
Response time:              0.01 secs
Transaction rate:        1752.52 trans/sec
Throughput:             0.04 MB/sec
Concurrency:                9.80
Successful transactions:       23736
Failed transactions:               0
Longest transaction:            0.43
Shortest transaction:           0.00
```


Читаем урлы из бд, без fasthttp. redirect_simple.

```sh
Transactions:              51102 hits
Availability:             100.00 %
Elapsed time:              29.89 secs
Data transferred:           1.28 MB
Response time:              0.01 secs
Transaction rate:        1709.67 trans/sec
Throughput:             0.04 MB/sec
Concurrency:                9.78
Successful transactions:       23085
Failed transactions:               0
Longest transaction:            0.43
Shortest transaction:           0.00
```

# Alloc = 4904600
# TotalAlloc = 287215784
# Sys = 13805816
# Lookups = 97012
# Mallocs = 3806390
# Frees = 3780593
# HeapAlloc = 4904600
# HeapSys = 8650752
# HeapIdle = 2080768
# HeapInuse = 6569984
# HeapReleased = 0
# HeapObjects = 25797
# Stack = 786432 / 786432
# MSpan = 103200 / 131072
# MCache = 9600 / 16384
# BuckHashSys = 1462954
# NextGC = 7694762


redirect_fasthttp с чтением из базы.

```sh
Transactions:              50785 hits
Availability:             100.00 %
Elapsed time:              30.00 secs
Data transferred:           0.00 MB
Response time:              0.01 secs
Transaction rate:        1692.83 trans/sec
Throughput:             0.00 MB/sec
Concurrency:                9.62
Successful transactions:       45875
Failed transactions:               0
Longest transaction:            0.45
Shortest transaction:           0.00
```

# runtime.MemStats
# Alloc = 5420616
# TotalAlloc = 84652880
# Sys = 13805816
# Lookups = 99936
# Mallocs = 1667624
# Frees = 1626690
# HeapAlloc = 5420616
# HeapSys = 8552448
# HeapIdle = 1671168
# HeapInuse = 6881280
# HeapReleased = 0
# HeapObjects = 40934
# Stack = 884736 / 884736
# MSpan = 102080 / 131072
# MCache = 9600 / 16384
# BuckHashSys = 1455098
# NextGC = 7407472

Сопоставимая производительность, но гораздо более эффективное потребление памяти. Можно сэкономить еще, отказавшись от ORM, можно положить в память всю таблицу с урлами, если есть существенный оверхед для запросов к бд, и получать обновления через механизм listen/notify.

Еще немного микрооптимизаций в фильнальной верии redirect_fasthttp.go

```sh
Transactions:              53606 hits
Availability:             100.00 %
Elapsed time:              29.95 secs
Data transferred:           0.00 MB
Response time:              0.06 secs
Transaction rate:        1789.85 trans/sec
Throughput:             0.00 MB/sec
Concurrency:               99.50
Successful transactions:       48276
Failed transactions:               0
Longest transaction:            0.66
Shortest transaction:           0.00
```
