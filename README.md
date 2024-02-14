# Instructions

Install the latest Go compiler for your operating system (if you don't have one) from https://go.dev/doc/install

The program doesn't have external dependencies so it can be compiled AS-IS

Compile the program by running `go build -o droplet-spell-checker .` to get executable file `droplet-spell-checker`

Run the program as expected `./droplet-spell-checker <dict-file> <file-to-check>` by passing 2 files

Program has basic error checking but some edge cases are not covered due to limited time constraint. 

## Extra files

Files `dict.txt` and `file.txt` are the test cases used during development. The `dict.txt` is a shrunk down version of the `dictionary.txt`

# Implementation

## Dictionary handling

We first read the dictionary file (first pass) to find out the number of lines -- we count `\n` characters.

Next, we allocate memory buffer to hold that many strings and read the file again while storing individual dictionary entries as individual strings.

Lastly, we check if the buffer is sorted (assumption needed for optimized search)

## Processing of the text file

The given text file is read on per-line basis while each line is seen as a collection of words separated by whitespaces and punctuation marks. We split the line into words and process all words of the same line in one program cycle.

Every word is analyzed for whether or not it needs to be compared against the dictionary. For example capitalized words that are not preceded by punctuation marks are treated as proper nouns and are ignored. 

Each word is lowercased and is attempted to be found in the dictionary.

### Dictionary search 

We use binary search over pre-sorted (by assumption) dictionary array of words. The algorithm is augmented to not just report the failure but to report the most likely matching position by returning the index in dictionary array with the negative sign. In this case, when a word is found in the dictionary -- the positive result corresponds to word index in the dictionary. Otherwise, a negative value, once negated, points to the closest match.

### Search result processing

When the current word is found in the dictionary we move on to the next word in the line.

When the word is not found, it needs to be reported. To report the word we need to identify
- context, i.e. substring within the original line
- list of suggestions
- position of the misspelled word (column) in the line
- line index in the file

While determining context it has to be noted that the same misspelled word can occur multiple time in the given line. So we first obtain a set of words surrounding the misspelled one and combine them together to create a regular expression that can be matched against the line to find the actual context. 

The list of suggestions is compiled from the neighboring entries of the closest matching word in the dictionary.

The line index within the file is the index of the currently processed line, no extra handling is needed.

The column position within the line is computed based on the position of the misspelled word in the "context" increased by the index of the context within the line.

## Result reporting

The result is reported as formatted JSON array where individual misspelled words (entries of the array) are reported with all required attributes.

# Credit

The implementation uses code snippets found in StackOverflow to utilize tried and tested solutions for binary search and efficient counting of lines in a text file. 

Credits are given in form of comment lines (with links to original SO posts) prior corresponding code snippets. 