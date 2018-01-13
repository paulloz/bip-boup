module.exports.command  = 'note'
module.exports.help     = 'Donne une note imparatiale sur un objet culturel.'
module.exports.callback = (message, words) => {
    var msg = "";
    switch(Math.floor(Math.random() * Math.floor(5))) {
        case 0: msg = "bnm nrv"; break;
        case 1: msg = "bnm neutre"; break;
        case 2: msg = "6/10 très bon"; break;
        case 3: msg = "FFXV GOTY"; break;
        case 4: msg = "C DE LA MERDE"; break;
    }
    message.reply(words[1] + "msg");
};
