package ru.bonifero.crypto.repository;

import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import ru.bonifero.crypto.entity.CryptoOperation;
import ru.bonifero.crypto.entity.OperationType;

public interface CryptoOperationRepository extends JpaRepository<CryptoOperation, Long> {

    Page<CryptoOperation> findByOperationType(OperationType operationType, Pageable pageable);
}
