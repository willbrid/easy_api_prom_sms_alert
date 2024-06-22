# Installation

Ces actions sont effectuées sous un serveur linux.

#### Via le package sous Linux

```
cd $HOME && mkdir -p alert-prometheus && cd alert-prometheus
```

```
curl -LO https://github.com/willbrid/easy_api_prom_sms_alert/releases/download/v<VERSION>/easy_api_prom_sms_alert_<VERSION>_linux_amd64.tar.gz
```

```
tar -xvzf easy_api_prom_sms_alert_<VERSION>_linux_amd64.tar.gz
```

```
vi config.yaml
```

```
easy_api_prom_sms_alert:
  simulation: true
  auth:
    enabled: true
    username: test
    password: test@test
  provider:
    url: "http://localhost:5797"
    content_type: "application/json"
    authentication:
      enabled: false
      authorization_type: ''
      authorization_credential: ''
    parameters: 
      from: 
        param_name: "from"
        param_value: "+1234567890"
        param_method: "post"
      to:
        param_name: "to"
        param_value: "administration"
        param_method: "query"
      message: 
        param_name: "content"
    timeout: 0s
  recipients: 
  - name: "administration"
    members:
    - "+1234567890"
    - "+0987654321"
  - name: "urgence"
    members:
    - "+1122334455"
    - "+5544332211"
```

```
./easy_api_prom_sms_alert_<VERSION>_linux_amd64 --config-file ./config.yaml
```

Remplacez **\<VERSION\>** par le numéro de version souhaité (supérieur ou égal à **1.1.8**).

#### Via docker

--- **Installation en utilisant le fichier de configuration par défaut et en activant le protocole https**

```
docker run -d -p 8000:5957 --name alert-sms-sender -e EASY_API_PROM_SMS_ALERT_PORT=5957 -e EASY_API_PROM_SMS_ALERT_ENABLE_HTTPS=true willbrid/easy-api-prom-sms-alert:latest
```

Dans cet exemple, le port par défaut **5957** interne au container est mappé au port externe **8000**. 

--- **Installation en utilisant un volume persistent pour le fichier config.yaml et en activant le protocole https**

L'idée ici est de pouvoir permettre la personnalisation du fichier de configuration **config.yaml**

```
cd $HOME && mkdir -p alert-prometheus && cd alert-prometheus
```

```
vi config.yaml
```

```
easy_api_prom_sms_alert:
  simulation: true
  auth:
    enabled: true
    username: test
    password: test@test
  provider:
    url: "http://localhost:5797"
    content_type: "application/json"
    authentication:
      enabled: false
      authorization_type: ''
      authorization_credential: ''
    parameters: 
      from: 
        param_name: "from"
        param_value: "+1234567890"
        param_method: "post"
      to:
        param_name: "to"
        param_value: "administration"
        param_method: "query"
      message: 
        param_name: "content"
    timeout: 0s
  recipients: 
  - name: "administration"
    members:
    - "+1234567890"
    - "+0987654321"
  - name: "urgence"
    members:
    - "+1122334455"
    - "+5544332211"
```

```
docker run -d -p 8000:5957 --name alert-sms-sender -v $HOME/alert-prometheus/config.yaml:/etc/easy-api-prom-sms-alert/config.yaml -e EASY_API_PROM_SMS_ALERT_PORT=5957 -e EASY_API_PROM_SMS_ALERT_ENABLE_HTTPS=true willbrid/easy-api-prom-sms-alert:latest
```
