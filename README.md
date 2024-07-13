# Easy_api_prom_sms_alert

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://opensource.org/licenses/MIT)

**Easy_api_prom_sms_alert** est un webhook pour prometheus permettant d'envoyer des alertes sms avec n'importe quel fournisseur de sms exposant une api.

## Problème

Lorsque **Prometheus** détecte des conditions anormales dans les systèmes surveillés, il déclenche des alertes pour informer **alertmanager** d'effectuer des notifications SMS. Cependant il existe plusieurs types de fournisseur de SMS avec leur propre spécification. Ainsi intégrer quelques uns dans alertmanager, rendrait la configuration complexe à gérer. 

## Solution

Avec **Easy_api_prom_sms_alert**, les utilisateurs auront la possibilité de choisir n'importe quel fournisseur de services SMS qui expose une API en **HTTP POST**. Cela leur donnera la liberté de sélectionner le fournisseur qui répond au mieux à leurs besoins en termes de coût, de fiabilité et de couverture géographique.

## Documentation

1- [Installation](https://github.com/willbrid/easy_api_prom_sms_alert/blob/main/docs/installation.md) <br>
2- [Utilisation](https://github.com/willbrid/easy_api_prom_sms_alert/blob/main/docs/utilisation.md) <br>
3- [Fichier de configuration](https://github.com/willbrid/easy_api_prom_sms_alert/blob/main/docs/configuration.md) <br>
4- [Exemple complet](https://github.com/willbrid/easy_api_prom_sms_alert/blob/main/docs/exemple.md)

## Licence

Ce projet est sous licence MIT - voir le fichier [LICENSE](https://opensource.org/licenses/MIT) pour plus de détails.