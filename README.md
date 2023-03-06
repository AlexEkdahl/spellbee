# SpellBee

This is a simple Go program that checks the correct article (a or an) to use before a given word. It uses the [TextGears API](https://textgears.com/api/) to perform grammar checking.

## Usage

To run the program, you need to have a TextGears API key and set it as an environment variable named `TEXT_GEARS_API_KEY`. You can get a free API key from [here](https://textgears.com/signup.php).

Then, you can run the program with a word as an argument, for example:

```bash
go run main.go apple
```

The program will output the correct article to use before the word, for example:

The word 'apple' should be preceded by 'an'.
If you don’t provide a word as an argument, the program will print an error message and exit.

You can also initialize the database used by the program by running it with the `--init` flag:

```bash
spellbee --init
```
This will create the database and the necessary table for caching results.

## Compilation
To compile the program, you can use the Makefile provided in the repository. Just run:

```bash
make build
```

This will create a binary file named sp in the bin folder. You can then run the binary file with a word as an argument, for example:
```bash
./bin/sp apple
```
## implementation

To see how spellbee is implemented in neovim, you can check out this file:

[spellbee.lua](https://github.com/AlexEkdahl/.dotfiles/blob/main/nvim/lua/spellbee.lua)

## Limitations
The program only works for single words and does not handle phrases or sentences. It also assumes that the word is in English and uses the British spelling. It may not be accurate for words that have different pronunciations or articles depending on the context.
