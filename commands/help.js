module.exports.command = 'help';
module.exports.help = 'Montre ce message d\'aide.';
module.exports.setup = true;
module.exports.callback = (config) => {
    return (message, words) => {
        const helpFor = (command) => `**${config.attentionChar}${command.command}** : ${command.help}`;
        if (words.length === 2) {
            const command = config.commands.find(command => command.command === words[1]);
            if (command != null)
                message.reply('\n' + helpFor(command));
            else
                message.reply(`La commande ${attentionChar}${words[1]} n'existe pas.`);
        } else {
            // Reply with the full help text
            // TODO Order alphabeticaly
            message.reply(['Voici la liste des commandes disponibles :'].concat(config.commands.map(command => `\t- ${helpFor(command)}`)).join('\n'));
        }
    };
};
