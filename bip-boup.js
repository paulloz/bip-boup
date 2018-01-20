const Fs = require('fs');
const Path = require('path');
const Winston = require('winston');
const Discord = require('discord.js');

const startup = () => {
    Winston.log('info', 'Starting up...');

    const bipboup = new Discord.Client();

    bipboup.config = {
        attentionChar : '!', // Define how people will get the bot's attention
        commands : []
    };

    bipboup.on('ready', () => {
        Winston.log('info', 'I\'m connected to the Discord network!');

        // Properly close connection on Ctrl-C
        process.on('SIGINT', () => {
            Winston.log('info', 'Shutting down...');
            bipboup.destroy().then(() => process.exit());
        });
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

    Fs.readdir(Path.join(__dirname, 'commands'), (err, files) => {
        if (err == null) {
            files.forEach(file => {
                if (Path.extname(file) == '.js') {
                    const {command, help, callback, setup} = require(Path.join(__dirname, 'commands', file));
                    if (command == null || help == null || callback == null) return;

                    if (bipboup.config.commands.find(c => c.command === command) == null) {
                        bipboup.config.commands.push({
                            command: command,
                            help: help,
                            callback: callback
                        });
                    }

                    if (setup != null)
                        setup(bipboup.config);
                }
            });
        }
    });
};

startup();
