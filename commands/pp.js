const Https = require('https')
const { URL } = require('url');
const Stream = require('stream');
const Fs = require('fs');
const Path = require('path');

module.exports.command = 'pp';
module.exports.help = '';
module.exports.admin = true;
module.exports.setup = true;
module.exports.callback = (bipboup) => {
    const ppdir =  Path.join(__dirname, '..', 'pp');
    if (!Fs.existsSync(ppdir))
        Fs.mkdirSync(ppdir);

    return (message, words) => {
        const upload = () => {
            let ok = false;
            message.attachments.array().forEach(attachment => {
                if (ok) return;
                if (attachment.width != null && attachment.width === attachment.height) {
                    url = new URL(attachment.proxyURL);
                    Https.get({
                        protocol: url.protocol,
                        hostname: url.hostname,
                        path: url.pathname + url.search,
                        headers: {
                            'User-Agent': 'Bip Boup/1.0.0'
                        }
                    }, result => {
                        let data = new Stream.Transform();

                        result.on('data', chunk => data.push(chunk));

                        result.on('end', () => {
                            let filename = attachment.filename;
                            let ext = Path.extname(filename);
                            if (words.length > 1) filename = words[1] + ext;
                            let path = Path.join(ppdir, filename);
                            while (Fs.existsSync(path)) {
                                filename = Path.basename(filename, ext).split('-');
                                filename = [filename[0]].concat([String(parseInt(filename[1] || 0) + 1)]).join('-');
                                path = Path.join(ppdir, filename + ext);
                            }

                            Fs.writeFileSync(path, data.read(), { encoding: 'binary', flag: 'w+' });

                            setPP(path);

                            ok = true;
                        });
                    });
                }
            });

            if (ok && message.deletable)
                setTimeout(() => message.delete(), 100);

            return ok;
        };

        const withFile = (filename, callback) => {
            Fs.readdir(ppdir, (err, files) => {
                for (let i = 0; i < files.length; ++i) {
                    if (Path.basename(files[i], Path.extname(files[i])) === filename) {
                        callback(Path.join(ppdir, files[i]));
                        break;
                    }
                }
            });
        };

        const listPP = () => {
            Fs.readdir(ppdir, (err, files) => {
                message.reply(files.map((file) => Path.basename(file, Path.extname(file))).join(', '));
            });
        };

        const rmPP = (f) => {
            if (f != null)
                Fs.unlink(f, (e) => message.react(e == null ? '👌' : '👎'));
            else
                withFile(words[2], rmPP);
        };

        const setPP = (f) => {
            if (f != null)
                bipboup.user.setAvatar(f).then(() => {
                    message.channel.send(':robot:');
                }).catch((error) => {
                    message.reply(':warning: Je ne peux pas changer de PP maintenant !');
                });
            else
                withFile(words[1], setPP);
        };

        if (message.attachments.array().length > 0) {
            upload();
        } else if (words.length >= 2) {
            switch (words[1])
            {
                case 'list':
                    listPP();
                    break;
                case 'rm':
                    rmPP();
                    break;
                default:
                    setPP();
                    break;
            }
        }
    };
};
