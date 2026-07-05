package ru.bonifero.crypto.dto;

import io.swagger.v3.oas.annotations.media.Schema;
import jakarta.validation.constraints.NotBlank;

@Schema(description = "Данные, полученные ранее из /api/crypto/encrypt, для дешифрования")
public record DecryptRequest(

        @NotBlank(message = "encryptedKey не должен быть пустым")
        @Schema(description = "AES-ключ, зашифрованный открытым RSA-ключом (Base64)")
        String encryptedKey,

        @NotBlank(message = "iv не должен быть пустым")
        @Schema(description = "Вектор инициализации AES/GCM (Base64)")
        String iv,

        @NotBlank(message = "cipherText не должен быть пустым")
        @Schema(description = "Зашифрованное сообщение (Base64)")
        String cipherText
) {
}
