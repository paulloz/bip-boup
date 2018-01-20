const chrono = require('chrono-node');
const moment = require('moment');

module.exports.command = 'rappel';
module.exports.help = 'Renvoie un message après un temps défini.';

module.exports.callback = (message, words) => {
    let messageContent;

    words = words.slice(1).join(' ');

    let matches = chrono.parse(words);
    if (matches.length > 0) {
        for (let match of matches) {
            if (Object.keys(match.tags).filter(k => k.startsWith('FR'))) {
                let sendAt = moment(chrono.parseDate(match.text, moment(message.createdTimestamp)));
                let time = sendAt - moment();

                if (messageContent == null) {
                    messageContent = words.substring(match.index + match.text.length);
                    if (match.index > messageContent.length)
                        messageContent = words.substring(0, match.index);
                }

                // TODO Do some smart things to clean the message
                messageContent = messageContent.trim();

                if (time > 0 && messageContent.length > 0) {
                    // TODO Store on disk
                    setTimeout(() => {
                        // TODO Remove from storage
                        message.channel.send(messageContent);
                    }, time);

                    message.react('👌');

                    return;
                }
            }
        }
    }

    message.react('👎');
};
