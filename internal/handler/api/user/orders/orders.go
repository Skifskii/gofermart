package orders

import (
	"errors"
	"gophermart/internal/middleware/authmw"
	"gophermart/internal/service/loyalty"
	"io"
	"net/http"
)

type OrderUploader interface {
	UploadOrderNum(userLogin, orderNum string) error
}

func New(oa OrderUploader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Получаем логин пользователя
		userLogin, ok := r.Context().Value(authmw.UserLoginKey).(string)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Не удалось прочитать тело запроса", http.StatusInternalServerError)
			return
		}

		orderNum := string(body)
		r.Body.Close()

		// Загружаем номер заказа в сервис
		if err := oa.UploadOrderNum(userLogin, orderNum); err != nil {
			if errors.Is(err, loyalty.ErrUserUploadedThisOrder) {
				// 200 - номер заказа уже был загружен этим пользователем
				w.WriteHeader(http.StatusOK)
				return
			}
			if errors.Is(err, loyalty.ErrWrongOrderNumFormat) {
				// 422 — неверный формат номера заказа
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			if errors.Is(err, loyalty.ErrAnotherUserUploadedThisOrder) {
				// 409 - номер заказа уже был загружен другим пользователем
				w.WriteHeader(http.StatusConflict)
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// 202 - новый номер заказа принят в обработку
		w.WriteHeader(http.StatusAccepted)
	}
}
