const Https = require('https');

module.exports.EmojiOrNothing = (channel, emoji) => channel.guild != null ? channel.guild.emojis.find('name', emoji) || '' : ''

module.exports.HttpsGet = (url, callback) => {
    return Https.get(url, result => {
        result.setEncoding('utf8');

        let htmlBody = '';
        result.on('data', data => htmlBody += data);
        result.on('end', () => {
            callback(htmlBody);
        });
    });
};
