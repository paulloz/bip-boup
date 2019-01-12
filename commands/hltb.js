let hltb = require('howlongtobeat');
let hltbService = new hltb.HowLongToBeatService();

module.exports.command = 'hltb'
module.exports.help = 'HowLongToBeat'
module.exports.callback = async (message, words) => {
    if (words.length > 1) {
        const result = await hltbService.search(words.slice(1).join(' '))
        if (result.length == 0) message.reply("oh no, j'ai rien trouvé ☹")
        else {
            game = result[0]
            const {
                gameplayMain: mainStory,
                gameplayCompletionist: completionist,
            } = game
            let res = game.name
            if (mainStory || completionist) res += " se finit en "
            else return "Pas d'info pour " + res
            if (mainStory) res += mainStory + "h si tu traces" + (completionist ? ", " : ".")
            if (completionist) res += completionist + "h si tu veux tout faire."
            message.reply(res)
        }
    } else {
        message.reply('Tu cherches quoi ?');
    }
};
