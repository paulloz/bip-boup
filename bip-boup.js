const Fs = require('fs');

const Discord = require('discord.js');
const bipboup = new Discord.Client();

bipboup.on('ready', () => {
    console.log('READY!');
});

Fs.readFile('.tosken', { encoding : 'utf-8' }, (err, data) => {
    if (err == null)
        bipboup.login(data.trimRight('\n'));
    else
        console.error(err.message);
});
