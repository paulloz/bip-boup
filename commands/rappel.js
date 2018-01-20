const chrono = require('chrono-node');
const moment = require('moment');

module.exports.command = 'rappel';
module.exports.help = 'Renvoie un message après un temps défini.';

const quotes = ['\'', '"', '`', '«»', '“”', '‘’', '‹›'];

module.exports.callback = (message, words) => {
    let messageContent;

    words = words.slice(1).join(' ');

    // TODO This is dumb, do better
    if (messageContent == null) {
        let quoted = [-1, -1];

        for (let quote of quotes) {
            let startQuote = quote[0], endQuote = quote[quote.length > 1 ? quote.length - 1 : 0];
            let startIdx = words.indexOf(startQuote);
            if (startIdx >= 0 && (quoted[0] < 0 || startIdx < quoted[0])) {
                let endIdx = words.lastIndexOf(endQuote)
                if (endIdx >= 0 && endIdx > startIdx) {
                    quoted = [startIdx, endIdx];
                }
            }
        }

        if (quoted[0] >= 0 && quoted[1] >= 0) {
            messageContent = words.substring(quoted[0], quoted[1] + 1);
            words = ((quoted[0] > words.length - quoted[1]) ? words.substring(0, quoted[0]) : words.substring(quoted[1] + 1)).trim();
        }
    }

    let matches = chrono.parse(words);
    if (matches.length > 0) {
        for (let match of matches) {
            if (Object.keys(match.tags).filter(k => k.startsWith('FR'))) {
                let sendAt = moment(chrono.parseDate(match.text, moment(message.createdTimestamp)));

                if (messageContent == null) {
                    messageContent = words.substring(match.index + match.text.length);
                    if (match.index > messageContent.length)
                        messageContent = words.substring(0, match.index);
                }

                // TODO Do some smart things to clean the message
                messageContent = messageContent.trim();

                if (messageContent.length > 0) {
                    // TODO Store on disk
                    setTimeout(() => {
                        // TODO Remove from storage
                        message.channel.send(messageContent);
                    }, sendAt - moment());

                    message.react('👌');

                    return;
                }
            }
        }
    }

    message.react('👎');
};
