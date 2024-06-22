# Fichier de configuration

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