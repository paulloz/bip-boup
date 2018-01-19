const { HttpsGetJson, BooleanEmoji } = require('../utils.js');

module.exports.command  = 'github'
module.exports.help     = 'Quelques informations sur mon code source.'

const baseAPI = "https://api.github.com/repos/paulloz/bip-boup";
const basePullRequestAPI = `${baseAPI}/pulls`;
const baseIssuesAPI = `${baseAPI}/issues`;

module.exports.callback = (message, words) => {

    const handlePRs = () => {
        const formatPR = (pr, short = true) => {
            const idTitleAuthor = `**#${pr.number}** ${pr.title} par **${pr.user.login}**`;
            if (short) {
                return `${idTitleAuthor}, <${pr.html_url}>`;
            }
            return `${idTitleAuthor} ${BooleanEmoji(pr.mergeable_state === 'clean')}\n` +
                   `${pr.body}\n` +
                   `<${pr.html_url}>`;
        }

        const handlePR = (num) => {
            HttpsGetJson(`${basePullRequestAPI}/${num}`, (json) => {
                // Handle not found
                if (!(json.message != null && json.message === 'Not Found')) {
                    message.channel.send(formatPR(json, false));
                }
            });
        };

        if (words.length <= 2) {
            HttpsGetJson(basePullRequestAPI, (json) => {
                let reply = `Il y a actuellement ${json.length} pull request${json.length > 1 ? 's' : ''} en attente${json.length > 0 ? ' :' : ''}`;
                for (let item of json)
                    reply += `\n\t- ${formatPR(item)}`
                message.channel.send(reply);
            });
        } else {
            for (let num of words.slice(2).map(x => parseInt(x)).filter(num => !isNaN(num))) {
                handlePR(num);
            }
        }
    };

    const handleIssues = () => {
        const formatIssue = (issue, short = true) => {
            const idTitleLabels = `**#${issue.number}** ${issue.title} ${issue.labels.map(label => `\`${label.name}\``).join(' ')}`
            if (short) {
                return `${idTitleLabels}, <${issue.html_url}>`;
            }
            return `${idTitleLabels}\n` +
                   `${issue.body}\n` +
                   `<${issue.html_url}>`;
        };

        const handleIssue = (num) => {
            HttpsGetJson(`${baseIssuesAPI}/${num}`, (json) => {
                // Handle not found
                if (!(json.message != null && json.message === 'Not Found')) {
                    // Handle PR
                    if (json.pull_request == null) {
                        message.channel.send(formatIssue(json, false));
                    }
                }
            })
        };

        if (words.length <= 2) {
            HttpsGetJson(baseIssuesAPI, (json) => {
                json = json.filter(item => item.pull_request == null);
                let reply = `Il y a actuellement ${json.length} issue${json.length > 1 ? 's' : ''} ouverte${json.length > 1 ? 's :' : json.length > 0 ? ' :' : ''}`;
                for (let item of json)
                    reply += `\n\t- ${formatIssue(item)}`
                message.channel.send(reply);
            });
        } else {
            for (let num of words.slice(2).map(x => parseInt(x)).filter(num => !isNaN(num))) {
                handleIssue(num);
            }
        }
    };

    if (words.length <= 1) {
        message.reply(`<https://github.com/paulloz/bip-boup.git>`);
    } else {
        switch (words[1].toLowerCase()) {
            case "pr": handlePRs(); break;
            case "issue": handleIssues(); break;
            default: break;
        }
    }
};
