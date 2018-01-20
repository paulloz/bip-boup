if (process.argv.length < 3) process.exit();

const Command = require(`../commands/${process.argv[2]}.js`);
const MessageMock = {
    reply : console.log,
    channel : { send : console.log },
    react : (str) => console.log(`-> ${str}`)
};

const callCommand = (words) => Command.callback(MessageMock, [process.argv[2]].concat(words));

if (process.argv.length > 3) {
    callCommand(process.argv.slice(3).filter(str => str.length > 0));
} else {
    const Readline = require('readline');

    const rl = Readline.createInterface({
        input: process.stdin,
        output: process.stdout
    });

    rl.question('> ', answer => {
        callCommand(answer.split(/\s+/).filter(str => str.length > 0));
        rl.close();
    });
}
