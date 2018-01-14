# bip-boup
A Discord bot doing ???

## Fonctionnement

* Lancer `npm install`
* Créer un fichier `.token` contenant le *Bot Token* de [l'application Discord](https://discordapp.com/developers/applications/me/)
* Lancer `npm start`
* Apprécier le bip boup

## Ajouter des commandes

Chaque fichier contenu dans le dossier [commands/](/commands) est chargé au lancement du robot.  
Les fichiers de commandes doivent exporter les propriétés suivantes :
* `command` : le mot devant déclencher la commande
* `help` : le texte d'aide accompagnant la commande
* `callback` : la fonction devant être appelée pour traiter la commandes

Exemple :
```javascript
module.exports.command  = 'spreadlove'
module.exports.help     = 'Répendre un peu d\'amour dans ce monde de brutes.'
module.exports.callback = (message, words) => {
    message.reply('❤️❤️❤️');
};
```

### Test

Il est possible de tester les commandes hors de l'environnement Discord via la commande `npm run test-command` :
```sh
$> nom run test-command foo bar

> bip-boup@x.x.x test-command .../bip-boup
> node tests/test-command.js "foo" "bar"

...   # <- Résultats simulé d'un message discord contenant "!foo bar"
```

```sh
$> nom run test-command foo

> bip-boup@x.x.x test-command .../bip-boup
> node tests/test-command.js "foo"

> bar # <- Écrire son message ici
...   # <- Résultats simulé d'un message discord contenant "!foo bar"
```

## Licence

Ce programme est publié sous licence MIT. Pour plus d'information, se référer au fichier [LICENSE](/LICENSE).
