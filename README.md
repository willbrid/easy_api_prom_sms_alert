# Easy-api-prom-alert-sms

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://opensource.org/licenses/MIT)

**Easy-api-prom-alert-sms** est un webhook pour prometheus permettant d'envoyer des alertes sms avec n'importe quel fournisseur de sms exposant une api.

## Problème

Lorsque **Prometheus** détecte des conditions anormales dans les systèmes surveillés, il déclenche des alertes pour informer **alertmanager** d'effectuer des notifications SMS. Cependant il existe plusieurs types de fournisseur de SMS avec leur propre spécification. Ainsi intégrer quelques uns dans alertmanager, rendrait la configuration complexe à gérer. 

## Solution

Avec **Easy-api-prom-alert-sms**, les utilisateurs auront la possibilité de choisir n'importe quel fournisseur de services SMS qui expose une API en **HTTP POST**. Cela leur donnera la liberté de sélectionner le fournisseur qui répond au mieux à leurs besoins en termes de coût, de fiabilité et de couverture géographique.

## Installation

- Via docker

--- **Installation en utilisant le fichier de configuration par défaut**

```
docker run -d -p 8000:5957 --name alert-sms-sender willbrid/easy-api-prom-sms-alert:v1.0.0 --config-file /etc/easy-api-prom-sms-alert/config.yaml --port 5957
```

Dans cet exemple, le port par défaut **5957** interne au container est mappé au port externe **8000**. 

--- **Installation en utilisant un volume persistent pour le fichier config.yaml et permettre sa configuration**

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
    authentication:
      enabled: false
      authorization: 
        type: ''
        credential: ''
    fields: 
      from: "from"
      from_value: "+1234567890"
      to: "to"
      message: "content"
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
docker run -d -p 8000:5957 -v $HOME/alert-prometheus/config.yaml:/etc/easy-api-prom-sms-alert/config.yaml --name alert-sms-sender willbrid/easy-api-prom-sms-alert:v1.0.0 --config-file /etc/easy-api-prom-sms-alert/config.yaml --port 5957
```

## Utilisation

- **Test à effectuer avec curl**

```
curl --location 'http://localhost:8000/api-alert' \
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

- **Intégration dans Alertmanager**

Pour intégrer **Easy-api-prom-alert-sms** dans **Alertmanager**, vous devez configurer un webhook en ajoutant un récepteur webhook à votre configuration d' **Alertmanager**.

```
receivers:
- name: 'admin'
  webhook_configs:
  - url: 'http://localhost:8000/api-alert'
    send_resolved: false
    http_config: 
      authorization:
        type: "Basic"
        credentials: dGVzdDp0ZXN0QHRlc3Q=
```

Pour visualiser le résultat en mode **simulation**, vous devez consulter les logs du conteneur.

```
docker container logs alert-sms-sender
```

## Fichier de configuration config.yaml

```
# Documentation sur le fichier de configuration
easy_api_prom_sms_alert:
  # Mode simulation du webhook : true -> les sms sont envoyés dans les logs et false (valeur en production) -> les sms sont envoyés via l'api
  simulation: true
  
  # Paramètre d'authentification au webhook
  auth:
    # Activation de l'authentification : true -> les paramètres username et password seront requis
    enabled: true
    # Nom d'utilisateur
    username: test
    # Mot de passe
    password: test@test

  # Paramètre du fournisseur de SMS
  provider:
    # Adresse de l'api du fournisseur
    url: "http://localhost:5797"
    # Paramètre d'authentification à l'api du fournisseur
    authentication:
      # Activation de l'authentification à l'api du fournisseur : 
      # - true -> l'api du fournisseur nécessite une authentification et dans ce cas la section authentication.authorization est obligatoire
      # - false -> l'api du fournisseur ne nécessite pas d'authentification
      enabled: false
      # Paramètre d'autorisation d'entête http : Authorization
      authorization:
        # Type d'entête parmi : Bearer, Basic, ApiKey
        type: ''
        # Chaine de caractères représentant la clé secret
        credential: ''
    # Paramètre des champs du corps de requête http de l'api    
    fields: 
      # Champ d'expéditeur
      from: "from"
      # Valeur du champ d'expéditeur
      from_value: "+1234567890"
      # Champ du destinataire
      to: "to"
      # Champ du contenu du SMS
      message: "content"
    # Paramètre de timeout à définir pour consommer l'api du fournisseur  
    timeout: 0s

  # Paramètre des différents destinataires qui recevront les SMS
  recipients:
  # nom de groupe du destinataire
  - name: "administration"
    # Membres du groupe de destinataire 
    members:
    - "+1234567890"
    - "+0987654321"
```

## Licence

Ce projet est sous licence MIT - voir le fichier [LICENSE](LICENSE) pour plus de détails.