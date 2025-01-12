## Описание проекта
Планировщик задач с функциями добавления, редактирования, удаления задач, а также поиска по задачам и возможностью аутентификации

## Список выполненных задач со звёздочкой
Были выполнены все дополнительные задачи, а именно:
    * Возможность задавать порт веб-сервера, путь к файлу БД и пароль через переменные окружения
    * Сложные правила повторения
        * Еженедельное повторение задачи по определенным дням недели
        * Ежемесячное повторение задачи по определенным числам месяца
    * Возможность поиска задач по заголовку и комментарию, по дате
    * Возможность аутентификации при указанном пароле

## Инструкция по запуску кода
### Запуск проекта
    Для запуска проекта локально, следует либо ввести в терминал
    в директории проекта команду <div>go run .</div>,
    либо же открыть файл go_final_project.exe
### Определение переменных окружения
    В проекте можно задать три переменных окружения:
        * PORT - порт веб-сервера
        * TODO_DBFILE - путь к файлу базы данных
        * TODO_PASSWORD - пароль для дальнейшей аутентификации
    их можно задать в файле .env в корне проекта
    Пример такого файла:
    <div>
    TODO_PASSWORD = "1234"
    PORT = "8080"
    TODO_DBFILE = "../scheduler.db"
    </div>
### Проект в браузере
    Для того чтобы после локального запуска получить доступ к сервису,
    в браузере следует указать адрес localhost:<указанный порт>
    По умолчанию значение порта равняется 7540

## Запуск тестов
    Тестирование всего приложения следует запускать с помощью команды <div>go test ./tests</div>
    При этом в файле tests/settings.go можно изменить некоторые параметры:
        * Port - для изменения порта по умолчанию
        * DBFile - изменение адреса файла базы данных
        * FullNextDate - true для тестов всех условий повторения задач, false - для ограниченного
        * Search - true для тестов поиска, false - без поиска
        * Token - JWT-токен для аутентификации

## Сборка проекта через Docker