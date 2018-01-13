const Fs = require('fs');
const Discord = require('discord.js');

const bipboup = new Discord.Client();

const attentionChar = '!';

console.log('Starting up...');

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
    let messageContent = message.cleanContent;
    console.log(`${message.author.username} : ${messageContent}`);
    if (messageContent.startsWith(attentionChar)) {
        // Clean message content
        messageContent = messageContent.trim(attentionChar).trim();

        // TODO Define commands properly in other files
        if (messageContent.startsWith('github')) {
            message.reply(`${message.author} <https://github.com/paulloz/bip-boup.git>`);
        }
    }
});

// Connect from the token found in the .token file
Fs.readFile('.token', { encoding : 'utf-8' }, (err, data) => {
    if (err == null)
        bipboup.login(data.trimRight('\n'));
    else
        console.error(err.message);
});
