# Diffie Hellman Key Exchange
Go program that implements the Diffie-Hellman key exchange protocol at a 1024-bit key strength, this consists of parameter generation, share exchange, and key derivation. Normally the D-H scheme is an interactive protocol between two parties. However, for this implementation files are used as the communication channel.

## DH-Alice 1

The first models Alice’s initial message to Bob, and outputs a secret key to be stored for later. 

The program generates both public and private keys for Alice, this includes generation of prime number `p` using other prime numbers for its eulers totient `\phi(p)`. The program can be executed as shown below:

    dh-alice1 <filename for message to Bob> <filename to store secret key>

Outputs decimal-formatted public key ( p, g, ga ) for Bob and writes the secret key (p, g, a) to a second file.

## DH-Bob

The second models Bob’s receipt of the message from Alice and outputs a response message back to Alice.

The program generates the public and private keys fro Bob and uses the private key and Alice's public key to generate the shared key. The program can bbe executed as follow:


    dh-bob <filename of message from Alice> <filename of message back to Alice>

Reads in Alice’s message, outputs ( p, g, gb ) to Alice, prints the shared secret gab.

## DH-Alice 2

The final program models Alice’s receipt of Bob’s response. For grading purposes you will hand in the following programs.

The program retrieves the previously generates secret key and Bob's pubblic key to generate the shared key gab. The program can bbe executed as follow:

    dh-alice2 <filename of message from Bob> <filename to read secret key>

Reads in Bob’s message and Alice’s stored secret, prints the shared secret gab.

Hence both the parties were able to derive the shared secret key gab.