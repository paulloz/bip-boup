module.exports.command  = 'note'
module.exports.help     = 'Donne une note imparatiale sur un objet culturel.'
module.exports.callback = (message, words) => {
    const notes = ["bnm nrv", "bnm neutre", "6/10 très bon :dreambuger:", "FFXV GOTY", "C DE LA MERDE"];
    const prefix = words.length >= 2 ? words.slice(1).join(' ') + ": " : "";
    message.channel.send(prefix + notes[Math.floor(Math.random() * Math.floor(notes.length))]);
};
