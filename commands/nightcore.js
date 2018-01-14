const Https = require('https');
const HtmlParser = require('htmlparser2');

module.exports.command = 'nightcore';
module.exports.help = 'Recherche du nightcore sur YouTube';

const URL = query => `https://www.youtube.com/results?sp=EgIQAQ%253D%253D&search_query=${query}`;
const RepliedURL = query => `https://www.youtube.com${query}`;
module.exports.callback = (message, words) => {
    if (words.length > 1) {
        Https.get(URL(words.join('+')), result => {
            result.setEncoding('utf8');

            let htmlBody = '';
            result.on('data', data => htmlBody += data);
            result.on('end', () => {
                let done = false;
                const parser = new HtmlParser.Parser({
                    onopentag: (name, attribs) => {
                        if (done) return;
                        if (name === 'a' && attribs.href.indexOf('/watch') === 0 && attribs.title != null) {
                            message.reply(RepliedURL(attribs.href));
                            done = true;
                        }
                    }
                });
                parser.write(htmlBody);
                parser.end();
            });
        });
    } else {
        message.reply('Tu cherches quoi ?');
    }
};
