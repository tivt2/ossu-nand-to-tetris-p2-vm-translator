VM translator - nand to tetris part 2 - week 1 and 2.
This are my results when developing the course project

VM translator is a program that reads a given text file, which contains the course VM commands,
it produces from it another text file that contains the course assembly commands,
it operates line by line, reading and parsing, and then write a single file in the end.

Giving a "file.vm" to be translated, it will generate a "file.asm" in the same directory of the provided file,
else if a "folder" is provided, it will translate all "file.vm" inside the folder to a "folder.asm" file,
where it will have a initial section for bootstraping the program.
