package ru.bonifero.crypto.dto;

import io.swagger.v3.oas.annotations.media.Schema;

@Schema(description = "Результат расчета хеша сообщения")
public record HashResponse(

        @Schema(description = "Использованный алгоритм", example = "SHA-256")
        String algorithm,

        @Schema(description = "Хеш в шестнадцатеричном представлении")
        String hashHex,

        @Schema(description = "Соль в шестнадцатеричном представлении (только для PBKDF2)", nullable = true)
        String saltHex,

        @Schema(description = "Количество итераций (только для PBKDF2)", nullable = true)
        Integer iterations
) {
}
