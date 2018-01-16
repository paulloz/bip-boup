const Https = require('https');
const HtmlParser = require('htmlparser2');

module.exports.command = 'fanfiction';
module.exports.help = 'Recherche une fanfiction sur fanfiction.net';

const URL = query => `https://www.fanfiction.net/search/?ready=1&keywords=${query}&categoryid=0&genreid1=0&genreid2=0&languageid=&3censorid=0&statusid=0&type=story&match=&sort=&ppage=1&characterid1=0&characterid2=0&characterid3=0&characterid4=0&words=0&formatid=0`;
const RepliedURL = query => `https://www.fanfiction.net${query}`;
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
                        if (name === 'a' && attribs.href.indexOf('/s') === 0) {
                            message.reply(RepliedURL(attribs.href));
                            done = true;
                        }
                    }
                });
                parser.write(htmlBody);
                parser.end();

                if (!done)
                    message.reply('Je n\'ai rien trouvé de satisfaisant :frowning:')
            });
        });
    } else {
        message.reply('Tu cherches quoi ?');
    }
};