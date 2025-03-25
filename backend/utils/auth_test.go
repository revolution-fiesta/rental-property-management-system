package utils

import (
    "errors"
    "testing"
)

type testUsername struct {
    username string
    wants    error
}

func TestCheckUsername(t *testing.T) {
    tests := []testUsername{
        {username: "12345654", wants: errors.New("Usernames start with numbers")},
        {username: "123456!4", wants: errors.New("Usernames start with numbers")},
        {username: "!!!3!!!6!!4", wants: errors.New("Usernames contain special characters")},
        {username: "!!!", wants: errors.New("The username length is too short")},
        {username: "!!!!!!!kfhkdjfhkdsjf!!!", wants: errors.New("The username length is too long")},
        {username: "u_inthds3", wants: nil},
        {username: "", wants: errors.New("The username is empty")},
        {username: "12345", wants: errors.New("The username length is too short")},
        {username: "1234543894872648734", wants: errors.New("The username length is too long")},
        {username: "u12345", wants: nil},
    }
    for _, test := range tests {
        got := CheckUsername(test.username)
        if (got != nil || test.wants != nil) && got.Error() != test.wants.Error() {
            t.Errorf("CheckUsername(%s) = %v; want %v", test.username, got, test.wants)
        }
    }

}




