package models

type Status struct {
	User   int32 `json:"user"`
	Forum  int32 `json:"forum"`
	Thread int32 `json:"thread"`
	Post   int32 `json:"post"`
}

func StatusInit() Status {
	return Status{
		User:   0,
		Forum:  0,
		Thread: 0,
		Post:   0,
	}
}
