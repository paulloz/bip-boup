module.exports.command  = '8ball'
module.exports.help     = 'Ton meilleur conseiller.'
module.exports.callback = (message, words) => {
    const fortunes = [
        "Essaye plus tard","Essaye encore","Pas d'avis","C'est ton destin","Le sort en est jeté",
        "Une chance sur deux","Repose ta question","D'après moi oui","C'est certain","Oui absolument",
        "Tu peux compter dessus","Sans aucun doute","Très probable","Oui","C'est bien parti","C'est non",
        "Peu probable","Faut pas rêver","N'y compte pas","Impossible"
    ];

    let n = 5;

    for (var i = 0; i < 10; i++) {
        roll(i*500,n);
        n += Math.floor(Math.random() * Math.floor(5)) - 2;
    }

    function roll(wt, n) {
        setTimeout(function() {
            if (wt >= 4500) {
                setTimeout(function() {
                    message.reply(fortunes[Math.floor(Math.random() * Math.floor(fortunes.length))]);
                    return;
                },1000);
            }
            let str = ".";
            str+= " ".repeat(Math.abs(n));
            message.channel.send(str + ":right_facing_fist::8ball::left_facing_fist:");
        }, wt);
    }
};
