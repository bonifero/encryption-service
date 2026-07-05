package ru.bonifero.crypto.dto;

import io.swagger.v3.oas.annotations.media.Schema;

@Schema(description = "Результат дешифрования сообщения")
public record DecryptResponse(

        @Schema(description = "Исходное сообщение в открытом виде")
        String message
) {
}
