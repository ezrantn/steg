# Steg

## Overview

Steg is a command-line tool for hiding and retrieving messages within PNG images using Least Significant Bit (LSB) steganography. This program provides encoding and decoding functionalities, allowing users to embed and extract hidden information from PNG files. Future plans for Steg include performance optimizations and support for additional file formats.

## Usage

The program allows encoding (hiding) a message in a PNG file and decoding (retrieving) a hidden message from a modified PNG file. Follow the instructions below to use the encoding and decoding features:

### Encoding a Message

To hide a message in a PNG file, use the `encode` option along with the input image, output file, secret key, and offset for added security.

```shell
go run . -i <input_image.png> -o <output_image.png> --key <secret_key> --offset <integer_offset> --payload <string> encode
```

Parameters:

- `-i <input_image.png>`: The path to the original PNG file where the message will be hidden.
- `-o <output_image.png>`: The output PNG file with the hidden message.
- `--key <secret_key>`: A string key used to encode the message.
- `--offset <integer_offset>`: An integer offset to add randomness in the message's position.
- `--payload <string>`: This is whatever message that you want to hide inside the image.

#### Example

```shell
go run . -i original.png -o output.png --key gophersValid --offset 1337 --payload "Hello there" encode
```

### Decoding a Message

To retrieve a hidden message from a PNG file, use the `decode` option along with the input PNG, output file, secret key, and offset used during encoding.

```shell
go run . -i <input_image.png> -o <decoded_output.png> --key <secret_key> --offset <integer_offset> decode
```

Parameters:

- `-i <input_image.png>`: The path to the PNG file containing the hidden message.
- `-o <decoded_output.png>`: The output file where the extracted message will be saved.
- `--key <secret_key>`: The same key used during encoding to retrieve the message.
- `--offset <integer_offset>`: The same offset used during encoding.

#### Example

```shell
go run . -i output.png -o decoded_output.png --key gophersValid --offset 1337 decode
```

### Viewing the Decoded Message

After decoding, the program will print the hidden message in the terminal as `Decoded Payload`. Alternatively, open the `decoded_output.png` to see the retrieved message if it’s stored in a separate text format.

## Program Structure

Steg is written in Go, focusing on two primary functions:

- **Encoding**: Embeds a message in a PNG file by modifying PNG chunks.
- **Decoding**: Extracts the message from the modified PNG file using the specified key and offset.

## Future Enhancements

- **Performance Optimization**: Improve the tool’s efficiency, especially with large files.
- **Support for Multiple Formats**: Expand support to additional image formats beyond PNG.
- **More Robust Encryption**: Currently, Steg uses XOR to encrypt the payload, but I'm working on implementing a more robust algorithm.

## Warning ⚠️

This CLI is not production-ready yet. Steg still has a long way to go before it can be used as a real steganography tool. This project is an experimental and learning process for me to understand and implement steganography concepts.

However, you’re welcome to try it out and explore its features. If you find a bug or have suggestions for improvement, please feel free to submit a pull request to this repository.