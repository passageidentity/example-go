# Using Passage with Go

Passage provides an SDK to easily authenticate HTTP requests. Example source code can be found on GitHub, [here](https://github.com/passageidentity/example-go).

### Setup

For this example app, we'll need to provide our application with a Passage App ID and API Key. Your App ID and API Key can be found in the [Passage Console](https://console.passage.id) in your App Settings. You'll need to change the following environment variables with your respective credentials (note that the default port for this application is 5000):

```go
	os.Setenv("PASSAGE_APP_ID", "[YOUR_APP_ID_HERE]")
	os.Setenv("PASSAGE_API_KEY", "[YOUR_PASSAGE_API_KEY_HERE]")
	os.Setenv("PORT", "5000")
```

### Run With Go

To run this example app, make sure you have [Go installed on your computer](https://golang.org/doc/install).

Run the following command:

```bash
go run main.go
```

### Run With Docker

Create your docker image with the following command:

```bash
$ docker build -t example-go .
```

Run your docker container using the example-go image:

```bash
$ docker run -p 5000:5000 example-go
```

### Authenticating an HTTP Request

A Go server can easily authenticate an HTTP request using the Passage SDK, as shown below.

```go
import (
	"net/http"

	"github.com/passageidentity/passage-go"
)

func exampleHandler(w http.ResponseWriter, r *http.Request) {

	// Authenticate this request using the Passage SDK.
	psg := passage.New("<Passage App Handle>")
	_, err := psg.AuthenticateRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Successful authentication. Proceed...

}
```

### Authorizing a User

It is important to remember that the `psg.AuthenticateRequest()` function validates that a request is properly authenticated, but makes no assertions about _who_ it is authorized for. To perform an authorization check, the Passage User Handle can be referenced.

```go
func exampleHandler(w http.ResponseWriter, r *http.Request) {

	// Authenticate this request using the Passage SDK.
	psg := passage.New("<Passage App Handle>")
	passageHandle, err := psg.AuthenticateRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Check that the user with `passageHandle` is allowed to perform
	// a certain action on a certain resource.
	err = authorizationCheck(passageHandle)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Successful authentication AND authorization. Proceed...

}
```

## Get User

To get user information, you can use the Passage SDK with an API key. This will authenticate your web server to Passage and grant you management
access to user information. API keys should never be hard-coded in source code, but stored in environment variables or a secrets storage mechanism.

```go

	psg := passage.New("<Passage App Handle>", "<API_KEY>")
	passageHandle, err := psg.AuthenticateRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

    //information regarding the user will exist in the user variable
	user, err := psg.GetUser(passageHandle)
	if err != nil {
		fmt.Println("Could not get user: ", err)
		return
	}
```

### Adding Authentication to the Frontend

The easiest way to add authentication to a web frontend is with a Passage Element. The HTML below will automatically embed a complete UI/UX for user sign-in and sign-up.

```html
<!-- Passage will populate this custom element with a complete authentication UI/UX. -->
<passage-auth app-id="<Passage App Handle>"></passage-auth>

<!-- Include the passage-web JavaScript from the Passage CDN. -->
<script src="https://cdn.passage.id/passage-web.js"></script>
```
