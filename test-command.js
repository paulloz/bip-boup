if (process.argv.length < 3) process.exit();

const Readline = require('readline');
const Command = require(`./commands/${process.argv[2]}.js`);

const rl = Readline.createInterface({
    input: process.stdin,
    output: process.stdout
});

rl.question('> ', answer => {
    Command.callback({
        reply : console.log,
        channel : { send : console.log }
    }, answer.split(/\s+/).filter(str => str.length > 0));
    rl.close();
});
