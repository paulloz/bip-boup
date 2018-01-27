module.exports.command = 'vote';
module.exports.help = 'Crée un sondage avec :thumbsup: et :thumbsdown:'

module.exports.callback = (message, words) => {
    if (words.length > 1) {
        message.channel.send(words.slice(1).join(' ')).then(_message => {
            _message.react('👍');
            _message.react('👎');

            if (message.deletable)
                message.delete();
        });
    }
};
