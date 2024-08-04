package myapp

import (
	"testing"
	"net/http/httptest"
	"net/http"

	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strings"
	"encoding/json"
	"strconv"
	"fmt"
)

func TestIndex(t *testing.T){
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	data, _ := ioutil.ReadAll(resp.Body)
	assert.Equal("Hello World", string(data))
}

//위에랑 똑같은데 경로만 다른것임
func TestUsers(t *testing.T){
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/users")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	data, _ := ioutil.ReadAll(resp.Body)
	assert.Contains(string(data), "No Users")
}

//restful api에서 get id 하는거
func TestGettUserInfo(t *testing.T){
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/users/89")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ := ioutil.ReadAll(resp.Body)
	assert.Contains(string(data), "No User Id:89")
	
}

//post
func TestCreateUser(t *testing.T){
	assert := assert.New(t)

	//test server
	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	resp, err := http.Post(ts.URL + "/users", "application/json",
		strings.NewReader(`{"first_name":"tucker", "last_name":"kim", "email":"tucker@naver.com"}`))

	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	user := new(User)
	err = json.NewDecoder(resp.Body).Decode(user)
	assert.NoError(err)
	assert.NotEqual(0, user.ID)

	id := user.ID
	resp, err =	http.Get(ts.URL + "/users/" + strconv.Itoa(id))
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	
	user2 := new(User)
	err = json.NewDecoder(resp.Body).Decode(user2)
	assert.NoError(err)
	assert.Equal(user.ID, user2.ID)
}

func TestDeleteUser(t *testing.T){
	assert := assert.New(t)

	//test server
	ts := httptest.NewServer(NewHandler())
	defer ts.Close()
																		//body
	req, _ := http.NewRequest("DELETE", ts.URL + "/users/1", nil)
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	data, _ := ioutil.ReadAll(resp.Body)
	assert.Contains(string(data), "No User Id:1")
	
	//405 error : no handelr
	// data, _ := ioutil.ReadAll(resp.Body)
	// log.Print(string(data))

	resp, err = http.Post(ts.URL + "/users", "application/json",
	strings.NewReader(`{"first_name":"tucker", "last_name":"kim", "email":"tucker@naver.com"}`))

	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	user := new(User)
	err = json.NewDecoder(resp.Body).Decode(user)
	assert.NoError(err)
	assert.NotEqual(0, user.ID)

	req, _ = http.NewRequest("DELETE", ts.URL + "/users/1", nil)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ = ioutil.ReadAll(resp.Body)
	assert.Contains(string(data), "Deleted User Id:1")
}

func TestUpdateUser(t *testing.T){
	assert := assert.New(t)

	//test server
	ts := httptest.NewServer(NewHandler())
	defer ts.Close()
																		//body
	req, _ := http.NewRequest("PUT", ts.URL + "/users", 
		strings.NewReader(`{"id":1, "first_name":"updated", "last_name":"updated", "email":"updated@naver.com"}`))

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	// data, _ := ioutil.ReadAll(resp.Body)
	// assert.Contains(string(data), "No User Id:1")


	resp, err = http.Post(ts.URL + "/users", "application/json",
	strings.NewReader(`{"first_name":"tucker", "last_name":"kim", "email":"tucker@naver.com"}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	user := new(User)
	err = json.NewDecoder(resp.Body).Decode(user)
	assert.NoError(err)
	assert.NotEqual(0, user.ID)

	//일케 보내면 스트럭쳐에서 보낸것만 받고 나머지는 안감 ;; 그래서 걍  default로 바꿈
	updateStr := fmt.Sprintf(`{"id":%d, "first_name":"jason", "last_name":""}`, user.ID)

	req, _ = http.NewRequest("PUT", ts.URL + "/users", 
	strings.NewReader(updateStr))

	resp, err = http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	updateUser := new(User)
	err = json.NewDecoder(resp.Body).Decode(updateUser)
	assert.NoError(err)
	assert.Equal(updateUser.ID, user.ID)
	assert.Equal("jason", updateUser.FirstName)
//	assert.Equal("", updateUser.LastName)
}


//restful api에서 get id 하는거
func TestUsers_WithUsersData(t *testing.T){
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()


	resp, err := http.Post(ts.URL + "/users", "application/json",
	strings.NewReader(`{"first_name":"tucker", "last_name":"kim", "email":"tucker@naver.com"}`))

	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)


	resp, err = http.Post(ts.URL + "/users", "application/json",
	strings.NewReader(`{"first_name":"jason", "last_name":"park", "email":"jay@naver.com"}`))

	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	resp, err =	http.Get(ts.URL+"/users")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	users := []*User{}
	err = json.NewDecoder(resp.Body).Decode(&users)
	assert.NoError(err)
	assert.Equal(2, len(users))
}