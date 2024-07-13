# Utilisation

Notre fichier de configuration précédent montre que l'authentification **Basic** est activée. Ainsi il faudrait générer le **base64** de la chaine **username:password** afin de l'utiliser dans le header **Authorization**.

```
echo -n test:test@test | base64
```

```
dGVzdDp0ZXN0QHRlc3Q=
```

#### **Test à effectuer avec curl**

```
curl -k --location 'https://localhost:8000/api-alert' \
--header 'Content-Type: application/json' \
--header 'Authorization: Basic dGVzdDp0ZXN0QHRlc3Q=' \
--data '{
  "version": "4",
  "groupKey": "{alertname=\"InstanceDown\"}",
  "status": "firing",
  "receiver": "webhook",
  "groupLabels": {
    "alertname": "InstanceDown"
  },
  "commonLabels": {
    "alertname": "InstanceDown",
    "job": "myjob",
    "severity": "critical"
  },
  "commonAnnotations": {
    "summary": "Instance xxx down",
    "description": "The instance xxx is down."
  },
  "externalURL": "http://prometheus.example.com",
  "alerts": [
    {
      "status": "firing",
      "labels": {
        "alertname": "InstanceDown",
        "instance": "localhost:9090",
        "team": "urgence",
        "job": "myjob",
        "severity": "critical"
      },
      "annotations": {
        "summary": "Instance localhost:9090 down",
        "description": "The instance localhost:9090 is down."
      },
      "startsAt": "2023-05-20T14:28:00.000Z",
      "endsAt": "0001-01-01T00:00:00Z",
      "generatorURL": "http://prometheus.example.com/graph?g0.expr=up%3D%3D0&g0.tab=1"
    }
  ]
}'
```

#### **Intégration dans Alertmanager**

Pour intégrer **Easy_api_prom_sms_alert** dans **Alertmanager**, vous devez configurer un webhook en ajoutant un récepteur webhook à votre configuration d' **Alertmanager**.

```
receivers:
- name: 'admin'
  webhook_configs:
  - url: 'https://localhost:8000/api-alert'
    send_resolved: false
    http_config: 
      basic_auth:
        username: test
        password: test@test
      tls_config:
        insecure_skip_verify: true
```

Pour visualiser le résultat en mode **simulation**, vous devez consulter les logs du conteneur.

```
docker container logs alert-sms-sender
```