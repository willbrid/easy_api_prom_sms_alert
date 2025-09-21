# Configuration

### Options de configuration

- **Mode binaire**

|Option          |Obligatoire|Description|
|----------------|-----------|-----------|
`--config-file       `|oui|option permettant de préciser l'emplacement du fichier de configuration
`--port`|non|option permettant de préciser le port (par défaut : `5957`)
`--enable-https      `|non|option permettant d'activer ou désactiver la communication tls (par défaut : `false`)
`--cert-file`|non|option permettant de préciser l'emplacement du fichier de certificat (obligatoire si l'option `--enable-https` est à `true`)
`--key-file`|non|option permettant de préciser l'emplacement du fichier de clé privée (obligatoire si l'option `--enable-https` est à `true`)

- **Mode conteneur**

|Variable d'environnement|Obligatoire|Description|
|------------------------|-----------|-----------|
`EASY_API_PROM_SMS_ALERT_CONFIG_FILE`|non|variable permettant de préciser l'emplacement du fichier de configuration dans le conteneur (par défaut: `/etc/easy-api-prom-sms-alert/config.yaml`). Il peut être écrasé avec un fichier externe si celui-ci est monté en volume avec le même nom et au même emplacement.
`EASY_API_PROM_SMS_ALERT_PORT`|non|variable permettant de préciser le port (par défaut : `5957`)
`EASY_API_PROM_SMS_ALERT_ENABLE_HTTPS`|non|variable permettant d'activer ou désactiver la communication tls (par défaut : `true`)
`EASY_API_PROM_SMS_ALERT_CERT_FILE`|non|variable permettant de préciser l'emplacement du fichier de certificat (obligatoire si la variable `EASY_API_PROM_SMS_ALERT_ENABLE_HTTPS` est à `true`, par défaut : `/etc/easy-api-prom-sms-alert/tls/server.crt`)
`EASY_API_PROM_SMS_ALERT_KEY_FILE`|non|variable permettant de préciser l'emplacement du fichier de clé privée (obligatoire si la variable `EASY_API_PROM_SMS_ALERT_ENABLE_HTTPS` est à `true`, par défaut : `/etc/easy-api-prom-sms-alert/tls/server.key`)

### Fichier de configuration

```
# Documentation sur le fichier de configuration
easy_api_prom_sms_alert:
  # Mode simulation du webhook : true -> les sms sont envoyés dans les logs et false (valeur en production) -> les sms sont envoyés 
  # via l'api du fournisseur
  simulation: true
  
  # Paramètre d'authentification au webhook
  auth:
    # Activation de l'authentification : true -> les paramètres username et password seront requis
    # Pour s'authentifier en header basic, il faudrait générer le base64 de la chaine username:password
    enabled: true
    # Nom d'utilisateur
    username: test
    # Mot de passe
    password: test@test

  # Paramètre du fournisseur de SMS
  provider:
    # Adresse de l'api du fournisseur
    url: "http://localhost:5797"
    # L'entête content-type acceptée par le fournisseur
    # Valeurs possibles : "application/json", "application/x-www-form-urlencoded"
    content_type: "application/json"
    # Paramètre d'activation de la vérification de certificat
    # - true -> le certificat de l'api en https du fournisseur ne sera pas vérifié
    # - false -> le certificat de l'api en https du fournisseur sera vérifié (valeur par défaut)
    insecure_skip_verify: false
    # Paramètre d'authentification à l'api du fournisseur
    authentication:
      # Activation de l'authentification à l'api du fournisseur : 
      # - true -> l'api du fournisseur nécessite une authentification et dans ce cas 
      #   les paramètres authorization_type et authorization_credential sont obligatoires
      # - false -> l'api du fournisseur ne nécessite pas d'authentification
      enabled: false
      # Paramètre d'autorisation d'entête http : Authorization
      # Type d'entête en exemple : Bearer, Basic, ApiKey
      authorization_type: ''
      # Chaine de caractères représentant la clé secret
      authorization_credential: ''
    # Paramètre des champs du corps de requête http de l'api    
    parameters:
      # Champ d'expéditeur
      from:
        # Nom du champ d'expéditeur
        param_name: "from"
        # Valeur du champ d'expéditeur
        param_value: "+1234567890"
        # méthode d'envoie du champ d'expéditeur : post ou query
        param_method: "post"
      # Champ du destinataire
      to:
        # Nom du champ destinataire
        param_name: "to"
        # Valeur par défaut du champ destinataire qui doit correspondre à l'un des noms des récipients configurés
        # dans le cas où le champ team est inexistant dans un champ alert
        param_value: "administration"
        # méthode d'envoie du champ destinataire : post ou query
        param_method: "query"
      # Champ du contenu du SMS
      message:
        # Nom du champ du contenu du SMS
        param_name: "content"
      # Paramètres supplémentaires du fournisseur. Ils peuvent être obligatoires ou non selon les spécifications d'intégration du fournisseur
      # Les valeurs ci-dessous sont des exemples. Il faudrait lire la documentation du fournisseur pour mieux configurer
      extra_params:
      - param_name: "pn1"
        param_value: "pv1"
        param_method: "post"
      - param_name: "pn2"
        param_value: "pv2"
        param_method: "query"
    # Paramètre de timeout à définir pour consommer l'api du fournisseur  (par défaut : 10s)
    timeout: 10s

  # Paramètre des différents destinataires qui recevront les SMS
  recipients:
  # nom de groupe du destinataire
  - name: "administration"
    # Membres du groupe de destinataire 
    members:
    - "+1234567890"
    - "+0987654321"
```