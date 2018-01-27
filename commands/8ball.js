module.exports.command  = '8ball'
module.exports.help     = 'Ton meilleur conseiller.'
module.exports.callback = (message, words) => {
    const fortunes = [
        "Essaye plus tard","Essaye encore","Pas d'avis","C'est ton destin","Le sort en est jeté",
        "Une chance sur deux","Repose ta question","D'après moi oui","C'est certain","Oui absolument",
        "Tu peux compter dessus","Sans aucun doute","Très probable","Oui","C'est bien parti","C'est non",
        "Peu probable","Faut pas rêver","N'y compte pas","Impossible"
    ];
    const gif = [
        "http://gph.is/2jWTPPf","http://gph.is/2iZqHD9","http://gph.is/2tdjYNH","http://gph.is/2sI6Igk",
        "http://gph.is/2mUXGgg","http://gph.is/2mRFnW0","http://gph.is/2pGyYTO"
    ];

    message.channel.send(gif[Math.floor(Math.random() * Math.floor(gif.length))]).then(thisMessage => {
        setTimeout(function() {
            thisMessage.delete();
            message.reply(fortunes[Math.floor(Math.random() * Math.floor(fortunes.length))]);
        }, 7000);
    });
};
