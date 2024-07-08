# Начало работы

## Инициализация файлов

Помимо структур и функций, облегчающих разработку, можно быстро создать следующие файлы:

- .gitignore
- .golangci.yaml - конфигурация для [golangci-lint](https://golangci-lint.run/)
- autocomplete.sh - добавляет автодополнение для [./helper.sh](./helper.sh). Нужно добавить в .bashrc
- helper.sh - набор функций для более удобной работы с утилитами из командной строки
- config - вариант хранения конфигурации проекта
  - aws - конфигурация aws
  - .gitignore
  - main.yaml - основной файл конфигурации
  - main-local.yaml - файл конфигурации с конфиденциальными данными
  - main-local.yaml.ex - пример файла с конфиденциальными данными

Для создания нужно добавить компонент `environment.Initializer` в компонент приложения
(конструктор - `environment.NewInitializer`).

Далее в коде вызвать:
```go
err := app.Initializer.Initialize(environment.InitializerConfig{
    CreateAutocompleteShell:   true,
    CreateGitignore:           true,
    CreateGolangCIConfig:      true,
    CreateHelperShell:         true,
    CreateReadmeMd:            true,
    CreateDefaultConfigFolder: true,
})
```

# Утилиты для разработки

## golangci-lint

https://golangci-lint.run/

Данная утилита позволяет проверят код множеством линтеров.

Для конфигурации используется файл `.golangci.yaml`.

## gofumpt

https://github.com/mvdan/gofumpt

Более строгая версия gofmt.
Желательно настроить вызов данной утилиты при сохранении файлов `.go`.

## helper.sh

`./helper.sh` - bash-скрипт, который облегчает вызов консольных утилит.

Для получения информации о командах:
```shell
./helper.sh --help
```

### Автодополнение

Чтобы включить автодополнение, нужно добавить строчку `source /path/to/autocomplete.sh"` в файл `.bashrc`