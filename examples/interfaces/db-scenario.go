package main

// Scenario: User Service with a Database (Tightly Coupled)
//
// You start simple with MySQL:
//

type MySQLDB struct{}

func (db MySQLDB) GetUser(id int) string {
	return "user from MySQL"
}

type UserService struct {
	db MySQLDB
}

func (s UserService) GetUserName(id int) string {
	return s.db.GetUser(id)
}

// New requirements come in:
// “We might switch to PostgreSQL later”
// “We need Redis caching”
// “We want fast unit tests (no real DB)”

// for switching DB now you rewrite service
// like below
//

type PostgresDB struct{}

type UserService struct {
	db PostgresDB // changed everywhere
}

// You must modify every place that used MySQLDB.
//
// Adding cache = messy logic

func (s UserService) GetUserName(id int) string {
	if cached := redis.Get(id); cached != "" {
		return cached
	}

	user := s.db.GetUser(id)
	redis.Set(id, user)
	return user
}

// Now your business logic is polluted with DB + cache logic.
//
// Testing becomes painful
//
// cant do

UserService{db: MySQLDB{}}

// it hits real db and slow / needs setup just for testing


// Same Problem Using Interface (Clean Design)
//
// Step 1: Define behavior

type UserRepository interface {
    GetUser(id int) string
}

// Step 2: Implement it

type MySQLDB struct{}

func (db MySQLDB) GetUser(id int) string {
    return "user from MySQL"
}

type PostgresDB struct{}

func (db PostgresDB) GetUser(id int) string {
    return "user from Postgres"
}

// Step 3: Use interface in service

type UserService struct {
    repo UserRepository
}

func (s UserService) GetUserName(id int) string {
    return s.repo.GetUser(id)
}

// now it is easier to scale and test
// Switch DB without touching service

service := UserService{repo: PostgresDB{}} // repo can be any new implementation like mongo or cassandra

// Add caching (cleanly)
// Wrap the repo:

type CachedRepo struct {
    next UserRepository
}

func (c CachedRepo) GetUser(id int) string {
    // pretend cache logic
    user := c.next.GetUser(id)
    return user
}

// Testing becomes trivial

type FakeRepo struct{}
func (f FakeRepo) GetUser(id int) string {
    return "fake user"
}

service := UserService{repo: FakeRepo{}}

// What Just Happened

// Before:

// UserService → MySQLDB (hard dependency)

// After:

// UserService → UserRepository (interface) → anything
