package ru.bonifero.crypto.dto;

import io.swagger.v3.oas.annotations.media.Schema;
import jakarta.validation.constraints.NotBlank;

@Schema(description = "Запрос с сообщением для шифрования или расчета хеша")
public record MessageRequest(

        @NotBlank(message = "Сообщение не должно быть пустым")
        @Schema(description = "Исходное сообщение в открытом виде", example = "Hello, world!")
        String message
) {
}
