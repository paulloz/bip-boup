const Fs = require('fs');
const Path = require('path');
const Discord = require('discord.js');

const bipboup = new Discord.Client();

// Define how people will get the bot's attention
const attentionChar = '!';
// TODO Change this to enable commands after a mention
const attentionRegexp = new RegExp(`^${attentionChar}(.*)`);

let commands = [];

bipboup.on('ready', () => {
    console.log('I\'m connected to the Discord guild!');

    // Properly close connection on Ctrl-C
    process.on('SIGINT', () => {
        console.log('Shuting down...');
        bipboup.destroy().then(() => process.exit());
    });
});

// Respond to messages
bipboup.on('message', message => {
    // Make sure we de not reply to our own messages
    if (message.author.id == bipboup.user.id) return;

    let messageContent = message.cleanContent.trim().match(attentionRegexp);

    if (messageContent != null) {
        // Split message content into words
        messageContent = messageContent[1].split(/\s+/).filter(word => word.length > 0);
        if (messageContent.length <= 0) return; // Safety

        if (messageContent[0] == 'help') {
            help(message, messageContent);
        } else {
            for (let command of commands) {
                if (messageContent[0] == command.command) {
                    command.callback(message, messageContent);
                    break;
                }
            }
        }
    }
});

// TODO Move this in a command file (but we need access to the `commands` array)
const help = (message, words) => {
    // TODO check words to see if we need only part of the help text
    // Build the help text
    let helpText = '';
    for (let command of commands) {
        helpText += `\n!${command.command} : ${command.help}`;
    }
    message.reply(helpText);
};

const startup = () => {
    console.log('Starting up...');

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
            bipboup.login(data.trimRight('\n'));
        else
            console.error(err.message);
    });
};

startup();
