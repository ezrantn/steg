# LSB Steganography in PNG Files

## Table of Contents
- [LSB Steganography in PNG Files](#lsb-steganography-in-png-files)
  - [Table of Contents](#table-of-contents)
  - [Introduction](#introduction)
  - [Project Overview](#project-overview)
  - [Steganography and LSB Technique](#steganography-and-lsb-technique)
    - [Traditional LSB in Images:](#traditional-lsb-in-images)
    - [Adaptation for PNG Chunks:](#adaptation-for-png-chunks)
  - [PNG File Structure](#png-file-structure)
  - [Implementation Details](#implementation-details)
  - [Advantages and Limitations](#advantages-and-limitations)
  - [Future Enhancements](#future-enhancements)
  - [Conclusion](#conclusion)

## Introduction

This document outlines a steganography project that focuses on hiding secret messages within PNG image files. The project utilizes the Least Significant Bit (LSB) technique, adapted to work with PNG file structures, to embed and extract hidden information.

## Project Overview

The goal of this project is to develop a web-based steganography tool that allows users to:
1. Upload a PNG image
2. Provide a secret message
3. Receive a modified PNG image with the message hidden inside
4. Extract hidden messages from modified PNG images

The project is implemented in Go and is designed to be integrated into a web application, making it accessible and user-friendly.


## Usage

The program allows encoding (hiding) a message in a PNG file and decoding (retrieving) a hidden message from a modified PNG file. Follow the instructions below to use the encoding and decoding features:

### Encoding a Message

To hide a message in a PNG file, use the `--encode` option along with the input image, output file, secret key, and offset for added security.

```shell
go run . -i <input_image.png> -o <output_image.png> --key <secret_key> --offset <integer_offset> --encode
```

Parameters:

- `-i <input_image.png>`: The path to the original PNG file where the message will be hidden.
- `-o <output_image.png>`: The output PNG file with the hidden message.
- `--key <secret_key>`: A string key used to encode the message.
- `--offset <integer_offset>`: An integer offset to add randomness in the message's position.

#### Example

```shell
go run . -i original.png -o output.png --key gophersValid --offset 1337 --encode
```

### Decoding a Message

To retrieve a hidden message from a PNG file, use the --decode option along with the input PNG, output file, secret key, and offset used during encoding.

```shell
go run . -i <input_image.png> -o <decoded_output.png> --key <secret_key> --offset <integer_offset> --decode
```

Parameters:

- `-i <input_image.png>`: The path to the PNG file containing the hidden message.
- `-o <decoded_output.png>`: The output file where the extracted message will be saved.
- `--key <secret_key>`: The same key used during encoding to retrieve the message.
- `--offset <integer_offset>`: The same offset used during encoding.

#### Example

```shell
go run . -i output.png -o decoded_output.png --key gophersValid --offset 1337 --decode
```

### Viewing the Decoded Message

After decoding, the program will print the hidden message in the terminal as `Payload Decode`. Alternatively, open the `decoded_output.png` to see the retrieved message if it’s stored in a separate text format.

## Steganography and LSB Technique

Steganography is the practice of concealing information within other non-secret data to avoid detection. The Least Significant Bit (LSB) technique is a common method used in image steganography.

### Traditional LSB in Images

1. In a typical image, each pixel is represented by color values (e.g., RGB).
2. Each color value is usually stored as an 8-bit integer (0-255).
3. The LSB technique replaces the least significant bit of these values with bits from the secret message.
4. This causes minimal change to the image, often imperceptible to the human eye.

Example:
- Original pixel: (138, 255, 74)
- Binary: (10001010, 11111111, 01001010)
- Hiding '101':
- Modified: (10001011, 11111111, 01001011)

New pixel: (139, 255, 75)

### Adaptation for PNG Chunks

My project adapts the LSB concept to work with PNG file structures:

1. Instead of modifying pixel data directly, we manipulate PNG chunks.
2. We create new chunks or modify existing ones to store our hidden data.
3. This approach maintains the integrity of the image data while still hiding information within the file structure.

## PNG File Structure

PNG files are structured as a series of chunks:

1. PNG Signature: 8 bytes that always begin a PNG file.
2. IHDR chunk: Contains basic information about the image.
3. Various other chunks (PLTE, IDAT, etc.): Contain palette, image data, and other information.
4. IEND chunk: Marks the end of the PNG datastream.

Each chunk consists of:
- Length (4 bytes)
- Chunk Type (4 bytes)
- Chunk Data (variable length)
- CRC (Cyclic Redundancy Check, 4 bytes)

My steganography method involves creating or modifying chunks to store hidden data while maintaining a valid PNG structure.

## Implementation Details

The project is implemented in Go and consists of several key components:

1. PNG Parsing:
   - Reading and validating PNG headers
   - Parsing individual chunks

2. Steganography Engine:
   - Embedding messages by creating new chunks or modifying existing ones
   - Extracting hidden messages from chunks

3. Web Server:
   - Handling file uploads and downloads
   - Processing steganography requests

4. User Interface:
   - Web-based interface for easy interaction with the tool

Key functions include:
- `encodeMessage`: Embeds a secret message into a PNG file
- `decodeMessage`: Extracts a hidden message from a PNG file
- `processImage`: Analyzes and modifies PNG chunks
- `writeData`: Writes modified PNG data back to a file or stream

## Advantages and Limitations

Advantages:
1. Maintains image quality as pixel data is not directly modified
2. Difficult to detect without knowing the specific implementation
3. Can potentially store larger amounts of data compared to traditional LSB

Limitations:
1. Limited to PNG file format
2. Potential increase in file size due to additional chunks
3. Vulnerable if the method of hiding data becomes known

## Future Enhancements

1. Implement encryption for the hidden messages
2. Extend support to other image formats
3. Develop techniques to make the steganography more resilient to image modifications
4. Implement a more sophisticated web interface with real-time processing

## Conclusion

This project demonstrates an innovative approach to steganography by adapting the LSB technique to work with PNG file structures. It provides a practical tool for hiding and extracting secret messages within PNG images, with potential applications in secure communication and data privacy.
