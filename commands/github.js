const { HttpsGetJson, BooleanEmoji } = require('../utils.js');

module.exports.command  = 'github'
module.exports.help     = 'Envoie l\'adresse du dépôt GitHub de mon code source.'

const baseAPI = "https://api.github.com/repos/paulloz/bip-boup";
const basePullRequestAPI = `${baseAPI}/pulls`;

module.exports.callback = (message, words) => {
    const formatPR = (pr) => `#${pr.number} ${pr.title} par ${pr.user.login} ${BooleanEmoji(pr.mergeable_state === 'clean')}, <${pr.html_url}>`;

    const handlePRs = () => {
        if (words.length <= 2) {
            HttpsGetJson(basePullRequestAPI, (json) => {
                let reply = `Il y a actuellement ${json.length} pull request${json.length > 1 ? 's' : ''} en attente${json.length > 0 ? ' :' : ''}`;
                for (let item of json)
                    reply += `\n\t* ${formatPR(item)}`
                message.channel.send(reply);
            });
        }
    };

    if (words.length <= 1) {
        message.reply(`<https://github.com/paulloz/bip-boup.git>`);
    } else {
        switch (words[1].toLowerCase()) {
            case "pr": handlePRs(); break;
            default: break;
        }
    }
};
