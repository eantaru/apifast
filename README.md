# APIFast Go Library

This library provides a simple API client for building HTTP requests with the `fasthttp` library in Go. This library supports GET, POST, PATCH, and DELETE requests with optional Basic Authentication, Bearer Token Authentication, custom headers, and more.


## Table of Contents
1. [Installation](#installation)
2. [Basic Usage Example (GET Request)](#basic-usage-example-get-request)
3. [Making a POST Request](#making-a-post-request)
4. [Using Basic Authentication](#using-basic-authentication)
5. [Using Bearer Token Authentication](#using-bearer-token-authentication)
6. [Adding Custom Headers](#adding-custom-headers)


## Installation

To install the package, run:

```bash
go get github.com/eantaru/apifast
```


## Usage

Basic Usage Example (GET Request)

```go
package main

import (
    "fmt"
    "time"
    "github.com/eantaru/apifast"
)

func main() {
    response, err := apifast.Build().
        Uri("https://jsonplaceholder.typicode.com/posts").
        Timeout(5 * time.Second). // Optional timeout
        Get()

    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Println("Response Code:", response.Code)
    fmt.Println("Response Body:", string(response.Body.([]byte)))
}
```

## Making a POST Request

```go
package main

import (
    "fmt"
    "time"
    "github.com/eantaru/apifast"
)

func main() {
    payload := []byte(`{"title": "foo", "body": "bar", "userId": 1}`)

    response, err := apifast.Build().
        Uri("https://jsonplaceholder.typicode.com/posts").
        Timeout(5 * time.Second).
        Payload(payload).
        Post()

    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Println("Response Code:", response.Code)
    fmt.Println("Response Body:", string(response.Body.([]byte)))
}

```


### Using Basic Authentication
If the API requires Basic Authentication, you can provide the username and password like this:


```go
package main

import (
    "fmt"
    "time"
    "github.com/eantaru/apifast"
)

func main() {
    auth := apifast.Auth{
        Username: "myUsername",
        Password: "myPassword",
    }

    response, err := apifast.Build().
        Uri("https://example.com/api").
        Timeout(5 * time.Second).
        Auth(auth).
        Get()

    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Println("Response Code:", response.Code)
    fmt.Println("Response Body:", string(response.Body.([]byte)))
}
```

### Using Bearer Token Authentication
For APIs requiring Bearer Token authentication, you can provide the token like this:

```go
package main

import (
    "fmt"
    "time"
    "github.com/eantaru/apifast"
)

func main() {
    auth := apifast.Auth{
        Token: "your-bearer-token",
    }

    response, err := apifast.Build().
        Uri("https://example.com/api").
        Timeout(5 * time.Second).
        Auth(auth).
        Get()

    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Println("Response Code:", response.Code)
    fmt.Println("Response Body:", string(response.Body.([]byte)))
}
```


### Adding Custom Headers
You can pass custom headers with the request by providing a list of headers:


```go
package main

import (
    "fmt"
    "time"
    "github.com/eantaru/apifast"
)

func main() {
    headers := []apifast.Header{
        {Tag: "Content-Type", Value: "application/json"},
        {Tag: "X-Custom-Header", Value: "CustomValue"},
    }

    response, err := apifast.Build().
        Uri("https://example.com/api").
        Timeout(5 * time.Second).
        Headers(headers).
        Get()

    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Println("Response Code:", response.Code)
    fmt.Println("Response Body:", string(response.Body.([]byte)))
}

```