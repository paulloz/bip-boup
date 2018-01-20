const Fs = require('fs');
const Chrono = require('chrono-node');
const Moment = require('moment');
const Uuidv1 = require('uuid/v1');

module.exports.command = 'rappel';
module.exports.help = 'Renvoie un message après un temps défini.';

const storage = '.queue.json';

const addOnStorage = (job) => {
    removeFromStorage(job.uuid); // Make sure there's no dups
    Fs.readFile(storage, 'utf8', (err, data) => {
        if (err == null || err.code === 'ENOENT') {
            job.sendAt = job.sendAt.valueOf();
            Fs.writeFileSync(storage, JSON.stringify(JSON.parse(data || '[]').concat([job])), 'utf8');
        }
    });
};

const removeFromStorage = (uuid) => {
    Fs.readFile(storage, 'utf8', (err, data) => {
        if (err == null) {
            try {
                Fs.writeFileSync(storage, JSON.stringify(JSON.parse(data).filter(x => x.uuid != uuid && moment(x.sendAt) >= moment())), 'utf8');
            } catch (e) {
                if (e instanceof SyntaxError)
                    Fs.unlinkSync(storage); // If there's an error inside the file, delete it
            }
        }
    });
};

const withStorageContent = (callback) => {
    Fs.readFile(storage, 'utf8', (err, data) => {
        if (err == null || err.code === 'ENOENT') callback(JSON.parse(data || '[]'));
    });
};

const parse = (toParse, createdTimeStamp) => {
    // TODO Do some smart things to clean the message
    const getMessageText = (full, sIdx, eIdx) => full.substring(...((sIdx > full.length - eIdx) ? [0, sIdx] : [eIdx])).trim();

    const getChronoMatches = (toParse, sentAt) => Chrono.parse(toParse).filter(
        m => Object.keys(m.tags).filter(k => k.startsWith('FR')).length > 0
    ).map(m => Object({
        sendAt: Moment(Chrono.parseDate(m.text, sentAt)),
        text: getMessageText(toParse, m.index, m.index + m.text.length)
    })).filter(m => m.sendAt > sentAt && m.text.length > 0);

    return getChronoMatches(toParse, Moment(createdTimeStamp)).shift();
};

const sendOnChannelAndRemoveFromDisk = (job, channel) => {
    if (channel == null) return;
    setTimeout(() => {
        removeFromStorage(job.uuid);
        channel.send(job.text);
    }, job.sendAt - Moment());
};

module.exports.setup = true;
module.exports.callback = (bipboup) => {
    withStorageContent((jobs) => {
        jobs.map(j => { j.sendAt = Moment(j.sendAt); return j; }).filter(j => j.sendAt >= Moment()).forEach((j) => {
            sendOnChannelAndRemoveFromDisk(j, bipboup.channels.get(j.channel));
        });
    });

    return (message, words) => {
        const sendAndSaveToDisk = (job) => {
            if (job == null) return false;

            job.uuid = job.uuid || Uuidv1();
            job.channel = message.channel.id;
            addOnStorage(job);
            sendOnChannelAndRemoveFromDisk(job, message.channel)
            return true;
        };

        message.react(sendAndSaveToDisk(parse(words.slice(1).join(' '), message.createdTimeStamp)) ? '👌' : '👎');
    };
};
