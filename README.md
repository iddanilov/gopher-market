# go-musthave-diploma-tpl

Схему с общей архитектурой решения:
https://app.diagrams.net/#Higitddanilov%2Fgopher-market%2Finit%2F%D0%94%D0%B8%D0%B0%D0%B3%D1%80%D0%B0%D0%BC%D0%BC%D0%B0%20%D0%B1%D0%B5%D0%B7%20%D0%BD%D0%B0%D0%B7%D0%B2%D0%B0%D0%BD%D0%B8%D1%8F.drawio

Шаблон репозитория для индивидуального дипломного проекта курса «Go-разработчик»


# Как запустить проект
1. Нужно поднять инфру

База данных: 
   >    make docker-compose-up
   
2. Выполнить миграцию

Установка инструмента
   > make install-goose

Запустить миграцию
> make test-migrations-up

3. Установить все модули 

> make go-vendor

4. Запустить

> make run


# Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```
git remote add -m master template https://github.com/yandex-praktikum/go-musthave-diploma-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/master .github
```

Затем добавьте полученные изменения в свой репозиторий.
