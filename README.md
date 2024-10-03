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
3. Enpoint обновления refresh токена реалтзован через считывание его из Cookie, т.к. это показалось более удобным.
4. Вместо ID у пользователей и больниц стоит UUID v4. Это помогло избежать лишнего запроса на вытягивание ID созданной сущности.
5. Все примеры запросов есть в файлах test_api.http. У каждого микросервиса есть такой файл. Там все запросы идут в том же порядке в каком идут в ТЗ.

### Буду рад любым комментариям и недочетам по выполненному заданию, спасибо
