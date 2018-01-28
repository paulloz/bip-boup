const { Permissions } = require('discord.js');
const Config = require('./config.js');

module.exports = (bipboup) => {
    // TODO Change this to enable commands after a mention
    const attentionRegexp = (char) => new RegExp(`^[${char}](\\w+)(\\s+.+)*`);

    // Respond to messages
    bipboup.on('message', message => {
        // Make sure we de not reply to our own messages
        if (message.author.id == bipboup.user.id) return;

        let messageContent = message.content.trim().match(attentionRegexp(Config.get('attention', message.guild)));

        if (messageContent != null) {
            // Split message content into words
            messageContent = [messageContent[1].trim()].concat((messageContent[2] || "").split(/\s+/).filter(word => word.length > 0));
            if (messageContent.length <= 0) return; // Safety

            const command = Config.hasCommand(messageContent[0]) ? Config.getCommand(messageContent[0]) : null;
            // TODO Check permissions by guild
            if (command != null && (!command.isAdmin() || bipboup.guilds.first().member(message.author).hasPermission(Permissions.FLAGS.ADMINISTRATOR)))
                command.call(message, messageContent);
        }
    });
};
