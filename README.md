# Easy-api-prom-alert-sms

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://opensource.org/licenses/MIT)

**Easy-api-prom-alert-sms** est un webhook pour prometheus permettant d'envoyer des alertes sms avec n'importe quel fournisseur de sms exposant une api.

## Problème

Lorsque **Prometheus** détecte des conditions anormales dans les systèmes surveillés, il déclenche des alertes pour informer **alertmanager** d'effectuer des notifications SMS. Cependant il existe plusieurs types de fournisseur de SMS avec leur propre spécification. Ainsi intégrer quelques uns dans alertmanager, rendrait la configuration complexe à gérer. 

## Solution

Avec **Easy-api-prom-alert-sms**, les utilisateurs auront la possibilité de choisir n'importe quel fournisseur de services SMS qui expose une API en **HTTP POST**. Cela leur donnera la liberté de sélectionner le fournisseur qui répond au mieux à leurs besoins en termes de coût, de fiabilité et de couverture géographique.

## Licence

Ce projet est sous licence MIT - voir le fichier [LICENSE](LICENSE) pour plus de détails.