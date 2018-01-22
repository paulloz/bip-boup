const Fs = require('fs');
const Path = require('path');
const Winston = require('winston');
const Discord = require('discord.js');
const Config = require('./config.js')

const startup = () => {
    Winston.log('info', 'Starting up...');

    const bipboup = new Discord.Client();
    let isInit = false;

    bipboup.config = {
        commands : []
    };

    const init = () => {
        Fs.readdir(Path.join(__dirname, 'commands'), (err, files) => {
            if (err == null) {
                files.forEach(file => {
                    if (Path.extname(file) == '.js') {
                        const {command, help, callback, setup} = require(Path.join(__dirname, 'commands', file));
                        if (command == null || help == null || callback == null) return;

                        if (!Config.hasCommand(command))
                            Config.addCommand(command, help, setup ? callback(bipboup) : callback);
                    }
                });
            }
        });

        isInit = true;
    };

    bipboup.on('ready', () => {
        Winston.log('info', 'I\'m connected to the Discord network!');

        // Properly close connection on Ctrl-C
        process.on('SIGINT', () => {
            Winston.log('info', 'Shutting down...');
            bipboup.destroy().then(() => process.exit());
        });

        if (!isInit)
            init();
    });

    require(Path.join(__dirname, 'on-message.js'))(bipboup);
    require(Path.join(__dirname, 'on-emoji.js'))(bipboup);

    // Connect from the token found in the .token file
    Fs.readFile('.token', { encoding : 'utf-8' }, (err, data) => {
        if (err == null)
            bipboup.login(data.trimRight());
        else
            Winston.log('error', err.message);
    });
};

startup();
