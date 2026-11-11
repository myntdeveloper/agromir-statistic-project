# Приём показаний ёмкостей (S7-1200 / WinCC → API)

Контроллер Siemens S7-1200 (или WinCC Runtime на ПК завода) раз в N минут отправляет HTTP POST с JSON на API. Отдельный агент на ПК завода не нужен.

```
S7-1200 / WinCC  ──HTTP POST──►  nginx :80  ──►  api :8000  ──►  PostgreSQL
```

## Куда отправлять

| Параметр | Значение |
|---|---|
| Метод | `POST` |
| URL | `http://168.113.176.218/api/v1/ingest/readings` |
| IP | `168.113.176.218` |
| Порт | `80` (обычный HTTP, без TLS) |
| Path | `/api/v1/ingest/readings` |

По IP nginx отдаёт сайт и API по обычному HTTP, без редиректа и без сертификатов. Это основной вариант для ПЛК: DNS, NTP и импорт сертификатов не нужны.

Запасной вариант — `https://agromir-stat.online/api/v1/ingest/readings` (порт `443`). Учтите:

- `http://agromir-stat.online/...` (по домену, порт 80) отвечает редиректом `301` на HTTPS, а ПЛК и WinCC редиректы для POST не обрабатывают. По HTTP обращайтесь только по IP.
- `https://168.113.176.218/...` не пройдёт проверку сертификата: он выдан на домен, а не на IP.
- Для HTTPS с S7-1200 нужно импортировать корневой сертификат Let's Encrypt **ISRG Root X1** (TIA Portal → Global security settings → Certificate manager) и задать DNS в ПЛК.

По HTTP ключ `AGENT_API_KEY` передаётся открытым текстом. Если у завода статический внешний IP, ограничьте приём только им (см. «Ограничение по IP завода»).

## Заголовки

```http
POST /api/v1/ingest/readings HTTP/1.1
Host: 168.113.176.218
Authorization: Bearer <AGENT_API_KEY>
Content-Type: application/json
Content-Length: <длина тела в байтах>
Connection: close
```

`<AGENT_API_KEY>` — значение `AGENT_API_KEY` из `services/api/.env` на сервере. Без заголовка или с неверным ключом сервер ответит `401`.

Если отправитель не умеет задавать `Authorization`, ключ можно передать в любом из заголовков (проверяются по порядку, префикс `Bearer ` необязателен):

| Заголовок | Пример |
|---|---|
| `Authorization` | `Authorization: Bearer <AGENT_API_KEY>` |
| `Header` | `Header: <AGENT_API_KEY>` или `Header: Bearer <AGENT_API_KEY>` |
| `X-API-Key` | `X-API-Key: <AGENT_API_KEY>` |

Все отклонённые запросы пишутся в лог API (`docker compose logs -f api`):

```
ingest auth: no key from 95.46.67.28 (expected one of headers [Authorization Header X-API-Key])
ingest auth: wrong key from 95.46.67.28 in header "Header" (AgRm***, len=33)
ingest from 95.46.67.28: invalid body: <причина>; body: <первые 1000 байт тела>
ingest from 95.46.67.28 (plant plant-01): unknown boiler_number: [25]
```

## Тело запроса (JSON)

### Минимальный пакет (рекомендуется для ПЛК)

Если поле `ts` не передать, сервер подставит `sent_at`, а если нет и его, то своё текущее время UTC. При отправке «здесь и сейчас» так проще всего: часы ПЛК синхронизировать не нужно.

```json
{
  "plant_id": "plant-01",
  "readings": [
    { "boiler_number": 1, "level_mm": 7301, "density": 1.35, "volume_m3": 159.98, "volume_l": 159980, "mass_t": 215.97, "income_month": 12.5, "outflow_month": 8.1, "income_year": 140.2, "outflow_year": 96.4 },
    { "boiler_number": 2, "level_mm": 7313, "density": 1.35, "volume_m3": 159.61, "volume_l": 159608, "mass_t": 215.47, "income_month": 0, "outflow_month": 0, "income_year": 0, "outflow_year": 0 }
  ]
}
```

### Полный пакет (со своим временем)

```json
{
  "plant_id": "plant-01",
  "sent_at": "2026-11-11T17:30:00Z",
  "readings": [
    {
      "boiler_number": 1,
      "ts": "2026-11-11T22:30:00+05:00",
      "level_mm": 7301,
      "density": 1.35,
      "volume_m3": 159.98,
      "volume_l": 159980.0,
      "mass_t": 215.97,
      "income_month": 12.5,
      "outflow_month": 8.1,
      "income_year": 140.2,
      "outflow_year": 96.4
    }
  ]
}
```

### Поля

| Поле | Тип | Обязательно | Описание |
|---|---|---|---|
| `plant_id` | string | да | Идентификатор завода, например `plant-01` |
| `sent_at` | string (RFC 3339) | нет | Время отправки пакета |
| `readings` | array | да, минимум 1 | Показания ёмкостей |
| `readings[].boiler_number` | int | да | Номер ёмкости, `1`–`24` |
| `readings[].ts` | string (RFC 3339) | нет | Время замера. Если не передать — `sent_at` / время сервера |
| `readings[].level_mm` | number | нет* | Уровень, мм |
| `readings[].density` | number | нет* | Плотность, т/м³ |
| `readings[].volume_m3` | number | нет* | Объём, м³ |
| `readings[].volume_l` | number | нет* | Объём, л |
| `readings[].mass_t` | number | нет* | Масса, т |
| `readings[].income_month` | number | нет* | Приход за месяц |
| `readings[].outflow_month` | number | нет* | Расход за месяц |
| `readings[].income_year` | number | нет* | Приход за год |
| `readings[].outflow_year` | number | нет* | Расход за год |

\* Пропущенное числовое поле сервер считает равным `0` и **перезапишет** им текущее значение. Поэтому отправляйте все поля в каждом пакете.

### Правила формата (частые ошибки ПЛК)

- **Дробный разделитель — только точка**: `1.35`, а не `1,35`. В WinCC при русской локали `CStr()` даёт запятую, её нужно заменить.
- Числа — без кавычек, без ведущего `+` и пробелов. Допустимо `7301`, `7301.0`, `1.35E+00`. Недопустимо `"7301"`, `+7301`, ` 7301`, `NaN`, `Inf`.
- Время — строго RFC 3339 с часовым поясом: `2026-11-11T17:30:00Z` или `2026-11-11T22:30:00+05:00`. Формат `2026-11-11 22:30:00` без `T` и без пояса сервер не примет. Пустую строку `""` тоже не передавать: либо корректное время, либо не передавать поле вообще.
- Кодировка — UTF-8. Лишние поля сервер игнорирует.
- Неизвестный `boiler_number` (не из `1`–`24`) отклоняет **весь** пакет.
- Можно отправлять все 24 ёмкости одним пакетом (около 6 КБ) или несколькими пакетами поменьше, если буфер ПЛК ограничен.

### Как часто отправлять

Каждый пакет обновляет текущие показания и добавляет точку на график (одна точка на ёмкость на каждое уникальное `ts`). Рекомендуемый интервал — 1–15 минут. Чаще раза в минуту отправлять не стоит: таблица графика быстро разрастается.

## Ответы сервера

| Код | Тело | Что значит |
|---|---|---|
| `200` | `{"accepted": 24}` | Пакет принят и сохранён |
| `400` | `{"error":"invalid request body","details":"..."}` | Битый JSON, неверный тип поля или формат времени. В `details` — причина |
| `400` | `{"error":"unknown boiler_number","numbers":[25]}` | В пакете есть ёмкость, которой нет в базе. Пакет не сохранён |
| `401` | `{"error":"unauthorized"}` | Нет заголовка `Authorization` или ключ неверный |
| `500` | `{"error":"internal server error"}` | Ошибка на сервере (БД). Повторить позже |

Если ответ не `200`, пакет не сохранён. Его можно отправить повторно: дублей на графике не будет, если `ts` тот же.

## Вариант 1: WinCC (VBScript, PC Runtime)

Работает в WinCC Classic и WinCC (TIA) RT Advanced на ПК. На панелях Comfort/Basic `CreateObject` недоступен. Скрипт вешается на циклический триггер, например раз в 5 минут. Имена тегов (`Boiler1_Level` и т. д.) замените на свои.

```vbscript
Function Num(v)
    Num = Replace(CStr(CDbl(v)), ",", ".")
End Function

Sub SendBoilerReadings()
    Const API_URL = "http://168.113.176.218/api/v1/ingest/readings"
    Const API_KEY = "<AGENT_API_KEY>"

    Dim i, json, sep, p
    json = "{""plant_id"":""plant-01"",""readings"":["
    sep = ""
    For i = 1 To 24
        p = "Boiler" & i & "_"
        json = json & sep & "{" & _
            """boiler_number"":" & i & "," & _
            """level_mm"":"      & Num(HMIRuntime.Tags(p & "Level").Read)        & "," & _
            """density"":"       & Num(HMIRuntime.Tags(p & "Density").Read)      & "," & _
            """volume_m3"":"     & Num(HMIRuntime.Tags(p & "VolumeM3").Read)     & "," & _
            """volume_l"":"      & Num(HMIRuntime.Tags(p & "VolumeL").Read)      & "," & _
            """mass_t"":"        & Num(HMIRuntime.Tags(p & "MassT").Read)        & "," & _
            """income_month"":"  & Num(HMIRuntime.Tags(p & "IncomeMonth").Read)  & "," & _
            """outflow_month"":" & Num(HMIRuntime.Tags(p & "OutflowMonth").Read) & "," & _
            """income_year"":"   & Num(HMIRuntime.Tags(p & "IncomeYear").Read)   & "," & _
            """outflow_year"":"  & Num(HMIRuntime.Tags(p & "OutflowYear").Read)  & "}"
        sep = ","
    Next
    json = json & "]}"

    Dim http
    Set http = CreateObject("MSXML2.ServerXMLHTTP.6.0")
    On Error Resume Next
    http.setTimeouts 5000, 5000, 10000, 10000
    http.open "POST", API_URL, False
    http.setRequestHeader "Content-Type", "application/json"
    http.setRequestHeader "Authorization", "Bearer " & API_KEY
    http.send json
    If Err.Number <> 0 Then
        HMIRuntime.Trace "boiler send failed: " & Err.Description & vbCrLf
    ElseIf http.status <> 200 Then
        HMIRuntime.Trace "boiler send status " & http.status & ": " & http.responseText & vbCrLf
    End If
    On Error GoTo 0
    Set http = Nothing
End Sub
```

На ПК с WinCC должен быть выход в интернет на `168.113.176.218:80` (проверьте фаервол завода).

## Вариант 2: напрямую с S7-1200

У S7-1200 нет встроенного HTTP-клиента. Используйте библиотеку Siemens **LHTTP** (Library for HTTP communication, SIOS) или собственную реализацию на `TCON` + `TSEND` / `TRCV`.

Требования:

- В свойствах CPU задан шлюз (gateway) с выходом в интернет.
- Соединение `TCON_IP_V4`: IP партнёра `168.113.176.218`, порт `80`, активное установление соединения со стороны ПЛК. TLS, DNS и сертификаты не нужны.

JSON собирается в `Array of Char` / `WString` (обычный `String` ограничен 254 символами). Числа — через `Real`→`String` с последующей заменой пробелов и запятых: формат должен соответствовать разделу «Правила формата». Поле `ts` проще не передавать.

Сырой HTTP-запрос, который должен уйти в сокет:

```http
POST /api/v1/ingest/readings HTTP/1.1
Host: 168.113.176.218
Authorization: Bearer <AGENT_API_KEY>
Content-Type: application/json
Content-Length: 157
Connection: close

{"plant_id":"plant-01","readings":[{"boiler_number":1,"level_mm":7301,"density":1.35,"volume_m3":159.98,"volume_l":159980,"mass_t":215.97,"income_month":0}]}
```

Строки заголовков разделяются `CR LF` (`$R$L` в SCL), после заголовков — пустая строка. `Content-Length` — точная длина тела в байтах.

## Настройка nginx (HTTP по IP)

В `infra/nginx/default.conf` есть отдельный блок `server { listen 80 default_server; server_name 168.113.176.218 _; ... }`. Он отдаёт сайт и `/api/` по HTTP без редиректа. Блок `server_name agromir-stat.online` на порту 80 по-прежнему редиректит на HTTPS и обслуживает проверку certbot.

После изменения конфига на сервере:

```bash
cd <путь к проекту>/infra
docker compose exec frontend nginx -t
docker compose exec frontend nginx -s reload
```

### Ограничение по IP завода

Чтобы по HTTP приём шёл только с завода, в блоке HTTP по IP добавьте отдельный location перед `location /api/`:

```nginx
location = /api/v1/ingest/readings {
    allow <внешний IP завода>;
    deny all;
    proxy_pass http://api:8000/api/v1/ingest/readings;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
}
```

## Проверка

С любого ПК (PowerShell):

```powershell
$body = '{"plant_id":"plant-01","readings":[{"boiler_number":1,"level_mm":7301,"density":1.35,"volume_m3":159.98,"volume_l":159980,"mass_t":215.97,"income_month":0,"outflow_month":0,"income_year":0,"outflow_year":0}]}'
Invoke-RestMethod -Method Post `
  -Uri http://168.113.176.218/api/v1/ingest/readings `
  -Headers @{ Authorization = "Bearer <AGENT_API_KEY>" } `
  -ContentType "application/json" -Body $body
```

Ожидаемый ответ: `accepted : 1`.

curl:

```bash
curl -X POST http://168.113.176.218/api/v1/ingest/readings \
  -H "Authorization: Bearer <AGENT_API_KEY>" \
  -H "Content-Type: application/json" \
  -d '{"plant_id":"plant-01","readings":[{"boiler_number":1,"level_mm":7301,"density":1.35,"volume_m3":159.98,"volume_l":159980,"mass_t":215.97}]}'
```

Когда пришёл последний пакет от завода (без авторизации):

```
GET http://168.113.176.218/api/v1/plants/plant-01/status
→ {"plant_id":"plant-01","last_accepted_at":"2026-11-11T17:30:00Z","last_count":24}
```

В логах API (`docker compose logs -f api`) на каждый принятый пакет пишется `accepted packet from plant-01: 24 records`.
