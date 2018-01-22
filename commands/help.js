const { GetInConf } = require('../utils.js');
const Config = require('../config.js');

module.exports.command = 'help';
module.exports.help = 'Montre ce message d\'aide.';
module.exports.setup = true;
module.exports.callback = (bipboup) => {
    return (message, words) => {
        const helpFor = (command) => `**${Config.get('attention', message.guild)}${command.command}** : ${command.help}`;
        if (words.length === 2) {
            const command = bipboup.config.commands.find(command => command.command === words[1]);
            if (command != null)
                message.reply('\n' + helpFor(command));
            else
                message.reply(`La commande ${attentionChar}${words[1]} n'existe pas.`);
        } else {
            // Reply with the full help text
            message.reply(['Voici la liste des commandes disponibles :'].concat(bipboup.config.commands.map(command => `\t- ${helpFor(command)}`)).join('\n'));
        }
    };
};
