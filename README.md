This started out as trying to track down why I'm seeing `Closed explicitly` returned from database queries when under load.  I believe this minimal reproduction shows this error as well as either `connection reset` or `read failed`, maybe due to mongo starving.

# Warning

- Heavy cpu usage
- Results don't happen on every run

# Running

```
docker-compose up -d mongo

sleep 10

docker-compose run tests
```

# Questions
- How is the best way to handle a flaky connection with the database?
- Is the above correct in that `Closed explicitly` is coming from a flaky
  mongo connection?


# Other notes

In our actual codebase the only error I'm seeing is 

`Closed explicitly`

I believe this is the same error we're seeing in this repo, with the difference in output due to having the following.

```
func NewSession(url string) (m *mgo.Session, err error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	url = strings.TrimSuffix(url, "?ssl=true")
	dialInfo, err := mgo.ParseURL(url)
	if err != nil {
		return m, err
	}

	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		if err != nil {
			fmt.Println(err.Error())
		}
		return conn, err
	}

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		return m, err
	}

	session.SetMode(mgo.Monotonic, true)

	return session, nil
}
```
