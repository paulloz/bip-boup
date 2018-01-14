module.exports.command = 'vote';
module.exports.help = 'Crée un sondage avec :thumbsup: et :thumbsdown:'

module.exports.callback = (message, words) => {
    if (words.length > 1) {
        message.channel.send(words.slice(1).join(' ')).then(message => {
            message.react('👍');
            message.react('👎');
        });
    }

    if (message.deletable)
        message.delete();
};
