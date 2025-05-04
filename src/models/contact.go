package models

type Contact struct {
	Id    int
	Name  string
	Email string
}

type Contacts = []Contact

type Data struct {
	Contacts Contacts
}

func NewContact(id int, name string, email string) Contact {
	return Contact{
		Id:    id,
		Name:  name,
		Email: email,
	}
}

func NewData() Data {
	return Data{
		Contacts: []Contact{
			NewContact(1, "John Doe", "jd@myemail.com"),
			NewContact(2, "Jane Doe", "jd2@myemail.com"),
		}}
}
