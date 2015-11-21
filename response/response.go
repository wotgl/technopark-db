package response

type Foo interface {
	Speak() string
}

type ArrayThreadsDetails struct {
	Date      string `json:"date"`
	Dislikes  int64  `json:"dislikes"`
	Forum     string `json:"forum"`
	Id        int64  `json:"id"`
	IsClosed  bool   `json:"isClosed"`
	IsDeleted bool   `json:"isDeleted"`
	Likes     int64  `json:"likes"`
	Message   string `json:"message"`
	Points    int    `json:"points"`
	Posts     int    `json:"posts"`
	Slug      string `json:"slug"`
	Title     string `json:"title"`
	User      string `json:"user"`
}

func (d ArrayThreadsDetails) Speak() string {
	return "Woof!"
}

type ErrorMessage struct {
	Msg string `json:"msg"`
}

func (d ErrorMessage) Speak() string { return "Woof!" }

type Open struct {
	Thread float64 `json:"thread"`
}

func (d Open) Speak() string { return "Woof!" }

type Remove struct {
	Thread float64 `json:"thread"`
}

func (d Remove) Speak() string { return "Woof!" }

type Restore struct {
	Thread float64 `json:"thread"`
}

func (d Restore) Speak() string { return "Woof!" }

type Subscribe struct {
	Thread float64 `json:"thread"`
	User   string  `json:"user"`
}

func (d Subscribe) Speak() string { return "Woof!" }

type Unsubscribe struct {
	Thread float64 `json:"thread"`
	User   string  `json:"user"`
}

func (d Unsubscribe) Speak() string { return "Woof!" }

type ThreadDetails struct {
	Date      string `json:"date"`
	Dislikes  int64  `json:"dislikes"`
	Forum     string `json:"forum"`
	Id        int64  `json:"id"`
	IsClosed  bool   `json:"isClosed"`
	IsDeleted bool   `json:"isDeleted"`
	Likes     int64  `json:"likes"`
	Message   string `json:"message"`
	Points    int64  `json:"points"`
	Posts     int64  `json:"posts"`
	Slug      string `json:"slug"`
	Title     string `json:"title"`
	User      string `json:"user"`
}

func (d ThreadDetails) Speak() string { return "Woof!" }
