const Fs = require('fs');

class Config {
    constructor() {
        this.defaults = {
            version: 2,
            attention: '!',
            mainchan: null
        };

        this.fileName = 'config.json';
        this.data = { };
        this.commands = [ ];

        this.readFromDisk();
    }

    writeDefaults(only) {
        only = only == null ? true : only;

        let data = { general : this.defaults };

        if (!only) {
            data = JSON.parse(Fs.readFileSync(this.fileName, 'utf8'))
            data.general = this.defaults;
        }

        Fs.writeFileSync(this.fileName, JSON.stringify(data), 'utf8');
        this.readFromDisk();
    }

    readFromDisk() {
        try {
            let data = JSON.parse(Fs.readFileSync(this.fileName, 'utf8'));
            if (data.general.version < this.defaults.version)
                this.writeDefaults(false);
            else
                Object.entries(data).forEach(entry => {
                    this.data[entry[0]] = Object.entries(entry[1]).filter(entry => {
                        return Object.keys(this.defaults).indexOf(entry[0]) >= 0;
                    }).reduce((acc, entry) => { acc[entry[0]] = entry[1]; return acc; }, { });
                });
        } catch(e) {
            if (e.code === 'ENOENT' || e.name === 'SyntaxError')
                writeDefaults();
        }
    }

    get(key, guild) {
        return this.data[(guild || {}).id || 'general'][key];
    }
}

const CONFIG_KEY = Symbol.for('bipboup.config');

let globalSymbols = Object.getOwnPropertySymbols(global);
let hasConfig = (globalSymbols.indexOf(CONFIG_KEY) > -1);

if (!hasConfig) {
    global[CONFIG_KEY] = new Config();
}

let singleton = { };

Object.defineProperty(singleton, 'instance', {
    get: function() {
        return global[CONFIG_KEY];
    }
});

Object.freeze(singleton);

module.exports = singleton.instance;
