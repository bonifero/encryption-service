package ru.bonifero.crypto.dto;

import io.swagger.v3.oas.annotations.media.Schema;
import ru.bonifero.crypto.entity.CryptoOperation;
import ru.bonifero.crypto.entity.OperationType;

import java.time.Instant;

@Schema(description = "Запись журнала выполненной криптографической операции")
public record CryptoOperationResponse(

        @Schema(description = "Идентификатор записи")
        Long id,

        @Schema(description = "Тип операции")
        OperationType operationType,

        @Schema(description = "Входные данные операции в формате JSON")
        String inputData,

        @Schema(description = "Результат операции в формате JSON")
        String outputData,

        @Schema(description = "Время выполнения операции")
        Instant createdAt
) {
    public static CryptoOperationResponse from(CryptoOperation entity) {
        return new CryptoOperationResponse(entity.getId(), entity.getOperationType(),
                entity.getInputData(), entity.getOutputData(), entity.getCreatedAt());
    }
}
