const Config = require('./config.js');

module.exports = (bipboup) => {
    bipboup.on('emojiCreate', (emoji) => {
        let chan = Config.get('mainchan', emoji.guild);
        if (chan != null)
            bipboup.channels[chan].send(`Nouvel emote : ${emoji} !`);
    });

    bipboup.on('emojiDelete', (emoji) => {
        let chan = Config.get('mainchan', emoji.guild);
        if (chan != null)
            bipboup.channels[chan].send(`Emote supprimé : ${emoji}`);
    });
};
