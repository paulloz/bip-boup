const { HttpsGetJson, BooleanEmoji } = require('../utils.js');

module.exports.command  = 'github'
module.exports.help     = 'Quelques informations sur mon code source.'

const baseAPI = "https://api.github.com/repos/paulloz/bip-boup";
const basePullRequestAPI = `${baseAPI}/pulls`;

module.exports.callback = (message, words) => {
    const formatPR = (pr, short = true) => {
        const idTitleAuthor = `**#${pr.number}** ${pr.title} par **${pr.user.login}**`;
        if (short) {
            return `${idTitleAuthor}, <${pr.html_url}>`;
        }
        return `${idTitleAuthor} ${BooleanEmoji(pr.mergeable_state === 'clean')}\n` +
               `*${pr.body}*\n` +
               `<${pr.html_url}>`;
    }

    const handlePRs = () => {
        if (words.length <= 2) {
            HttpsGetJson(basePullRequestAPI, (json) => {
                let reply = `Il y a actuellement ${json.length} pull request${json.length > 1 ? 's' : ''} en attente${json.length > 0 ? ' :' : ''}`;
                for (let item of json)
                    reply += `\n\t- ${formatPR(item)}`
                message.channel.send(reply);
            });
        } else {
            handlePR();
        }
    };

    const handlePR = () => {
        if (words.length > 2) {
            for (let num of words.slice(2).map(x => parseInt(x))) {
                if (!isNaN(num)) {
                    HttpsGetJson(`${basePullRequestAPI}/${num}`, (json) => {
                        // Handle not found
                        if (json.message != null && json.message === 'Not Found') {
                        } else {
                            message.channel.send(formatPR(json, false));
                        }
                    });
                }
            }
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
