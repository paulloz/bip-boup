const chrono = require('chrono-node');
const moment = require('moment');

module.exports.command = 'rappel';
module.exports.help = 'Renvoie un message après un temps défini.';

module.exports.callback = (message, words) => {
    const parseAndHandle = (toParse) => {
        // TODO Do some smart things to clean the message
        const getMessageText = (full, sIdx, eIdx) => full.substring(...((sIdx > full.length - eIdx) ? [0, sIdx] : [eIdx])).trim();

        const getChronoMatches = (toParse, sentAt) => chrono.parse(toParse).filter(
            m => Object.keys(m.tags).filter(k => k.startsWith('FR')).length > 0
        ).map(m => Object({
            sendAt: moment(chrono.parseDate(m.text, sentAt)),
            text: getMessageText(toParse, m.index, m.index + m.text.length)
        })).filter(m => m.sendAt > sentAt && m.text.length > 0);

        let match;
        if ((match = getChronoMatches(toParse, moment(message.createdTimeStamp)).shift()) != null) {
            // TODO Store on disk
            setTimeout(() => {
                // TODO Remove from storage
                message.channel.send(match.text);
            }, match.sendAt - moment());

            return true;
        }
        return false;
    };

    message.react(parseAndHandle(words.slice(1).join(' ')) ? '👌' : '👎');
};
