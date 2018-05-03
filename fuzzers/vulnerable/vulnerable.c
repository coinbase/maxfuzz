// vulnerable.c taken from https://github.com/ThalesIgnite/afl-training under
// the MIT License, copied below from the repository:
//
// Copyright 2016, 2017 Thales UK Ltd
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

#include <string.h>
#include <stdio.h>
#include <unistd.h>
#include <stdlib.h>

int main(int argc, char *argv[])
{
	char input[100] = {0};
	char *out;

	// Slurp input
	if (read(STDIN_FILENO, input, 100) < 0) {
		fprintf(stderr, "Couldn't read stdin.\n");
	}
	if(input[0] == 'c') {
		// count characters
		out = malloc(sizeof(input) - 1 + 3); // enough space for 2 digits + a space + input-1 chars
		sprintf(out, "%lu ", strlen(input) - 1);
		strcat(out, input+1);
		printf("%s", out);
		free(out);
	} else if ((input[0] == 'e') && (input[1] == 'c')) {
		// echo input
		printf(input + 2);
	} else if (strncmp(input, "head", 4) == 0) {
		// head
		if (strlen(input) > 5) {
			input[input[4]] = '\0'; // truncate string at specified position
			printf("%s", input+4);
		} else {
			fprintf(stderr, "head input was too small\n");
		}
	} else if (strcmp(input, "surprise!\n") == 0) {
		// easter egg!
		*(char *)1=2;
	} else {
		fprintf(stderr, "Usage: %s\nText utility - accepts commands on stdin and prints results to stdout:\n", argv[0]);
		fprintf(stderr, "\tInput           | Output\n");
		fprintf(stderr, "\t----------------+-----------------------\n");
		fprintf(stderr, "\tec<string>      | <string> (simple echo)\n");
		fprintf(stderr, "\thead<N><string> | The first <N> bytes of <string>\n");
		fprintf(stderr, "\tc<string>       | The length of <string>, followed by <string>\n");
	}
}
