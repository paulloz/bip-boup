const Config = require('./config.js');

module.exports = (bipboup) => {
    bipboup.on('guildMemberAdd', (member) => {
        let chan = Config.get('mainchan', member.guild);
        if (chan != null)
            bipboup.channels.get(chan).send(`Tiens, une nouvelle tête. Bonjour ${member} !`);
    });
};
