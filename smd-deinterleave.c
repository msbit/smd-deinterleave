#include <stdint.h>
#include <stdio.h>

#define SIZE 16384

typedef struct {
  uint8_t size_of_file;    // Byte 00h : Size of file.
  uint8_t file_data_type;  // Byte 01h : File data type.
  uint8_t status_flags_1;  // Byte 02h : Status flags.
  uint8_t status_flags_2;  //
  uint8_t status_flags_3;  //
  uint8_t status_flags_4;  //
  uint8_t status_flags_5;  //
  uint8_t status_flags_6;  //
  uint8_t identifier_1;    // Byte 08h : Identifier 1.
  uint8_t identifier_2;    // Byte 09h : Identifier 2.
  uint8_t file_type;       // Byte 0Ah : File type.
} smd_header_t;

int main(int argc, char **argv) {
  int result = 0;
  if (argc != 3) {
    fprintf(stderr, "usage: %s <input-file> <output-file>\n", argv[0]);
    result = -1;
    goto defer_none;
  }

  uint8_t inputBuffer[SIZE];
  uint8_t outputBuffer[SIZE];

  FILE *input = fopen(argv[1], "r");
  if (input == NULL) {
    perror("fopen input");
    result = -1;
    goto defer_none;
  }

  smd_header_t header;
  if (fread(&header, sizeof(header), 1, input) < 1) {
    perror("fread header");
    result = -1;
    goto defer_close_input;
  }

  if (header.identifier_1 != 0xaa || header.identifier_2 != 0xbb || header.file_type != 0x06) {
    fprintf(
      stderr,
      "invalid file: identifier_1 = 0x%02x identifier_2 = 0x%02x file_type = 0x%02x\n",
      header.identifier_1,
      header.identifier_2,
      header.file_type
    );
    result = -1;
    goto defer_close_input;
  }

  if (fseek(input, 512, SEEK_SET) == -1) {
    perror("fseek");
    goto defer_close_input;
  }

  FILE *output = fopen(argv[2], "w");
  if (output == NULL) {
    perror("fopen output");
    result = -1;
    goto defer_close_input;
  }

  for (int i = 0; i < header.size_of_file; i++) {
    if (fread(inputBuffer, 1, SIZE, input) < 1) {
      perror("fread chunk");
      result = -1;
      goto defer_close_output;
    }

    for (int j = 0, k = 0; j < (SIZE / 2); j++, k += 2) {
      outputBuffer[k] = inputBuffer[j + (SIZE / 2)];
      outputBuffer[k + 1] = inputBuffer[j];
    }

    if (fwrite(outputBuffer, 1, SIZE, output) < 1) {
      perror("fwrite chunk");
      result = -1;
      goto defer_close_output;
    }
  }

defer_close_output:
  fclose(output);

defer_close_input:
  fclose(input);

defer_none:
  return result;
}
