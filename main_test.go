package main

import (
	"fmt"
	"os"
	"sync"
	"testing"

	uuid "github.com/satori/go.uuid"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func Collection(s *mgo.Session) *mgo.Collection {
	return s.DB("testing").C("collection")
}

func TestHighVolumeSpec2(t *testing.T) {
	url := os.Getenv("TEST_MONGO_URL")

	// Get the session
	session, err := mgo.Dial(url)
	if err != nil {
		fmt.Println(err.Error())
		t.FailNow()
		return
	}

	session.SetMode(mgo.Monotonic, true)

	// Collection(session).RemoveAll(nil)
	fmt.Println("Starting reproduction")

	count := 30000
	wg := &sync.WaitGroup{}
	wg.Add(count)

	for i := 0; i < count; i++ {
		index := i
		go func() {
			err := UpdateIdentity(session)
			if err != nil {
				fmt.Println(err.Error())
				t.FailNow()
				panic(err)
			} else {
				if index%1000 == 0 {
					fmt.Println(index / 1000)
				}
			}
			wg.Done()
		}()
	}

	wg.Wait()
}

func UpdateIdentity(session *mgo.Session) (err error) {
	s := session.Clone()
	defer s.Close()
	c := Collection(s)

	query := &bson.M{
		"someVal":       uuid.NewV4().String(),
		"somerOtherVal": "test",
	}
	mod := &bson.D{}

	_, err = c.Upsert(query, mod)
	return err
}
