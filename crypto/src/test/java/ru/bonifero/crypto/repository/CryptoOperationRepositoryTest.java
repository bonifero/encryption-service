package ru.bonifero.crypto.repository;

import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.orm.jpa.DataJpaTest;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Sort;
import ru.bonifero.crypto.entity.CryptoOperation;
import ru.bonifero.crypto.entity.OperationType;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

@DataJpaTest
class CryptoOperationRepositoryTest {

    @Autowired
    private CryptoOperationRepository repository;

    @Test
    void save_thenFindById_returnsPersistedOperation() {
        CryptoOperation saved = repository.save(new CryptoOperation(OperationType.HASH, "\"hello\"", "{\"hashHex\":\"abc\"}"));

        CryptoOperation found = repository.findById(saved.getId()).orElseThrow();

        assertEquals(OperationType.HASH, found.getOperationType());
        assertEquals("\"hello\"", found.getInputData());
        assertEquals("{\"hashHex\":\"abc\"}", found.getOutputData());
    }

    @Test
    void findByOperationType_filtersAndPaginatesCorrectly() {
        repository.save(new CryptoOperation(OperationType.HASH, "in1", "out1"));
        repository.save(new CryptoOperation(OperationType.ENCRYPT, "in2", "out2"));
        repository.save(new CryptoOperation(OperationType.HASH, "in3", "out3"));

        Page<CryptoOperation> hashOnly = repository.findByOperationType(
                OperationType.HASH, PageRequest.of(0, 10, Sort.by(Sort.Direction.DESC, "createdAt")));

        assertEquals(2, hashOnly.getTotalElements());
        assertTrue(hashOnly.getContent().stream().allMatch(op -> op.getOperationType() == OperationType.HASH));
    }

    @Test
    void findAll_respectsPageSize() {
        for (int i = 0; i < 5; i++) {
            repository.save(new CryptoOperation(OperationType.HASH, "in" + i, "out" + i));
        }

        Page<CryptoOperation> firstPage = repository.findAll(PageRequest.of(0, 2));

        assertEquals(2, firstPage.getContent().size());
        assertEquals(5, firstPage.getTotalElements());
        assertEquals(3, firstPage.getTotalPages());
    }
}
