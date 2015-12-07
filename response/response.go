package response

type RespStruct interface {
	Foo() bool
}
type UserUnfollow struct {
	Username      string      `json:"username"`
	About         string      `json:"about"`
	Name          string      `json:"name"`
	Subscriptions interface{} `json:"subscriptions"`
	Id            int64       `json:"id"`
	Followers     interface{} `json:"followers"`
	Following     interface{} `json:"following"`
	Isanonymous   bool        `json:"isAnonymous"`
	Email         string      `json:"email"`
}

func (instance *UserUnfollow) Foo() bool { return true }

type UserFollow struct {
	Username      string      `json:"username"`
	About         string      `json:"about"`
	Name          string      `json:"name"`
	Subscriptions interface{} `json:"subscriptions"`
	Id            int64       `json:"id"`
	Followers     interface{} `json:"followers"`
	Following     interface{} `json:"following"`
	Isanonymous   bool        `json:"isAnonymous"`
	Email         string      `json:"email"`
}

func (instance *UserFollow) Foo() bool { return true }

type UserDetails struct {
	Username      string      `json:"username"`
	About         string      `json:"about"`
	Name          string      `json:"name"`
	Subscriptions interface{} `json:"subscriptions"`
	Id            int64       `json:"id"`
	Followers     interface{} `json:"followers"`
	Following     interface{} `json:"following"`
	Isanonymous   bool        `json:"isAnonymous"`
	Email         string      `json:"email"`
}

func (instance *UserDetails) Foo() bool { return true }

type UserListFollowing struct {
	Users []UserDetails
}

func (instance *UserListFollowing) Foo() bool { return true }

type UserListPosts struct {
	Users []UserDetails
}

func (instance *UserListPosts) Foo() bool { return true }

type UserCreate struct {
	Username    string `json:"username"`
	About       string `json:"about"`
	Name        string `json:"name"`
	Id          int64  `json:"id"`
	IsAnonymous bool   `json:"isAnonymous"`
	Email       string `json:"email"`
}

func (instance *UserCreate) Foo() bool { return true }

type UserListFollowers struct {
	Users []UserDetails
}

func (instance *UserListFollowers) Foo() bool { return true }

type UserUpdateProfile struct {
	Username      string      `json:"username"`
	About         string      `json:"about"`
	Name          string      `json:"name"`
	Subscriptions interface{} `json:"subscriptions"`
	Id            int64       `json:"id"`
	Followers     interface{} `json:"followers"`
	Following     interface{} `json:"following"`
	Isanonymous   bool        `json:"isAnonymous"`
	Email         string      `json:"email"`
}

func (instance *UserUpdateProfile) Foo() bool { return true }

type ForumListUsers struct {
	Users []UserDetails
}

func (instance *ForumListUsers) Foo() bool { return true }

type ForumListThreads struct {
	Users []ThreadList
}

func (instance *ForumListThreads) Foo() bool { return true }

type ForumDetails struct {
	User       UserDetails `json:"user"`
	Id         int64       `json:"id"`
	Short_Name string      `json:"short_name"`
	Name       string      `json:"name"`
}

func (instance *ForumDetails) Foo() bool { return true }

type ForumListPosts struct {
	Posts []PostDetails
}

func (instance *ForumListPosts) Foo() bool { return true }

type ForumCreate struct {
	User       string `json:"user"`
	Id         int64  `json:"id"`
	Short_Name string `json:"short_name"`
	Name       string `json:"name"`
}

func (instance *ForumCreate) Foo() bool { return true }

type ThreadBoolBasic struct {
	Thread float64 `json:"thread"`
}

func (instance *ThreadBoolBasic) Foo() bool { return true }

type ThreadOpen struct {
	Thread int64 `json:"thread"`
}

func (instance *ThreadOpen) Foo() bool { return true }

type ThreadUpdate struct {
	Forum     string `json:"forum"`
	Title     string `json:"title"`
	Posts     int64  `json:"posts"`
	Dislikes  int64  `json:"dislikes"`
	Slug      string `json:"slug"`
	Isclosed  bool   `json:"isClosed"`
	Points    int64  `json:"points"`
	Likes     int64  `json:"likes"`
	Date      string `json:"date"`
	Message   string `json:"message"`
	Id        int64  `json:"id"`
	Isdeleted bool   `json:"isDeleted"`
	User      string `json:"user"`
}

func (instance *ThreadUpdate) Foo() bool { return true }

type ThreadUnsubscribe struct {
	User   string `json:"user"`
	Thread int64  `json:"thread"`
}

func (instance *ThreadUnsubscribe) Foo() bool { return true }

type ThreadRestore struct {
	Thread int64 `json:"thread"`
}

func (instance *ThreadRestore) Foo() bool { return true }

type ThreadDetails struct {
	Forum     string `json:"forum"`
	Title     string `json:"title"`
	Posts     int64  `json:"posts"`
	Dislikes  int64  `json:"dislikes"`
	Slug      string `json:"slug"`
	Isclosed  bool   `json:"isClosed"`
	Points    int64  `json:"points"`
	Likes     int64  `json:"likes"`
	Date      string `json:"date"`
	Message   string `json:"message"`
	Id        int64  `json:"id"`
	Isdeleted bool   `json:"isDeleted"`
	User      string `json:"user"`
}

func (instance *ThreadDetails) Foo() bool { return true }

type ThreadRemove struct {
	Thread int64 `json:"thread"`
}

func (instance *ThreadRemove) Foo() bool { return true }

type ThreadVote struct {
	Forum     string `json:"forum"`
	Title     string `json:"title"`
	Posts     int64  `json:"posts"`
	Dislikes  int64  `json:"dislikes"`
	Slug      string `json:"slug"`
	Isclosed  bool   `json:"isClosed"`
	Points    int64  `json:"points"`
	Likes     int64  `json:"likes"`
	Date      string `json:"date"`
	Message   string `json:"message"`
	Id        int64  `json:"id"`
	Isdeleted bool   `json:"isDeleted"`
	User      string `json:"user"`
}

func (instance *ThreadVote) Foo() bool { return true }

type ThreadClose struct {
	Thread int64 `json:"thread"`
}

func (instance *ThreadClose) Foo() bool { return true }

type ThreadList struct {
	Threads []ThreadDetails
}

func (instance *ThreadList) Foo() bool { return true }

type ThreadSubscribe struct {
	User   string `json:"user"`
	Thread int64  `json:"thread"`
}

func (instance *ThreadSubscribe) Foo() bool { return true }

type ThreadListPosts struct {
	Posts []PostDetails
}

func (instance *ThreadListPosts) Foo() bool { return true }

type ThreadCreate struct {
	Forum     string `json:"forum"`
	Title     string `json:"title"`
	Slug      string `json:"slug"`
	IsClosed  bool   `json:"isClosed"`
	User      string `json:"user"`
	Date      string `json:"date"`
	Message   string `json:"message"`
	Id        int64  `json:"id"`
	IsDeleted bool   `json:"isDeleted"`
}

func (instance *ThreadCreate) Foo() bool { return true }

type PostUpdate struct {
	Parent        int64  `json:"parent"`
	Forum         string `json:"forum"`
	Isapproved    bool   `json:"isApproved"`
	User          string `json:"user"`
	Dislikes      int64  `json:"dislikes"`
	Isspam        bool   `json:"isSpam"`
	Thread        int64  `json:"thread"`
	Points        int64  `json:"points"`
	Ishighlighted bool   `json:"isHighlighted"`
	Isedited      bool   `json:"isEdited"`
	Date          string `json:"date"`
	Message       string `json:"message"`
	Id            int64  `json:"id"`
	Isdeleted     bool   `json:"isDeleted"`
	Likes         int64  `json:"likes"`
}

func (instance *PostUpdate) Foo() bool { return true }

type PostRestore struct {
	Post int64 `json:"post"`
}

func (instance *PostRestore) Foo() bool { return true }

type PostDetails struct {
	Parent        int64  `json:"parent"`
	Forum         string `json:"forum"`
	Isapproved    bool   `json:"isApproved"`
	User          string `json:"user"`
	Dislikes      int64  `json:"dislikes"`
	Isspam        bool   `json:"isSpam"`
	Thread        int64  `json:"thread"`
	Points        int64  `json:"points"`
	Ishighlighted bool   `json:"isHighlighted"`
	Isedited      bool   `json:"isEdited"`
	Date          string `json:"date"`
	Message       string `json:"message"`
	Id            int64  `json:"id"`
	Isdeleted     bool   `json:"isDeleted"`
	Likes         int64  `json:"likes"`
}

func (instance *PostDetails) Foo() bool { return true }

type PostRemove struct {
	Post int64 `json:"post"`
}

func (instance *PostRemove) Foo() bool { return true }

type PostVote struct {
	Parent        int64  `json:"parent"`
	Forum         string `json:"forum"`
	Isapproved    bool   `json:"isApproved"`
	User          string `json:"user"`
	Dislikes      int64  `json:"dislikes"`
	Isspam        bool   `json:"isSpam"`
	Thread        int64  `json:"thread"`
	Points        int64  `json:"points"`
	Ishighlighted bool   `json:"isHighlighted"`
	Isedited      bool   `json:"isEdited"`
	Date          string `json:"date"`
	Message       string `json:"message"`
	Id            int64  `json:"id"`
	Isdeleted     bool   `json:"isDeleted"`
	Likes         int64  `json:"likes"`
}

func (instance *PostVote) Foo() bool { return true }

type PostList struct {
	Posts []PostDetails
}

func (instance *PostList) Foo() bool { return true }

type PostCreate struct {
	Parent        int64  `json:"parent"`
	Forum         string `json:"forum"`
	Isapproved    bool   `json:"isApproved"`
	Isspam        bool   `json:"isSpam"`
	Thread        int64  `json:"thread"`
	Ishighlighted bool   `json:"isHighlighted"`
	Isedited      bool   `json:"isEdited"`
	Date          string `json:"date"`
	Message       string `json:"message"`
	Id            int64  `json:"id"`
	Isdeleted     bool   `json:"isDeleted"`
	User          string `json:"user"`
}

func (instance *PostCreate) Foo() bool { return true }
