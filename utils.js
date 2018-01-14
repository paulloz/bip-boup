module.exports.EmojiOrNothing = (channel, emoji) => channel.guild != null ? channel.guild.emojis.find('name', emoji) || '' : ''
