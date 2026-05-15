package request

// UpdateProfile is a partial update — only fields explicitly supplied are
// applied. Pointer types let us distinguish "field omitted" (nil) from
// "field cleared" (empty string).
type UpdateProfile struct {
	Name        *string `json:"name,omitempty"         binding:"omitempty,personname"`
	Surname     *string `json:"surname,omitempty"      binding:"omitempty,personname"`
	PhoneNumber *string `json:"phone_number,omitempty" binding:"omitempty,phone"`
	Address     *string `json:"address,omitempty"      binding:"omitempty,max=500"`
}
