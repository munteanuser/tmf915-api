package router

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/tmf915-api/internal/handlers"
	"github.com/tmf915-api/internal/middleware"
)

const swaggerHTML = `<!DOCTYPE html>
<html>
<head>
  <title>TMF915 AI Management API</title>
  <meta charset="utf-8"/>
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    // Load the spec from /swagger.json (which is already host-patched by the server)
    fetch('/swagger.json')
      .then(function(r){ return r.json(); })
      .then(function(spec){
        SwaggerUIBundle({
          spec: spec,
          dom_id: '#swagger-ui',
          presets: [SwaggerUIBundle.presets.apis, SwaggerUIBundle.SwaggerUIStandalonePreset],
          layout: "BaseLayout",
          deepLinking: true,
          persistAuthorization: true,
          requestInterceptor: function(request) {
            // Inject stored API key automatically
            var key = localStorage.getItem('tmf915_api_key');
            if (key) request.headers['X-API-Key'] = key;
            return request;
          }
        });
      });
  </script>
  <style>
    /* Persistent API-key helper bar */
    #apikey-bar { background:#1b1b1b; color:#fff; padding:8px 16px; font-family:monospace; font-size:13px; display:flex; align-items:center; gap:12px; }
    #apikey-bar input { flex:1; max-width:420px; padding:4px 8px; border-radius:4px; border:none; font-size:13px; }
    #apikey-bar button { padding:4px 12px; border-radius:4px; border:none; cursor:pointer; background:#4CAF50; color:#fff; font-size:13px; }
  </style>
  <div id="apikey-bar">
    <span>X-API-Key:</span>
    <input id="apikey-input" type="text" placeholder="paste your API key here"/>
    <button onclick="var v=document.getElementById('apikey-input').value; localStorage.setItem('tmf915_api_key',v); alert('Key saved — all Try-it-out requests will include it.')">Save</button>
    <span id="apikey-status" style="color:#aaa;font-size:12px"></span>
  </div>
  <script>
    var stored = localStorage.getItem('tmf915_api_key');
    if (stored) {
      document.getElementById('apikey-input').value = stored;
      document.getElementById('apikey-status').textContent = '(key loaded from storage)';
    }
  </script>
</body>
</html>`

func New(
	log zerolog.Logger,
	apiKey string,
	swaggerPath string,
	aiModelSpec *handlers.AiModelSpecificationHandler,
	aiModel *handlers.AiModelHandler,
	aiContractSpec *handlers.AiContractSpecificationHandler,
	aiContract *handlers.AiContractHandler,
	aiContractViolation *handlers.AiContractViolationHandler,
	alarm *handlers.AlarmHandler,
	rule *handlers.RuleHandler,
	monitor *handlers.MonitorHandler,
	topic *handlers.TopicHandler,
	event *handlers.EventHandler,
	hub *handlers.HubHandler,
) http.Handler {
	mux := http.NewServeMux()

	// Swagger UI (CDN-based, no binary assets needed)
	mux.HandleFunc("GET /docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(swaggerHTML))
	})
	mux.HandleFunc("GET /swagger.json", func(w http.ResponseWriter, r *http.Request) {
		raw, err := os.ReadFile(swaggerPath)
		if err != nil {
			http.Error(w, "swagger.json not found", http.StatusNotFound)
			return
		}
		// Patch host and scheme so Swagger UI "Try it out" hits this server
		var spec map[string]interface{}
		if err := json.Unmarshal(raw, &spec); err != nil {
			http.Error(w, "invalid swagger.json", http.StatusInternalServerError)
			return
		}
		spec["host"] = r.Host
		spec["schemes"] = []string{"http"}
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		json.NewEncoder(w).Encode(spec)
	})

	base := "/tmf-api/AiM/v4"

	// AiModelSpecification
	mux.HandleFunc("GET "+base+"/aiModelSpecification", aiModelSpec.List)
	mux.HandleFunc("POST "+base+"/aiModelSpecification", aiModelSpec.Create)
	mux.HandleFunc("GET "+base+"/aiModelSpecification/{id}", aiModelSpec.Get)
	mux.HandleFunc("PATCH "+base+"/aiModelSpecification/{id}", aiModelSpec.Update)
	mux.HandleFunc("DELETE "+base+"/aiModelSpecification/{id}", aiModelSpec.Delete)

	// AiModel
	mux.HandleFunc("GET "+base+"/aiModel", aiModel.List)
	mux.HandleFunc("POST "+base+"/aiModel", aiModel.Create)
	mux.HandleFunc("GET "+base+"/aiModel/{id}", aiModel.Get)
	mux.HandleFunc("PATCH "+base+"/aiModel/{id}", aiModel.Update)
	mux.HandleFunc("DELETE "+base+"/aiModel/{id}", aiModel.Delete)

	// AiContractSpecification
	mux.HandleFunc("GET "+base+"/aiContractSpecification", aiContractSpec.List)
	mux.HandleFunc("POST "+base+"/aiContractSpecification", aiContractSpec.Create)
	mux.HandleFunc("GET "+base+"/aiContractSpecification/{id}", aiContractSpec.Get)
	mux.HandleFunc("PATCH "+base+"/aiContractSpecification/{id}", aiContractSpec.Update)
	mux.HandleFunc("DELETE "+base+"/aiContractSpecification/{id}", aiContractSpec.Delete)

	// AiContract
	mux.HandleFunc("GET "+base+"/aiContract", aiContract.List)
	mux.HandleFunc("POST "+base+"/aiContract", aiContract.Create)
	mux.HandleFunc("GET "+base+"/aiContract/{id}", aiContract.Get)
	mux.HandleFunc("PATCH "+base+"/aiContract/{id}", aiContract.Update)
	mux.HandleFunc("DELETE "+base+"/aiContract/{id}", aiContract.Delete)

	// AiContractViolation
	mux.HandleFunc("GET "+base+"/aiContractViolation", aiContractViolation.List)
	mux.HandleFunc("POST "+base+"/aiContractViolation", aiContractViolation.Create)
	mux.HandleFunc("GET "+base+"/aiContractViolation/{id}", aiContractViolation.Get)

	// Alarm
	mux.HandleFunc("GET "+base+"/alarm", alarm.List)
	mux.HandleFunc("POST "+base+"/alarm", alarm.Create)
	mux.HandleFunc("GET "+base+"/alarm/{id}", alarm.Get)
	mux.HandleFunc("PATCH "+base+"/alarm/{id}", alarm.Update)
	mux.HandleFunc("DELETE "+base+"/alarm/{id}", alarm.Delete)

	// Rule
	mux.HandleFunc("GET "+base+"/rule", rule.List)
	mux.HandleFunc("POST "+base+"/rule", rule.Create)
	mux.HandleFunc("GET "+base+"/rule/{id}", rule.Get)
	mux.HandleFunc("PATCH "+base+"/rule/{id}", rule.Update)
	mux.HandleFunc("DELETE "+base+"/rule/{id}", rule.Delete)

	// Monitor
	mux.HandleFunc("GET "+base+"/monitor", monitor.List)
	mux.HandleFunc("POST "+base+"/monitor", monitor.Create)
	mux.HandleFunc("GET "+base+"/monitor/{id}", monitor.Get)
	mux.HandleFunc("PATCH "+base+"/monitor/{id}", monitor.Update)
	mux.HandleFunc("DELETE "+base+"/monitor/{id}", monitor.Delete)

	// Topic
	mux.HandleFunc("GET "+base+"/topic", topic.List)
	mux.HandleFunc("POST "+base+"/topic", topic.Create)
	mux.HandleFunc("GET "+base+"/topic/{id}", topic.Get)
	mux.HandleFunc("PATCH "+base+"/topic/{id}", topic.Update)
	mux.HandleFunc("DELETE "+base+"/topic/{id}", topic.Delete)

	// Events (nested under topic)
	mux.HandleFunc("GET "+base+"/topic/{topicId}/event", event.List)
	mux.HandleFunc("GET "+base+"/topic/{topicId}/event/{id}", event.Get)

	// Hub
	mux.HandleFunc("GET "+base+"/hub", hub.ListHubs)
	mux.HandleFunc("POST "+base+"/hub", hub.CreateHub)
	mux.HandleFunc("GET "+base+"/hub/{id}", hub.GetHub)
	mux.HandleFunc("DELETE "+base+"/hub/{id}", hub.DeleteHub)

	// Listener
	mux.HandleFunc("POST "+base+"/listener", hub.RegisterListener)
	mux.HandleFunc("DELETE "+base+"/listener/{id}", hub.UnregisterListener)

	// Management SPA (no auth — API key entered inside the app)
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("spa"))))
	mux.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/app/", http.StatusMovedPermanently)
	})

	// Health check (no auth)
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Chain middleware: recovery → logging → api key
	var h http.Handler = mux
	h = middleware.Recovery(log)(h)
	h = middleware.Logging(log)(h)
	h = middleware.APIKey(apiKey)(h)

	return h
}
