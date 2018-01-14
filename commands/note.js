const { EmojiOrNothing } = require('../utils.js');

module.exports.command  = 'note'
module.exports.help     = 'Donne une note imparatiale sur un objet culturel.'
module.exports.callback = (message, words) => {
    const notes = [
                   "bnm nrv", "bnm neutre",
                   `6/10 très bon ${EmojiOrNothing(message.channel, 'dreamburger')}`,
                   "FFXV GOTY", "C DE LA MERDE",
                   `pour les fans de ${EmojiOrNothing(message.channel, 'jul')}`,
                   "bnm content", "rice", "TÛT TÛT LES JALOUX", "Telerama a aimé, moi non",
                   "Ça vaut vraiment pas le premier...", "chef d'œuvre", "vous feriez mieux de regarde Mad Max",
                   "quelques rares défauts mais un excellent moment",
                   `${EmojiOrNothing(message.channel, 'baby')} qui pleure`,
                   "RT car c triss", "j'en sais rien, je connais pas", "bon somnifère", "mieux que du café"
               ];
    const prefix = words.length >= 2 ? words.slice(1).join(' ') + ": " : "";
    let seed = Math.floor(Math.random() * Math.floor(notes.length));
    if (notes[seed] === "rice") {
        message.channel.send(prefix + "8/10");
        message.channel.send(prefix + "with rice 10/10");
    } else {
        message.channel.send(prefix + notes[seed]);
    }
};
