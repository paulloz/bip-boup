// Define how people will get the bot's attention
const attentionChar = '!';
// TODO Change this to enable commands after a mention
const attentionRegexp = new RegExp(`^${attentionChar}(.*)`);

module.exports = (bipboup, commands) => {
    // TODO Move this in a command file (but we need access to the `commands` array)
    const help = (message, words) => {
        // TODO check words to see if we need only part of the help text
        // Reply with the help text
        // TODO Order alphabeticaly
        message.reply([''].concat(commands.map(command => `!${command.command} : ${command.help}`)).join('\n'));
    };

    // Respond to messages
    bipboup.on('message', message => {
        // Make sure we de not reply to our own messages
        if (message.author.id == bipboup.user.id) return;

        let messageContent = message.cleanContent.trim().match(attentionRegexp);

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
