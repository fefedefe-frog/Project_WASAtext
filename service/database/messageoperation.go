package database

func (db *appdbimpl) InsertMessage(message Message, chatId int) error {
	//TODO implement me
	panic("implement me")
}

func (db *appdbimpl) RemoveMessage(msgId int, chatId int) error {
	//TODO implement me
	panic("implement me")
}

func (db *appdbimpl) ForwardMessage(msgId int, chatIdToForwatd int) error {
	//TODO implement me
	panic("implement me")
}

func (db *appdbimpl) GetMessageComments(msgId int) ([]Comment, error) {
	//TODO implement me
	panic("implement me")
}

func (db *appdbimpl) CommentMessage(msgId int, comment Comment) error {
	//TODO implement me
	panic("implement me")
}

func (db *appdbimpl) UncommentMessage(msgId int, commentId int) error {
	//TODO implement me
	panic("implement me")
}