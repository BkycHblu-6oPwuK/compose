## Создание корневого сертификата (CA)

Операции по созданию, актуально только, если вы хотите самостоятельно создать, сейчас же все создано и находится в директории ./certs

1. Создайте частный ключ для корневого центра сертификации (CA):

```bash
openssl genrsa -out rootCA.key 2048
```

2. Создайте корневой сертификат

```bash
openssl req -x509 -new -nodes -key rootCA.key -sha256 -days 9999 -out rootCA.crt
```

## Создание SSL-сертификата для localhost

1. Создайте частный ключ для localhost:
```bash
openssl genrsa -out localhost.key 2048
```
2. Создайте запрос на сертификат (CSR) для localhost:
```bash
openssl req -new -key localhost.key -out localhost.csr
```

3. Создайте файл конфигурации для расширений сертификата (SAN - Subject Alternative Name)

```
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
DNS.2 = app # название контейнера
```

4. Подпишите сертификат с помощью CA

```bash
openssl x509 -req -in localhost.csr -CA rootCA.crt -CAkey rootCA.key -CAcreateserial -out localhost.crt -days 9999 -sha256 -extfile localhost.ext
```

## Импорт корневого сертификата в доверенное хранилище на хосте linux

В контейнерах эта операция уже выполняется, команда актуальна только, если вы в целом работаете из под linux

```
sudo cp ./rootCA.crt /usr/local/share/ca-certificates/
sudo update-ca-certificates
```

## Импорт в Windows

1. скачайте rootCA.crt и скопируйте в подходящее для вас место
2. Откройте меню "Пуск" и введите certmgr.msc, затем нажмите Enter. Это откроет управление сертификатами.
3. В левой панели найдите папку "Доверенные корневые центры сертификации". Раскройте ее, чтобы увидеть подкатегории.
4. Щелкните правой кнопкой мыши на "Сертификаты" в "Доверенные корневые центры сертификации" и выберите "Все задачи" → "Импорт".
5. В мастере импорта сертификатов:
    1.  Нажмите "Далее".
    2.  Нажмите "Обзор" и выберите файл rootCA.crt.
    3.  Убедитесь, что выбран пункт "Местный компьютер", если он не выбран по умолчанию, затем нажмите "Далее".
    4.  Выберите "Поместить все сертификаты в следующий хранилище" и выберите "Доверенные корневые центры сертификации".
    5.  Нажмите "Далее", затем "Готово".
6. Проверьте, что сертификат успешно установлен. Перезапустите браузеры, теперь при заходе на сайт не должно быть сообщения об недоверии.