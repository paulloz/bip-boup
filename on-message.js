module.exports = (bipboup) => {
    // TODO Change this to enable commands after a mention
    const attentionRegexp = new RegExp(`^${bipboup.config.attentionChar}(.*)`);

    // Respond to messages
    bipboup.on('message', message => {
        // Make sure we de not reply to our own messages
        if (message.author.id == bipboup.user.id) return;

        let messageContent = message.content.trim().match(attentionRegexp);

        if (messageContent != null) {
            // Split message content into words
            messageContent = messageContent[1].split(/\s+/).filter(word => word.length > 0);
            if (messageContent.length <= 0) return; // Safety

            for (let command of bipboup.config.commands) {
                if (messageContent[0] == command.command) {
                    command.callback(message, messageContent);
                    break;
                }
            }
        }
    });
};
