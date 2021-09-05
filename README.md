
#### Реализованная функциональность

-	Скраппинг интернета
-	Генерация писем на основе шаблона
-	Парсинг реестра gov.ru, проверка добросовестности поставщика
-	Веб морда

#### Особенность проекта в следующем:

-   Автоматическая, конфигурируемая по шаблонам рассылка идеально подходящим поставщикам
-   Поточный поиск. Поддержка как горячих запросов, так и холодного мониторинга
-   Генерация сложных запросов с ОКПД2, оцентой добросовестности и прочим на базе пары слов
-   Оценка цены, опыта, качества поставщика на основе его истории
-   Безопасность


#### Основной стек технологий:
- Фронтенд:
  -   HTML, CSS, TypeScript.
  -   LESS, SCSS, PostCSS.
  -   Webpack
  -   React
  -   Redux
  - Redux-thunk
  - Axios
- Бэкенд
  - Golang
- Общее
-   	Git

#### Демо

Демо сервиса доступно по адресу: [https://unicorn-leaders.vercel.app](https://unicorn-leaders.vercel.app)

## СРЕДА ЗАПУСКА

1.  Развертывание сервиса производится на Unix-based системах (фронтенд запускался на MacOS);
2.	Для демо пакета под дистрибутивы не собрано, нужно собирать бинарь, любой Linux дистрибутив, на который встает go, либо MacOS.

## УСТАНОВКА

### Установка фронтенда

Для установки нужен пакет `yarn`.

Выполните

```
git clone https://github.com/un1c0rn-dev/digital_breakthrough_2021
cd digital_breakthrough_2021/frontend
yarn install

```

### Установка бэкенда

cd ```project_dir``` && go build .

### Запуск фронтенда

Выполните

```
cd digital_breakthrough_2021/frontend
yarn start
```

### Запуск бэкенда

```
(cd <project_dir> && echo "[Unit]
Description=Unicorn dev webscrapper
After=network.target

[Service]
ExecStart=$(pwd)/unicorn.dev.web-scrap -api-keys $(pwd)/Configs/api_keys.json -port 443 -use-tls -tls-key $(pwd)/Configs/server.key -tls-crt $(pwd)/Configs/server.crt
ExecReload=$(pwd)/unicorn.dev.web-scrap -api-keys $(pwd)/Configs/api_keys.json -port 443 -use-tls -tls-key $(pwd)/Configs/server.key -tls-crt $(pwd)/Configs/server.crt
KillMode=process
Restart=on-failure

[Install]
WantedBy=multi-user.target
" > /etc/systemd/system/unicorn-scap.service)

systemctl enable unicorn-scap.service 
systemctl start unicorn-scap.service 
```
Приветствуется запуск в докере

### Установка зависимостей проекта

sudo apt install golang

## РАЗРАБОТЧИКИ

#### Тимур Черных backend [https://t.me/j35uScHr1St](https://t.me/j35uScHr1St)
#### Стефан Тюрин backend [https://t.me/Rumb0](https://t.me/Rumb0)
#### Дмитрий Никулин frontend [https://t.me/LaExplorad0ra](https://t.me/LaExplorad0ra)
