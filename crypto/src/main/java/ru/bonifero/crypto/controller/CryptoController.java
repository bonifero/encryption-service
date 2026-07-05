package ru.bonifero.crypto.controller;

import com.fasterxml.jackson.databind.ObjectMapper;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.data.domain.Sort;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.server.ResponseStatusException;
import ru.bonifero.crypto.crypto.CryptoService;
import ru.bonifero.crypto.crypto.HashResult;
import ru.bonifero.crypto.crypto.HybridEncryptionResult;
import ru.bonifero.crypto.dto.CryptoOperationResponse;
import ru.bonifero.crypto.dto.DecryptRequest;
import ru.bonifero.crypto.dto.DecryptResponse;
import ru.bonifero.crypto.dto.EncryptResponse;
import ru.bonifero.crypto.dto.HashRequest;
import ru.bonifero.crypto.dto.HashResponse;
import ru.bonifero.crypto.dto.MessageRequest;
import ru.bonifero.crypto.entity.CryptoOperation;
import ru.bonifero.crypto.entity.OperationType;
import ru.bonifero.crypto.exception.CryptoServiceException;
import ru.bonifero.crypto.repository.CryptoOperationRepository;

@RestController
@RequestMapping("/api/crypto")
@Tag(name = "Crypto", description = "Шифрование, дешифрование и расчет хеша сообщений")
public class CryptoController {

    private final CryptoService cryptoService;
    private final CryptoOperationRepository repository;
    private final ObjectMapper objectMapper;

    public CryptoController(CryptoService cryptoService, CryptoOperationRepository repository,
                             ObjectMapper objectMapper) {
        this.cryptoService = cryptoService;
        this.repository = repository;
        this.objectMapper = objectMapper;
    }

    @Operation(summary = "Зашифровать сообщение",
            description = "Гибридное шифрование AES-256/GCM + RSA-OAEP. Возвращает зашифрованный AES-ключ, IV и шифротекст.")
    @PostMapping("/encrypt")
    public EncryptResponse encrypt(@Valid @RequestBody MessageRequest request) {
        HybridEncryptionResult result = cryptoService.encrypt(request.message());
        EncryptResponse response = new EncryptResponse(result.encryptedKey(), result.iv(), result.cipherText());
        persist(OperationType.ENCRYPT, request.message(), response);
        return response;
    }

    @Operation(summary = "Расшифровать сообщение",
            description = "Дешифрует сообщение, зашифрованное методом /encrypt приватным RSA-ключом сервиса.")
    @PostMapping("/decrypt")
    public DecryptResponse decrypt(@Valid @RequestBody DecryptRequest request) {
        HybridEncryptionResult encrypted = new HybridEncryptionResult(
                request.encryptedKey(), request.iv(), request.cipherText());
        String message = cryptoService.decrypt(encrypted);
        DecryptResponse response = new DecryptResponse(message);
        persist(OperationType.DECRYPT, request, response);
        return response;
    }

    @Operation(summary = "Рассчитать хеш сообщения",
            description = "Поддерживаемые алгоритмы: SHA_256, PBKDF2 (со случайной солью).")
    @PostMapping("/hash")
    public HashResponse hash(@Valid @RequestBody HashRequest request) {
        HashResult result = cryptoService.hash(request.message(), request.algorithm());
        HashResponse response = new HashResponse(result.algorithm(), result.hashHex(),
                result.saltHex(), result.iterations());
        persist(OperationType.HASH, request.message(), response);
        return response;
    }

    @Operation(summary = "Получить журнал выполненных криптографических операций",
            description = "Постраничный список записей с входными данными и результатом каждой операции. "
                    + "Параметр type опционален и фильтрует по конкретному методу (ENCRYPT, DECRYPT, HASH).")
    @GetMapping("/logs")
    public Page<CryptoOperationResponse> logs(
            @RequestParam(required = false) OperationType type,
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "20") int size) {
        Pageable pageable = PageRequest.of(page, size, Sort.by(Sort.Direction.DESC, "createdAt"));
        Page<CryptoOperation> result = (type != null)
                ? repository.findByOperationType(type, pageable)
                : repository.findAll(pageable);
        return result.map(CryptoOperationResponse::from);
    }

    @Operation(summary = "Получить одну запись журнала операций по идентификатору")
    @GetMapping("/logs/{id}")
    public CryptoOperationResponse logById(@PathVariable Long id) {
        return repository.findById(id)
                .map(CryptoOperationResponse::from)
                .orElseThrow(() -> new ResponseStatusException(HttpStatus.NOT_FOUND,
                        "Запись с id=" + id + " не найдена"));
    }

    private void persist(OperationType type, Object input, Object output) {
        try {
            String inputJson = objectMapper.writeValueAsString(input);
            String outputJson = objectMapper.writeValueAsString(output);
            repository.save(new CryptoOperation(type, inputJson, outputJson));
        } catch (Exception e) {
            throw new CryptoServiceException("Ошибка при сохранении результата операции в БД", e);
        }
    }
}
