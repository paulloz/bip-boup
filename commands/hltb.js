let hltb = require('howlongtobeat');
let hltbService = new hltb.HowLongToBeatService();

module.exports.command = 'hltb'
module.exports.help = 'HowLongToBeat'
module.exports.callback = (message, words) => {
    const getRes = (stg) => {
        if (stg.length < 1) return "oh no, j'ai rien trouvé ☹"
        let game = stg[0]
        let res = game.name
        let mainStory = game.gameplayMain
        let completionist = game.gameplayCompletionist
        if (mainStory || completionist) res += " se finit en "
        else return "Pas d'info pour " + res
        if (mainStory) res += mainStory + "h si tu traces" + (completionist ? ", " : ".")
        if (completionist) res += completionist + "h si tu veux tout faire."
        return res;
    }
    if (words.length > 1) {
        hltbService.search(words.slice(1).join(' ')).then(result => message.reply(getRes(result)))
    } else {
        message.reply('Tu cherches quoi ?');
    }
};
