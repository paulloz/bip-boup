const { HttpsGetJson, BooleanEmoji, Plural } = require('../utils.js');
const Discord = require('discord.js');

module.exports.command  = 'github'
module.exports.help     = 'Quelques informations sur mon code source.'

const baseAPI = "https://api.github.com/repos/paulloz/bip-boup";
const basePullRequestAPI = `${baseAPI}/pulls`;
const baseIssuesAPI = `${baseAPI}/issues`;

module.exports.callback = (message, words) => {
    const embed = (json) => {
        return new Discord.RichEmbed({
            author: {
                name: json.user.login,
                icon_url: json.user.avatar_url
            },
            title: `#${json.number} ${json.title}`,
            url: json.html_url,
            description: json.body
        });
    };

    const handlePRs = () => {
        const handlePR = (num) => {
            HttpsGetJson(`${basePullRequestAPI}/${num}`, (json) => {
                // Handle not found
                if (!(json.message != null && json.message === 'Not Found')) {
                    message.channel.send(embed(json).setColor(
                        json.mergeable_state == 'clean' ? 'GREEN' : 'RED'
                    ).setFooter(`${json.commits} ${Plural('commit', json.commits)} | ${json.changed_files} ${Plural('fichier', json.changed_files)} | +${json.additions} -${json.deletions}`));
                }
            });
        };

        if (words.length <= 2) {
            HttpsGetJson(basePullRequestAPI, (json) => {
                let embed = new Discord.RichEmbed();
                for (let item of json)
                    embed.addField(`#${item.number} ${item.title} par *${item.user.login}*`, item.html_url);
                message.channel.send(embed.setFooter(`${json.length} pull ${Plural('request', json)} actuellement en attente.`));
            });
        } else {
            for (let num of words.slice(2).map(x => parseInt(x)).filter(num => !isNaN(num))) {
                handlePR(num);
            }
        }
    };

    const handleIssues = () => {
        const handleIssue = (num) => {
            HttpsGetJson(`${baseIssuesAPI}/${num}`, (json) => {
                // Handle not found and PR
                if (!(json.message != null && json.message === 'Not Found') && json.pull_request == null) {
                    console.log(json);
                    message.channel.send(embed(json));
                }
            })
        };

        if (words.length <= 2) {
            HttpsGetJson(baseIssuesAPI, (json) => {
                json = json.filter(item => item.pull_request == null);
                let embed = new Discord.RichEmbed();
                console.log(embed.fields.length);
                for (let item of json)
                    embed.addField(`#${item.number} ${item.title}`, item.html_url);
                message.channel.send(embed.setFooter(`${json.length} ${Plural('ticket', json)} actuellement ${Plural('ouvert', json)}.`));
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
            case "issue":
            case "issues":
            case "ticket":
            case "tickets": handleIssues(); break;
            default: break;
        }
    }
};
