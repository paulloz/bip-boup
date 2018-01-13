module.exports.command  = 'github'
module.exports.help     = 'Envoie l\'adresse du dépôt GitHub de mon code source.'
module.exports.callback = (message, words) => {
    message.reply(`<https://github.com/paulloz/bip-boup.git>`);
};
