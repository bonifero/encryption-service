package ru.bonifero.crypto.entity;

import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.EnumType;
import jakarta.persistence.Enumerated;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Id;
import jakarta.persistence.Table;

import java.time.Instant;

@Entity
@Table(name = "crypto_operation")
public class CryptoOperation {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Enumerated(EnumType.STRING)
    @Column(nullable = false)
    private OperationType operationType;

    @Column(nullable = false, columnDefinition = "TEXT")
    private String inputData;

    @Column(nullable = false, columnDefinition = "TEXT")
    private String outputData;

    @Column(nullable = false)
    private Instant createdAt;

    protected CryptoOperation() {
    }

    public CryptoOperation(OperationType operationType, String inputData, String outputData) {
        this.operationType = operationType;
        this.inputData = inputData;
        this.outputData = outputData;
        this.createdAt = Instant.now();
    }

    public Long getId() {
        return id;
    }

    public OperationType getOperationType() {
        return operationType;
    }

    public String getInputData() {
        return inputData;
    }

    public String getOutputData() {
        return outputData;
    }

    public Instant getCreatedAt() {
        return createdAt;
    }
}
