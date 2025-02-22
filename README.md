### sap segmentation

Поднять локальное окружение:

```
make docker-compose up -d
```

Запустить, подключившись к локальному окружению:

```
make run-local
```

Прогнать тесты:

```
make test
```

Прогнать интеграционные тесты (нужна база из локального окружения):

```
make test-all
```