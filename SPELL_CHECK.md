# Make a Spell Checker!

Write a program that checks spelling. The input to the program is a dictionary file containing a list of valid words and a file containing the text to be checked.
You can use the `dictionary.txt` file included here as your dictionary.

The program should run on the command line like so:

```sh
my-cool-spellchecker dictionary.txt file-to-check.txt
# output here
```

## The Features

Your program should support the following features (time permitting):

- The program outputs a list of incorrectly spelled words.
- For each misspelled word, the program outputs a list of suggested words.
- The program includes the line and column number of the misspelled word
- The program prints the misspelled word along with some surrounding context.
- The program handles proper nouns (person or place names, for example) correctly.
