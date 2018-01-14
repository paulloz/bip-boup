const { EmojiOrNothing } = require('../utils.js');

module.exports.command  = 'note'
module.exports.help     = 'Donne une note imparatiale sur un objet culturel.'
module.exports.callback = (message, words) => {
    const notes = [
        "bnm nrv", "bnm neutre", "bnm content",
        `6/10 très bon ${EmojiOrNothing(message.channel, 'dreamburger')}`,
        "FFXV GOTY", "C DE LA MERDE",
        `pour les fans de ${EmojiOrNothing(message.channel, 'jul')}`,
        "bnm content", ["8/10", "with rice 10/10"], "TÛT TÛT LES JALOUX", "Telerama a aimé, moi non",
        "Ça vaut vraiment pas le premier...", "chef d'œuvre", "vous feriez mieux de regarder Mad Max",
        "quelques rares défauts mais un excellent moment",
        `${EmojiOrNothing(message.channel, 'baby')} qui pleure`,
        "RT car c triss", "j'en sais rien, je connais pas", "bon somnifère", "mieux que du café"
    ];

    const prefix = words.length >= 2 ? words.slice(1).join(' ') + ": " : "";
    const n = Math.floor(Math.random() * Math.floor(notes.length));

    if (typeof(notes[n]) === typeof([])) {
        notes[n].forEach(note => message.channel.send(prefix + note));
    } else {
        message.channel.send(prefix + notes[n]);
    }
};
