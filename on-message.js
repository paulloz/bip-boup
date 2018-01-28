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

            if (Config.hasCommand(messageContent[0]))
                Config.getCommand(messageContent[0]).call(message, messageContent);
        }
    });
};
