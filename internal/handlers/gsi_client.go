package handlers

import "net/http"

// gisClientScript is a minimal mock of the Google Identity Services browser
// client (https://accounts.google.com/gsi/client), sufficient to drive local
// development/test credential flows. It is not a replacement for the real SDK.
const gisClientScript = `(function () {
  "use strict";

  function requestCredential(clientID) {
    return fetch("/gsi/credential", {
      method: "POST",
      credentials: "same-origin",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ client_id: clientID }),
    })
      .then(function (res) {
        if (!res.ok) {
          throw new Error("mock GIS credential request failed: " + res.status);
        }
        return res.json();
      })
      .then(function (data) {
        return data.credential;
      });
  }

  function renderSignInButton(container, clientID, callbackName) {
    var button = document.createElement("button");
    button.type = "button";
    button.textContent = "Sign in (mock GIS)";
    button.setAttribute("aria-label", "Sign in with mock Google Identity Services");
    button.addEventListener("click", function () {
      requestCredential(clientID).then(function (credential) {
        var callback = window[callbackName];
        if (typeof callback === "function") {
          callback({ credential: credential });
        }
      });
    });
    container.appendChild(button);
  }

  function init() {
    var onload = document.getElementById("g_id_onload");
    if (!onload) {
      return;
    }

    var clientID = onload.getAttribute("data-client_id");
    var callbackName = onload.getAttribute("data-callback");
    if (!clientID || !callbackName) {
      return;
    }

    var containers = document.querySelectorAll(".g_id_signin");
    for (var i = 0; i < containers.length; i++) {
      renderSignInButton(containers[i], clientID, callbackName);
    }
  }

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", init);
  } else {
    init();
  }
})();
`

// GISClientHandler serves the mock Google Identity Services browser compatibility script.
type GISClientHandler struct{}

// NewGISClientHandler creates a new GISClientHandler
func NewGISClientHandler() *GISClientHandler {
	return &GISClientHandler{}
}

// ServeHTTP handles requests for the mock GIS client script
func (h *GISClientHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(gisClientScript))
}
