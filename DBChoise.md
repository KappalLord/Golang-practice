# MongoDB: 
NoSQL БД, использует формат хранения данных BSON (Binary JSON). 
MongoDB хранит данные в виде документов, категории походят на папки с документами.

## Плюсы:
1. **Гибкость:** если нужна гибкая схема(не имеем четкой структуры данных) и высокая производительность - монго очень хорошо справляется.
2. **Масштабируемость:** Поддерживает горизонтальное масштабирование, можно спавнить новые машины и тем самым увеличивать мощность.

3. **Обработка JSON:** монго использует BSON,что по факту является представление json.Следовательно, хорошо подходит для работы веб сервиса.

4. **Уровень согласованности:** можно выбрать уровень согласованности в зависимости от важности данных. Если не так критична сохранность данных, можно добиться куда более высокой производительности.


## Минусы:
1. **Транзакции:** Сильно хуже справляется с транзакциями, по сравнению с реляционными БД.
2. **Меньше свободы в запросах:** куда меньше гибкости в составлении запросов, например, join не применить.
3. **Размер данных:** Данные занимают побольше места, так как у документов хранятся и имена полей.

## Сложность операций*:
- Чтение(SELECT): O(log n)
- Запись(INSERT): O(n)  
- Изменение(UPDATE): O(log n)
- Удаление(DELETE): O(log n)

# PostgreSQL:
PostgreSQL —  созданная коммьюнити объектно-реляционная СУБД с открытым исходным кодом. Поддерживает и реляционные(дефолт SQL), и нереляционные(JSON) запросы.

# Плюсы:
1. **Все возможности SQL:** Поддержка различных сложных запросов и тому подобного.
2. **Транзакции:** Сильно упрощает написание кода, ибо не надо по к.д. писать обработку ошибок и т.д.
3. **Надёжность:** Засчет тех же самых транзакций обеспечивается и надёжность, ведь транзакции соответствуют требованиям ACID:
   - **Атомарность (Atomicity)**:
      - Транзакция выполняется целиком или не выполняется вовсе. Все изменения данных происходят, только если транзакция завершается успешно.

   - **Согласованность (Consistency)**:
      - Транзакция переводит базу данных из одного согласованного состояния в другое. Все бизнес-правила и ограничения должны соблюдаться до и после транзакции.

   - **Изолированность (Isolation)**:
      - Результаты выполнения транзакции не видны другим транзакциям до её завершения. Это предотвращает взаимное влияние параллельно выполняющихся транзакций.

   - **Долговечность (Durability)**:
      - После завершения транзакции её результаты сохраняются в базе данных даже в случае сбоев системы.

Все работает благодаря блокировки на чтение - несколько транзакций могут параллельно считывать данные, но не записывать. блокировка на запись - одна транзакция записывает данные, остальные ждут, пока закончит запись.

## Минусы:
1. **Производительность при больших объемах данных:** Начинает довольно медленно работать, когда необходимо работать с уже действительно большим объемом данных.
2. **Зависимость от ресурсов:** Требует больше ресурсов по сравнению с некоторыми легковесными СУБД. Можно масштабировать только вертикально, то есть апгрейдить железки компьютера, что не очень удобно.

## Сложность операций*:
- Чтение (SELECT): O(log n) 
- Запись (INSERT): O(n)
- Изменение(UPDATE): O(log n)
- Удаление(DELETE): O(log n)

# ClickHouse:
ClickHouse — это колоночная СУБД с открытым исходным кодом, основная фишка - хранит данные в колонках, а не в строках. 
## Плюсы:
1. **Высокая производительность с большими данными:** засчет своих особенностей хороша в анализ большого объема данных.
2. **Масштабируемость:** Легко масштабируется как вертикально, так и горизонтально, оч удобно.
3. **Сжатие данных:** Хороший алгоритм сжатия, освобождает много пространства.

## Минусы:
1. **Не подходит для транзакций:** ClickHouse не предназначен для работы с транзакциями, ибо создан для других целей.
2. **Сложность точечной работы:** Исходя из особенностей бд, взаимодействие с отдельными строками создаёт определенные трудности.
3. **Зависимость от оперативки:** скорость работы сильно зависит от объема и быстроты оперативки.
4. **Баги и уязвимости:** довольно часто возникают ошибки, а также существуют весомые уязвимости, так что безопасной данную БД явно не назовешь.

## Сложность операций*:
- Чтение (SELECT): O(log n) 
- Запись (INSERT): O(n)
- Изменение(UPDATE): O(n)
- Удаление(DELETE): O(n)


# Вердикт

**MongoDB** - довольно удобная база данных для использования в определенных сценариях, где необходимо обеспечить высокую производительность для обработки данных в риал тайме. Скорее как вспомогательный инструмент для таких случаев, в остальном, имеет смысл использовать классическую реляционную БД. 

**PostgreSQL** - хорошая БД для хранения всей основной необходимой информации. Очень надёжная, можно делать сложные SQL запросы, но и при этом и довольно требовательна. В случае, если необходимо обработать много информации, функционирует довольно медленно. 

**ClickHouse** - БД, созданная для определенных целей, таких как анализ большого массива данных. Если мы говорим о динамически обновляемых базах данных, то кликхаус для работы не подойдёт.


# Какие бд для чего используем

1. **Хранить данные юзеров, балансы монеток и т.д.**

Для данных целей я бы использовал PostgreSQL, ибо у данных явно будет существовать определенная структура, и сохранность конкретно этих данных для нас критична. По идее тут не должно быть сильно много запросов на минуту времени, следовательно недостатки не так критичны в данном случае.

2. **Хранить гео и мак адреса, которые насобирал юзеры**

Для этого я бы использовал MongoDB. Во первых, монго хорошо работает с геодатой, что нам и нужно. Также придется обрабатывать большой поток данных в риал тайме, и сохранность этих данных не столь критично, то есть сможем выдать большую производительность и комфортно обработать данные. Также горизонтальное масштабирование позволит довольно просто сохранить скорость обработки данных в случае роста объема данных.

