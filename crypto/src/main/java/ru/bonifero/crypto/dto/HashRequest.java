package ru.bonifero.crypto.dto;

import io.swagger.v3.oas.annotations.media.Schema;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import ru.bonifero.crypto.crypto.HashAlgorithm;

@Schema(description = "Запрос на расчет хеша сообщения")
public record HashRequest(

        @NotBlank(message = "Сообщение не должно быть пустым")
        @Schema(description = "Исходное сообщение в открытом виде", example = "Hello, world!")
        String message,

        @NotNull(message = "Алгоритм хеширования обязателен")
        @Schema(description = "Алгоритм расчета хеша", example = "SHA_256")
        HashAlgorithm algorithm
) {
}
