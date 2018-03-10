package frameworks

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nightlyone/lockfile"
	"github.com/riomhaire/lightauthuserapi/entities"
	"github.com/riomhaire/lightauthuserapi/usecases"
)

const (
	usernameField = 0
	passwordField = 1
	enabledField  = 2
	rolesField    = 3
	claim1Field   = 4
	claim2Field   = 5
	roleNameField = 0
)

// This is a test implementation for test purposes
type CSVReaderDatabaseInteractor struct {
	registry    *usecases.Registry
	userdb      map[string]entities.User
	roledb      []entities.Role
	names       []string
	initialized bool
	mux         sync.Mutex
}

func NewCSVReaderDatabaseInteractor(registry *usecases.Registry) *CSVReaderDatabaseInteractor {
	d := CSVReaderDatabaseInteractor{}
	d.userdb = make(map[string]entities.User)
	d.roledb = make([]entities.Role, 0)
	d.registry = registry

	return &d
}

func (db *CSVReaderDatabaseInteractor) LookupUserByName(username string) (entities.User, error) {
	db.lazyLoad()
	if val, ok := db.userdb[username]; ok {
		return val, nil
	} else {
		return entities.User{}, errors.New("Unknown user")
	}
}

func (db *CSVReaderDatabaseInteractor) CreateUser(user entities.User) error {
	db.lazyLoad()
	if _, ok := db.userdb[user.Username]; ok {
		return errors.New("User exists")
	}
	db.userdb[user.Username] = user
	db.writeUsers()
	db.rebuildNameIndex()
	return nil
}

func (db *CSVReaderDatabaseInteractor) LookupUserNames(search string, page int, pageSize int) ([]string, error) {
	db.lazyLoad()

	if len(search) > 0 {
		// How many results
		minReq := min(pageSize, len(db.names))
		if minReq == -1 {
			minReq = len(db.names)
		}
		i := 0
		matchNames := make([]string, 0)
		before := time.Now()
		for len(matchNames) < minReq && i < len(db.names) {
			potentialMatch := db.names[i]
			if strings.Contains(potentialMatch, search) {
				matchNames = append(matchNames, potentialMatch)
			}
			i++
		}
		now := time.Now()
		diff := now.Sub(before)
		db.registry.Logger.Log("DEBUG", fmt.Sprintf("Search for '%v' and %v hits took %v", search, len(matchNames), diff))

		return matchNames, nil
	}

	return db.names, nil
}

func (db *CSVReaderDatabaseInteractor) UpdateUser(user entities.User) error {
	db.lazyLoad()
	if _, ok := db.userdb[user.Username]; ok {
		db.userdb[user.Username] = user
	} else {
		return errors.New("User Does Not Exists")
	}
	// Flush to file
	db.writeUsers()

	return nil
}

func (db *CSVReaderDatabaseInteractor) DeleteUser(user string) error {
	db.lazyLoad()
	if _, ok := db.userdb[user]; ok {
		delete(db.userdb, user)
	} else {
		return errors.New("User Does Not Exists")
	}
	// Flush to file
	db.writeUsers()
	db.rebuildNameIndex()
	return nil
}

func (db *CSVReaderDatabaseInteractor) LookupRoleNames() ([]string, error) {
	db.lazyLoad()
	var roles []string
	for _, r := range db.roledb {
		roles = append(roles, r.Name)
	}
	return roles, nil
}

// Initiaizes data structues - IE Read user DB
func (db *CSVReaderDatabaseInteractor) loadUsers() (map[string]entities.User, error) {
	filename := db.registry.Configuration.UserStore
	users := make(map[string]entities.User)

	db.registry.Logger.Log("INFO", fmt.Sprintf("Reading User Database %s", filename))
	// If filename is none - dont load (test usage)
	if strings.Compare("NONE", strings.ToUpper(filename)) == 0 {
		return users, nil
	}

	csvfile, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
		return users, err
	}
	defer csvfile.Close()
	r := csv.NewReader(csvfile)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
		return users, err
	}
	// Create user map
	for index, row := range records {
		if index > 0 && len(row) > 0 {
			user := entities.User{}
			user.Username = row[usernameField]
			user.Password = row[passwordField]

			v, _ := strconv.ParseBool(row[enabledField])
			user.Enabled = v
			roles := strings.Split(row[rolesField], ":")
			user.Roles = roles
			user.Claim1 = row[claim1Field]
			user.Claim2 = row[claim2Field]

			// Add
			users[user.Username] = user
		}
	}
	db.userdb = users
	db.rebuildNameIndex()
	db.registry.Logger.Log("INFO", fmt.Sprintf("#Number of users = %v", len(users)))
	return users, nil
}

func (db *CSVReaderDatabaseInteractor) rebuildNameIndex() {
	names := make([]string, 0)

	// recreate search index
	for k := range db.userdb {
		names = append(names, k)
	}
	// Sort names
	sort.Strings(names)
	db.names = names
}

func (db *CSVReaderDatabaseInteractor) writeUsers() error {
	db.registry.Logger.Log("INFO", fmt.Sprintf("Writing User Database %s", db.registry.Configuration.UserStore))
	// If filename is none - dont load (test usage)
	if strings.Compare("NONE", strings.ToUpper(db.registry.Configuration.UserStore)) == 0 {
		return nil
	}
	// Setup file lock
	lock, err := lockfile.New(filepath.Join(os.TempDir(), "userstore.lock"))
	if err != nil {
		fmt.Printf("Cannot init lock. reason: %v", err)
		return nil // May loose - need to retry
	}
	err = lock.TryLock()

	// Error handling is essential, as we only try to get the lock.
	if err != nil {
		fmt.Printf("Cannot lock %q, reason: %v", lock, err)
		panic(err) // handle properly please!
	}

	defer lock.Unlock()

	csvfile, err := os.Create(db.registry.Configuration.UserStore)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer csvfile.Close()

	csvfile.Write([]byte("username,password,enabled,roles,claim1,claim2\n"))
	// Iterate through
	for _, v := range db.userdb {
		record := fmt.Sprintf("%v,%v,%v,%v,%v,%v\n", v.Username, v.Password, v.Enabled, strings.Join(v.Roles, ":"), v.Claim1, v.Claim2)
		csvfile.Write([]byte(record))
	}

	return nil
}

// Initiaizes data structues - IE Read roles DB
func (db *CSVReaderDatabaseInteractor) loadRoles() ([]entities.Role, error) {
	filename := db.registry.Configuration.RoleStore
	roles := make([]entities.Role, 0)

	db.registry.Logger.Log("INFO", fmt.Sprintf("Reading Roles Database %s", filename))
	// If filename is none - dont load (test usage)
	if strings.Compare("NONE", strings.ToUpper(filename)) == 0 {
		return roles, nil
	}

	csvfile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return roles, err
	}
	defer csvfile.Close()
	r := csv.NewReader(csvfile)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
		return roles, err
	}
	// Create roles
	for index, row := range records {
		if index > 0 && len(row) > 0 {
			role := entities.Role{}
			role.Name = row[roleNameField]
			// Add
			roles = append(roles, role)
		}
	}
	db.registry.Logger.Log("INFO", fmt.Sprintf("#Number of Roles = %v", len(roles)))
	return roles, nil

}

// Function loads the datastore if it has not aleady been loaded
func (db *CSVReaderDatabaseInteractor) lazyLoad() {
	if db.initialized {
		return
	}
	before := time.Now()
	db.mux.Lock()
	if db.initialized == false {
		db.initialized = true
		db.userdb, _ = db.loadUsers()
		db.roledb, _ = db.loadRoles()
	}
	db.mux.Unlock()
	now := time.Now()
	diff := now.Sub(before)
	db.registry.Logger.Log("DEBUG", fmt.Sprintf("Load user file took %v", diff))
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
