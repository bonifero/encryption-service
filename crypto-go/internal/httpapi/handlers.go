package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"crypto-go/internal/cryptoutil"
	"crypto-go/internal/model"
	"crypto-go/internal/repository"
)

type Handler struct {
	crypto *cryptoutil.Service
	repo   *repository.OperationRepository
}

func NewHandler(crypto *cryptoutil.Service, repo *repository.OperationRepository) *Handler {
	return &Handler{crypto: crypto, repo: repo}
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		Status:    status,
		Error:     message,
	})
}

func (h *Handler) Encrypt(w http.ResponseWriter, r *http.Request) {
	var req MessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Некорректное тело запроса")
		return
	}
	if strings.TrimSpace(req.Message) == "" {
		writeError(w, http.StatusBadRequest, "message: Сообщение не должно быть пустым")
		return
	}

	result, err := h.crypto.Encrypt(req.Message)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	response := EncryptResponse{EncryptedKey: result.EncryptedKey, IV: result.IV, CipherText: result.CipherText}
	if err := h.persist(r, model.OperationEncrypt, req, response); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, response)
}

func (h *Handler) Decrypt(w http.ResponseWriter, r *http.Request) {
	var req DecryptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Некорректное тело запроса")
		return
	}
	if strings.TrimSpace(req.EncryptedKey) == "" {
		writeError(w, http.StatusBadRequest, "encryptedKey: encryptedKey не должен быть пустым")
		return
	}
	if strings.TrimSpace(req.IV) == "" {
		writeError(w, http.StatusBadRequest, "iv: iv не должен быть пустым")
		return
	}
	if strings.TrimSpace(req.CipherText) == "" {
		writeError(w, http.StatusBadRequest, "cipherText: cipherText не должен быть пустым")
		return
	}

	message, err := h.crypto.Decrypt(cryptoutil.HybridEncryptionResult{
		EncryptedKey: req.EncryptedKey, IV: req.IV, CipherText: req.CipherText,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	response := DecryptResponse{Message: message}
	if err := h.persist(r, model.OperationDecrypt, req, response); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, response)
}

func (h *Handler) Hash(w http.ResponseWriter, r *http.Request) {
	var req HashRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Некорректное тело запроса")
		return
	}
	if strings.TrimSpace(req.Message) == "" {
		writeError(w, http.StatusBadRequest, "message: Сообщение не должно быть пустым")
		return
	}

	algorithm := cryptoutil.HashAlgorithm(req.Algorithm)
	if algorithm != cryptoutil.AlgorithmSHA256 && algorithm != cryptoutil.AlgorithmPBKDF2 {
		writeError(w, http.StatusBadRequest, "algorithm: Алгоритм хеширования обязателен (SHA_256 или PBKDF2)")
		return
	}

	result, err := h.crypto.Hash(req.Message, algorithm)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	response := HashResponse{
		Algorithm: result.Algorithm, HashHex: result.HashHex,
		SaltHex: result.SaltHex, Iterations: result.Iterations,
	}
	if err := h.persist(r, model.OperationHash, req, response); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, response)
}

func (h *Handler) Logs(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	page := 0
	if v := query.Get("page"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			page = parsed
		}
	}
	size := 20
	if v := query.Get("size"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			size = parsed
		}
	}

	var typeFilter *model.OperationType
	if v := query.Get("type"); v != "" {
		t := model.OperationType(v)
		if !t.Valid() {
			writeError(w, http.StatusBadRequest, "type: неизвестный тип операции "+v)
			return
		}
		typeFilter = &t
	}

	operations, total, err := h.repo.FindPage(r.Context(), typeFilter, page, size)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	content := make([]CryptoOperationResponse, 0, len(operations))
	for _, op := range operations {
		content = append(content, toResponse(op))
	}

	totalPages := int64(0)
	if size > 0 {
		totalPages = (total + int64(size) - 1) / int64(size)
	}

	writeJSON(w, http.StatusOK, PageResponse{
		Content:          content,
		TotalElements:    total,
		TotalPages:       totalPages,
		Size:             size,
		Number:           page,
		NumberOfElements: len(content),
		First:            page == 0,
		Last:             int64(page+1) >= totalPages,
		Empty:            len(content) == 0,
	})
}

func (h *Handler) LogByID(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id должен быть числом")
		return
	}

	op, err := h.repo.FindByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusNotFound, "Запись с id="+idParam+" не найдена")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, toResponse(op))
}

func toResponse(op model.CryptoOperation) CryptoOperationResponse {
	return CryptoOperationResponse{
		ID:            op.ID,
		OperationType: string(op.OperationType),
		InputData:     op.InputData,
		OutputData:    op.OutputData,
		CreatedAt:     op.CreatedAt,
	}
}

func (h *Handler) persist(r *http.Request, opType model.OperationType, input, output any) error {
	inputJSON, err := json.Marshal(input)
	if err != nil {
		return err
	}
	outputJSON, err := json.Marshal(output)
	if err != nil {
		return err
	}
	_, err = h.repo.Save(r.Context(), model.CryptoOperation{
		OperationType: opType,
		InputData:     string(inputJSON),
		OutputData:    string(outputJSON),
		CreatedAt:     time.Now().UTC(),
	})
	return err
}
