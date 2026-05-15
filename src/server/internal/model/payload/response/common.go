package response

type Success struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Errors  any    `json:"error,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type Paginate struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Meta    Meta   `json:"meta"`
}

type Meta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func NewSuccess(code int, message string, data any) Success {
	return Success{Code: code, Message: message, Data: data}
}

func NewError(code int, message string, errs any) Error {
	return Error{Code: code, Message: message, Errors: errs}
}

func NewPaginate(code int, message string, data any, page, pageSize int, total int64) Paginate {
	totalPages := 0
	if pageSize > 0 {
		totalPages = int((total + int64(pageSize) - 1) / int64(pageSize))
	}
	return Paginate{
		Code:    code,
		Message: message,
		Data:    data,
		Meta: Meta{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}
}
