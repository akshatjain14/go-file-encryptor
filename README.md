# go-file-encryptor

This Go project provides a program to encrypt files in a given source directory, store the encrypted files in a destination directory, and verify the integrity of the encrypted files using SHA-256 hashing. It can also decrypt the encrypted files to a specified location.

## Features
- File Encryption: Encrypts files using AES (Advanced Encryption Standard) in CFB (Cipher Feedback) mode.
- File Decryption: Decrypts files encrypted by the program.
- Integrity Check: Computes SHA-256 hashes of the encrypted files and verifies their integrity before decryption.

## Installation

1. Clone the repository:
git clone https://github.com/akshatjain14/go-file-encryptor.git
2. Install dependencies:
go mod tidy
3. Run the application:
go run main.go

