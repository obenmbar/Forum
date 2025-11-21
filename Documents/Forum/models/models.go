package forumino

type ErrorPage struct {
	Messege string
	Code int
}
type MY_Post struct {
	Title string
	Contenue string
	Category []string
}
type Category struct {
 Categoryid map[string]int
}
var CategorID = &Category{
	Categoryid: make(map[string]int),
}
type PageData struct {
	ErrorMessege error
	Post MY_Post
	CSRFToken string
}
type ReactionData struct {
	ContentType string
	ContentId int
	ContentReaction string
}