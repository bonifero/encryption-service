package ru.bonifero.crypto.dto;

import io.swagger.v3.oas.annotations.media.Schema;

@Schema(description = "Результат гибридного шифрования сообщения (AES-256/GCM + RSA-OAEP)")
public record EncryptResponse(

        @Schema(description = "AES-ключ, зашифрованный открытым RSA-ключом (Base64)")
        String encryptedKey,

        @Schema(description = "Вектор инициализации AES/GCM (Base64)")
        String iv,

        @Schema(description = "Зашифрованное сообщение (Base64)")
        String cipherText
) {
}
