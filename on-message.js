// Define how people will get the bot's attention
const attentionChar = '!';
// TODO Change this to enable commands after a mention
const attentionRegexp = new RegExp(`^${attentionChar}(.*)`);

module.exports = (bipboup, commands) => {
    // TODO Move this in a command file (but we need access to the `commands` array)
    const help = (message, words) => {
        const helpFor = (command) => `${attentionChar}${command.command} : ${command.help}`;
        if (words.length === 2) {
            const command = commands.find(command => command.command === words[1]);
            if (command != null)
                message.reply('\n' + helpFor(command));
            else
                message.reply(`La commande ${attentionChar}${words[1]} n'existe pas.`);
        } else {
            // Reply with the full help text
            // TODO Order alphabeticaly
            message.reply(['Voici la liste des commandes disponibles :'].concat(commands.map(command => helpFor(command))).join('\n'));
        }
    };

    // Respond to messages
    bipboup.on('message', message => {
        // Make sure we de not reply to our own messages
        if (message.author.id == bipboup.user.id) return;

        let messageContent = message.content.trim().match(attentionRegexp);

        if (messageContent != null) {
            // Split message content into words
            messageContent = messageContent[1].split(/\s+/).filter(word => word.length > 0);
            if (messageContent.length <= 0) return; // Safety

            if (messageContent[0] == 'help') {
                help(message, messageContent);
            } else {
                for (let command of commands) {
                    if (messageContent[0] == command.command) {
                        command.callback(message, messageContent);
                        break;
                    }
                }
            }
        }
    });
};
