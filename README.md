# Using Passage with Go

Passage provides an SDK to easily authenticate HTTP requests. Example source code can be found on GitHub, [here](https://github.com/passageidentity/example-go).

### Configuring a Go Server

The `passage-go` SDK depends on a `PASSAGE_PUBLIC_KEY` environment variable being set. An app's `PASSAGE_PUBLIC_KEY` can be copied off of the Passage Console.

### Authenticating an HTTP Request

A Go server can easily authenticate an HTTP request using the Passage SDK, as shown below.

```go
import (
	"net/http"

	"github.com/passageidentity/passage-go"
)

func exampleHandler(w http.ResponseWriter, r *http.Request) {

	// Authenticate this request using the Passage SDK.
	psg := passage.New()
	_, err := psg.AuthenticateRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Successful authentication. Proceed...

}
```

### Authorizing a User

It is important to remember that the `psg.AuthenticateRequest()` function validates that a request is properly authenticated, but makes no assertions about *who* it is authorized for. To perform an authorization check, the Passage User Handle can be referenced.

```go
func exampleHandler(w http.ResponseWriter, r *http.Request) {

	// Authenticate this request using the Passage SDK.
	psg := passage.New()
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

### Adding Authentication to the Frontend

The easiest way to add authentication to a web frontend is with a Passage Element. The HTML below will automatically embed a complete UI/UX for user sign-in and sign-up.

```html
<!-- Passage will populate this div with a complete authentication UI/UX. -->
<div id="passage-auth" data-app="<Passage App Handle>"></div>

<!-- Include the passage-web JavaScript from the Passage CDN. -->
<script src="https://cdn.passage.id/passage-web.js"></script>
```