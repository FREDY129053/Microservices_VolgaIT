# Microservices_VolgaIT
Задание полуфинального этапа дисциплины Backend-разработка Web API

## Основное задание
1. Account Doc URL: http://localhost:8081/api/docs/index.html
2. Hospital Doc URL: http://localhost:8082/api/Hospitals/docs/index.html
3. Timetable Doc URL: http://localhost:8083/api/docs/index.html
4. Document Doc URL: http://localhost:8084/api/History/docs/index.html


## Дополнительная инофрмация
1. Схема базы данных доступна в репозитории, полная схема доступна по ссылке: https://dbdesigner.page.link/aVCiLaDshYx5QuPR9.
2. Endpoint верификации токена не использовался, т.к. по сути(по коду) он проверяется в middleware и смысла делать запрос нет.
3. Токены оба хранятся в Cookie.
4. Enpoint обновления refresh токена реализован через считывание его из Cookie, т.к. это показалось более удобным.
5. Вместо ID у пользователей и больниц стоит UUID v4. Это помогло избежать лишнего запроса на вытягивание ID созданной сущности.
6. Все примеры запросов есть в файлах test_api.http. У каждого микросервиса есть такой файл. Там все запросы идут в том же порядке в каком идут в ТЗ.
7. Все пользователи по ТЗ доступны в базе данных с соответствующими ролями. Так же у всех заданы имя и фамилия.

### Если сборка не получается, то в папке есть дамп базы данных. Создать БД и запустить каждый микросервис через go run main.go. Тогда они будут доступны на http://0.0.0.0:808x/...

### Буду рад любым комментариям и недочетам по выполненному заданию, спасибо
