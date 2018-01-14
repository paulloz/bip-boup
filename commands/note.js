module.exports.command  = 'note'
module.exports.help     = 'Donne une note imparatiale sur un objet culturel.'
module.exports.callback = (message, words) => {
    const notes = ["bnm nrv", "bnm neutre", "6/10 très bon", "FFXV GOTY", "C DE LA MERDE"];
    if (words.length < 2) return;
    message.channel.send(words.slice(1).join(' ') + ": " + notes[Math.floor(Math.random() * Math.floor(notes.length))]);
};
