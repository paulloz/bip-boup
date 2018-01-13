module.exports.command  = 'note'
module.exports.help     = 'Donne une note imparatiale sur un objet culturel.'
module.exports.callback = (message, words) => {
    const notes = ["bnm nrv", "bnm neutre", "6/10 très bon", "FFXV GOTY", "C DE LA MERDE"];
    message.reply(words.length > 0 ? words[1] : "" + notes[Math.floor(Math.random() * Math.floor(notes.length))]);
};
