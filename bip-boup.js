const Fs = require('fs');
const Discord = require('discord.js');

const bipboup = new Discord.Client();

// Define how people will get the bot's attention
const attentionChar = '!';
const attentionRegexp = new RegExp(`^${attentionChar}(.*)`);

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
    let messageContent = message.cleanContent.trim().match(attentionRegexp);

    if (messageContent != null) {
        // Split message content into words
        messageContent = messageContent[1].split(/\s+/).filter(word => word.length > 0);
        if (messageContent.length <= 0) return; // Safety

        // TODO Define commands properly in other files
        if (messageContent[0] == 'github') {
            message.reply(`<https://github.com/paulloz/bip-boup.git>`);
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
