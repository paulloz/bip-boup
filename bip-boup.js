const Fs = require('fs');
const Path = require('path');
const Discord = require('discord.js');

const startup = () => {
    console.log('Starting up...');

    const bipboup = new Discord.Client();

    let commands = [];

    bipboup.on('ready', () => {
        console.log('I\'m connected to the Discord guild!');

        // Properly close connection on Ctrl-C
        process.on('SIGINT', () => {
            console.log('Shuting down...');
            bipboup.destroy().then(() => process.exit());
        });
    });

    require(Path.join(__dirname, 'on-message.js'))(bipboup, commands);

    Fs.readdir(Path.join(__dirname, 'commands'), (err, files) => {
        if (err == null) {
            files.forEach(file => {
                if (Path.extname(file) == '.js') {
                    const {command, help, callback} = require(Path.join(__dirname, 'commands', file));
                    if (command == null || help == null || callback == null) return;
                    // TODO Make sure there's no duplicates
                    commands.push({
                        command: command,
                        help: help,
                        callback: callback
                    });
                }
            });
        }
    });

    // Connect from the token found in the .token file
    Fs.readFile('.token', { encoding : 'utf-8' }, (err, data) => {
        if (err == null)
            bipboup.login(data.trimRight());
        else
            console.error(err.message);
    });
};

startup();
