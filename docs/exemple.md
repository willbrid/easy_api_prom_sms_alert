# Exemple complet

### Mise en place d'une sandbox

Nous mettrons en place un serveur vagrant **Rocky linux 8** sous une machine hôte Ubuntu 20.04 (ou plus) avec déjà **virtualbox7** et **vagrant** installé.

```
mkdir $HOME/easy_api_prom_sms_alert && cd $HOME/easy_api_prom_sms_alert
```

```
wget https://download.virtualbox.org/virtualbox/7.0.12/VBoxGuestAdditions_7.0.12.iso
```

```
vi Vagrantfile
```

```
# -*- mode: ruby -*-
# vi: set ft=ruby :

VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.vbguest.auto_update = false
  config.vbguest.no_remote = true
  config.vbguest.iso_path = "./VBoxGuestAdditions_7.0.12.iso"

  # General Vagrant VM configuration
  config.ssh.insert_key = false
  config.vm.synced_folder ".", "/vagrant", disabled: true
  config.vm.provider :virtualbox do |v|
    v.memory = 4096
    v.cpus = 2
    v.linked_clone = true
  end

  # Monitoring Server
  config.vm.define "monitoring-server" do |srv|
    srv.vm.box = "willbrid/rockylinux8"
    srv.vm.box_version = "0.0.2"
    srv.vm.hostname = "monitoring-server"
    srv.vm.network :private_network, ip: "192.168.56.211"
  end
end
```

```
vagrant up
```

### Installation de prometheus, alertmanager, node-exporter et easy_api_prom_sms_alert sur notre serveur monitoring-server

Dans cette section, nous allons installer sous forme conteneurisée sur le serveur vagrant **monitoring-server** : 
- **prometheus** : pour le monitoring
- **alertmanager** : pour recevoir les alertes prometheus et déclenchées les notifications
- **node-exporter** : pour collecter les métriques systèmes sur le serveur vagrant **monitoring-server**
- **easy_api_prom_sms_alert** : webhook qui sera configuré au niveau d'alertmanager qui lui permettra d'envoyer des alertes SMS

```
vagrant ssh monitoring-server
```

##### Mise en place de node-exporter

```
podman run -d --net="host" \
       --pid="host" \
       -v "/:/host:ro,rslave" \
       -u root \
       --name node-exporter \
       quay.io/prometheus/node-exporter:v1.7.0 \
       --path.rootfs=/host
```

##### Mise en place d'alertmanager

```
mkdir -p $HOME/monitoring/alertmanager && mkdir $HOME/monitoring/alertmanager/data && cd $HOME/monitoring/alertmanager
```

```
vi alertmanager.yml
```

```
route:
  receiver: 'admin'
  repeat_interval: 1h

receivers:
- name: 'admin'
  webhook_configs:
  - url: 'https://192.168.56.211:5797/api-alert'
    send_resolved: false
    http_config: 
      basic_auth:
        username: test
        password: test@test
      tls_config:
        insecure_skip_verify: true
```

```
podman run -d --net=host \
       -v $HOME/monitoring/alertmanager/alertmanager.yml:/config/alertmanager.yml:z \
       -v $HOME/monitoring/alertmanager/data:/data:z \
       --name alertmanager \
       prom/alertmanager:v0.26.0 \
       --config.file=/config/alertmanager.yml \
       --log.level=debug
```

##### Mise en place de Prometheus

```
mkdir -p $HOME/monitoring/prometheus && mkdir $HOME/monitoring/prometheus/data && mkdir $HOME/monitoring/prometheus/rules && cd $HOME/monitoring/prometheus
```

```
vi prometheus.yml
```

```
global:
  scrape_interval:     10s 
  evaluation_interval: 10s
  external_labels:
    cluster: CLUSTER_A
    replica: 0

rule_files:
- "/etc/prometheus/rules/*"

alerting:
 alertmanagers:
 - static_configs:
   - targets: ['192.168.56.211:9093']

scrape_configs:
  - job_name: 'prometheus'
    scrape_interval: 10s
    static_configs:
    - targets: ['192.168.56.211:9090']

  - job_name: 'node-exporter'
    static_configs:
    - targets: ['192.168.56.211:9100']
```

```
vi rules/nodes-rules.yml
```

```
groups:
- name: NODE_CLUSTER_A
  rules:
  - alert: NodeDown
    expr: up{job="node-exporter"} == 0
    for: 1m
    labels:
      severity: critical
      team: urgence
    annotations:
      summary: "Node is down"
      description: "The node has been down for the last 1 minute."
```

```
podman run -d --net=host \
    -v $HOME/monitoring/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml:z \
    -v $HOME/monitoring/prometheus/data:/prometheus:z \
    -v $HOME/monitoring/prometheus/rules:/etc/prometheus/rules:z \
    -u root \
    --name prometheus \
    quay.io/prometheus/prometheus:v2.38.0 \
    --config.file=/etc/prometheus/prometheus.yml \
    --storage.tsdb.retention.time=4h \
    --storage.tsdb.path=/prometheus \
    --storage.tsdb.max-block-duration=2h \
    --storage.tsdb.min-block-duration=2h \
    --web.listen-address=:9090 \
    --web.external-url=http://192.168.56.211:9090 \
    --web.enable-lifecycle \
    --web.enable-admin-api
```

##### Mise en place de easy_api_prom_sms_alert

```
mkdir -p $HOME/monitoring/alert && cd $HOME/monitoring/alert
```

Par la suite, il faudrait utiliser l'un des fichiers de configuration (à personnaliser) ci-dessous pour paramètrer l'intégration avec un fournisseur : **Twilio**, **WhatsApp**,... 

```
vi $HOME/monitoring/alert/config.yaml
```

**Intégration avec Twilio**
```
easy_api_prom_sms_alert:
  simulation: false
  auth:
    enabled: true
    username: test
    password: test@test
  provider:
    url: "https://api.twilio.com/2010-04-01/Accounts/XXXXXXX/Messages"
    content_type: "application/x-www-form-urlencoded"
    authentication:
      enabled: true
      authorization_type: 'Basic'
      authorization_credential: 'YYYYYYY'
    parameters: 
      from: 
        param_name: "From"
        param_value: "+xxxxxxx"
        param_method: "post"
      to:
        param_name: "To"
        param_value: "urgence"
        param_method: "post"
      message: 
        param_name: "Body"
    timeout: 0s
  recipients: 
  - name: "urgence"
    members:
    - "+yyyyyyy"
    - "+zzzzzzz"
```

**XXXXXXX** est la chaine **SID** à récupérer sur la plateforme **Twilio**. <br>
**YYYYYYY** est le base64 de la chaine **SID:TOKEN** à récupérer sur la plateforme **Twilio**. <br>
**+xxxxxxx** est le numéro émetteur. <br>
**+yyyyyyy** et **+zzzzzzz** sont les numéros de téléphone qui recevront les alertes sms.

**Référence** : [Twilio Documentation](https://www.twilio.com/en-us/blog/send-sms-twilio-shell-script-curl)

**Intégration avec WhatsApp**
```
easy_api_prom_sms_alert:
  simulation: false
  auth:
    enabled: true
    username: test
    password: test@test
  provider:
    url: "https://api.wassenger.com/v1/messages"
    content_type: "application/json"
    authentication:
      enabled: true
      authorization_type: 'Token'
      authorization_credential: 'API_TOKEN'
    parameters: 
      from: 
        param_name: "reference"
        param_value: "prometheus"
        param_method: "post"
      to:
        param_name: "phone"
        param_value: "urgence"
        param_method: "post"
      message: 
        param_name: "message"
    timeout: 0s
  recipients: 
  - name: "urgence"
    members:
    - "+xxxxxxx"
    - "+yyyyyyy"
```

**API_TOKEN** est le token récupérable sur la plateforme **Wassenger**. <br>
**+xxxxxxx** et **+yyyyyyy** sont les comptes whatsapp qui recevront les alertes sms.

**Référence** : [Wassenger Documentation](https://app.wassenger.com/docs/)

- **Démarrage de easy_api_prom_sms_alert avec l'un des contenus d'intégration**

```
podman run -d --net=host \
       --name alert-sms-sender \
       -v $HOME/monitoring/alert/config.yaml:/etc/easy-api-prom-sms-alert/config.yaml:z 
       -e EASY_API_PROM_SMS_ALERT_PORT=5957 
       -e EASY_API_PROM_SMS_ALERT_ENABLE_HTTPS=true 
       willbrid/easy-api-prom-sms-alert:latest
```

### Test

Afin de simuler l'arrêt du serveur **monitoring-server**, nous stoppons le container **node-exporter**.

```
podman container stop node-exporter
```

Après une minute, nous verrons une alerte sms sur notre téléphone (**intégration avec Twilio**) ou sur notre compte whatsApp (**intégration avec WhatsApp**).