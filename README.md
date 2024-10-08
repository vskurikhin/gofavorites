# Техническое задание

## Микро-сервис по работе с выбранными биржевыми инструментами GoFavorites

### Общие требования

GoFavorites представляет собой микро-сервис, позволяющую пользователю надёжно и безопасно хранить выбранные биржевые
инструменты и прочую приватную информацию.

#### Сервер должен реализовывать следующую бизнес-логику:

* обращение к внешнему микро-сервису для получения информации о зарегистрированных, аутентифицированных и
  авторизированных пользователях;
* хранение приватных данных;
* синхронизация данных между несколькими шардами (блоками) автоматизированной системы мобильных инвестиций для одного
  владельца;
* передача приватных данных владельцу по запросу.

#### Функции, реализация которых остаётся на усмотрение исполнителя:

* создание, редактирование и удаление данных на стороне сервера;
* формат работы с внешним микро-сервисом по регистрации аутентификация и авторизация пользователей;
* выбор хранилища и формат хранения данных;
* обеспечение безопасности передачи и хранения данных;
* протокол взаимодействия клиента и сервера;
* механизмы аутентификации пользователя и авторизации доступа к информации.

#### Дополнительные требования:

* сервер должен распространяться в виде docker-приложения с возможностью запуска на платформах Linux и Mac OS;
* сервер должен давать пользователю возможность получить информацию о версии и дате сборки бинарного файла клиента.

### Типы хранимой информации

* глобальный идентификатор пользователя (является персональными данными);
* произвольные текстовые данные;
* данные об биржевых инструментах.

Для любых данных должна быть возможность хранения произвольной текстовой метаинформации (принадлежность данных к
веб-сайту, личности или банку, списки одноразовых кодов активации и прочее).

### Абстрактная схема взаимодействия с системой

Ниже описаны базовые сценарии взаимодействия пользователя с системой. Они не являются исчерпывающими — решение отдельных
сценариев (например, разрешение конфликтов данных на сервере) остаётся на усмотрение исполнителя.

#### Для нового выбранного биржевого инструмента пользователя:

1) Пользователь проходит процедуру первичной регистрации и аутентификации на внешнем микро-сервисе.
2) Пользователь проходит процедуру авторизации на микро-сервисе GoFavorites.
3) Пользователь имеет возможность добавить новый биржевой инструмент на микро-сервисе GoFavorites.
4) Добавленный инструмент сохранятся как на текущем шарде (блоке) так и во внешнюю систему синхронизации между шардами (
   блоками).
5) Для обеспечения синхронизации информации о биржевых инструментах пользователя необходимо предусмотреть
   версионирование данных.

#### Для существующего биржевого инструмента пользователя:

1) Пользователь проходит процедуру первичной аутентификации на внешнем микро-сервисе.
2) Пользователь проходит процедуру авторизации на микро-сервисе GoFavorites.
3) При перемещении пользователя между шардами (блоками) необходимо синхронизировать биржевые инструменты пользователя по
   версии с помощью микро-сервиса GoFavorites и внешней системы синхронизации.
4) Пользователь имеет возможность получить актуальную информацию о своих биржевых инструментах на микро-сервисе
   GoFavorites.
5) Пользователь имеет возможность удалить биржевые инструменты.

### Тестирование и документация

Код всей системы должен быть покрыт юнит-тестами не менее чем на 80%. Каждая экспортированная функция, тип, переменная,
а также пакет системы должны содержать исчерпывающую документацию.

### Необязательные функции

Перечисленные ниже функции необязательны к имплементации, однако позволяют лучше оценить степень экспертизы исполнителя.
Исполнитель может реализовать любое количество из представленных ниже функций на свой выбор:

* поддержка данных типа OTP (one time password);
* использование бинарного протокола;
* наличие функциональных и/или интеграционных тестов;
* описание протокола взаимодействия клиента и сервера в формате Swagger.

[«GoFavorites»](https://github.com/vskurikhin/gofavorites)