package ru.bonifero.crypto.crypto;

public record HashResult(String algorithm, String hashHex, String saltHex, Integer iterations) {
}
