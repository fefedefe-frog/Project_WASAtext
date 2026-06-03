# Project_WASAtext

Progetto universitario per il corso **Web and Software Architecture (WASA)** — Sapienza Università di Roma.

WASAText è un'applicazione di messaggistica istantanea ispirata a WhatsApp, che consente conversazioni private e di gruppo tramite browser. Gli utenti possono scambiarsi messaggi di testo e foto, reagire ai messaggi con emoticon, gestire gruppi e personalizzare il proprio profilo.

---

## Requisiti della consegna

Il progetto è sviluppato nell'ambito dell'esame WASA e richiede di:

1. Definire le API tramite lo standard **OpenAPI**
2. Progettare e sviluppare il **backend** in **Go**
3. Progettare e sviluppare il **frontend** in **JavaScript** (Vue 3)
4. Creare un'immagine **Docker** per il deploy
### Funzionalità richieste
 
| Funzionalità | Descrizione | Rispettata |
|---|---|:---:|
| Lista conversazioni | Visualizzata in ordine cronologico inverso, con username/nome gruppo, foto, data e anteprima dell'ultimo messaggio (testo o icona foto) | ⚠️ (side effect nel computed) |
| Nuova conversazione | Avviabile con qualsiasi utente; ricerca utenti per username | ✅ |
| Gruppi | Creazione con più utenti, aggiunta membri da parte di qualsiasi membro, impossibilità di unirsi autonomamente o vedere gruppi di cui non si fa parte, possibilità di abbandonare un gruppo | ✅ |
| Visualizzazione messaggi | In ordine cronologico inverso, con timestamp, contenuto (testo o foto), username mittente per i messaggi ricevuti, checkmark per i messaggi inviati | ✅ |
| Checkmark | Uno = messaggio consegnato a tutti i destinatari; due = messaggio letto da tutti all'interno della conversazione | ✅ |
| Azioni sui messaggi | Invio, risposta (reply), inoltro (forward), eliminazione dei propri messaggi | Invio ⚠️ (no testo+foto insieme) · Reply ✅ · Forward ⚠️ (solo chat esistenti) · Elimina ✅ |
| Reazioni | Aggiunta e rimozione di emoticon su qualsiasi messaggio, con visualizzazione dei nomi degli utenti che hanno reagito | ❌ |
| Profilo utente | Login tramite username (senza password), aggiornamento username (se non già in uso), aggiornamento foto profilo | ✅ |
| Gestione gruppi | Rinomina, cambio foto, aggiunta/rimozione membri | ✅ |
| Login semplificato | Inserendo uno username si effettua login se esiste, registrazione se è nuovo — nessuna password, nessun cookie, autenticazione Bearer con identificatore utente | ✅ |
| CORS | Il backend deve rispondere alle pre-flight request con `Allow-All-Origins` e `Max-Age: 1` | ✅ |
 
---

## Struttura della repository

```
Project_WASAtext/
├── cmd/
│   └── webapi/         # Entry point del server Go (main, configurazione, avvio HTTP)
├── service/            # Logica di business del backend
│   ├── api/            # Handler delle route REST (un file per gruppo di endpoint)
│   └── database/       # Layer di accesso al database SQLite (query, modelli)
├── webui/              # Frontend Vue 3 (Single Page Application)
│   ├── src/
│   │   ├── views/      # Pagine principali (Login, Conversazioni, Chat, Profilo)
│   │   ├── components/ # Componenti riutilizzabili (MessageBubble, GroupCard, ecc.)
│   │   └── router/     # Vue Router per la navigazione client-side
│   └── ...
├── doc/                # Specifica OpenAPI (api.yaml) delle REST API
├── demo/               # Screenshot o materiale dimostrativo
├── vendor/             # Dipendenze Go vendored
├── node_modules/       # Dipendenze npm (frontend)
├── Dockerfile.backend  # Immagine Docker per il server Go
├── Dockerfile.frontend # Immagine Docker per il frontend Vue (servito via nginx)
├── docker-compose.yml  # Orchestrazione: avvia backend + frontend con un solo comando
├── go.mod / go.sum     # Moduli Go
├── package.json        # Dipendenze npm e script di build
└── .golangci.yml       # Configurazione linter Go (golangci-lint)
```

### Backend
Il backend espone una **REST API** che segue la specifica OpenAPI definita in `doc/api.yaml`. Usa **Bearer Authentication** con identificatore utente (nessuna password, nessun cookie di sessione come richiesto da consegna). Il database è **SQLite**, gestito direttamente tramite il package `database/sql`.

Gli endpoint coprono le operazioni principali dell'app: login, gestione utenti, conversazioni, messaggi, gruppi e upload di foto.

### Frontend

Per il frontend si è scelto di usare un sistema SPA costruita con **Vue 3** e comunica col backend tramite chiamate REST. Gestisce la visualizzazione delle conversazioni, l'invio di messaggi, la gestione dei gruppi e le impostazioni del profilo utente. Il frontend è containerizzato e servito tramite **nginx**.

### Deploy

L'intera applicazione si avvia con:

```bash
docker compose up
```

Backend disponibile su `localhost:3000`, frontend su `localhost:8080`.

---

## API implementate

Tutti gli `operationId` richiesti dalla consegna sono presenti nella specifica OpenAPI:

| operationId | Descrizione |
|---|---|
| `doLogin` | Login / registrazione tramite username |
| `setMyUserName` | Aggiornamento username |
| `setMyPhoto` | Aggiornamento foto profilo |
| `getMyConversations` | Lista conversazioni dell'utente |
| `getConversation` | Messaggi di una singola conversazione |
| `sendMessage` | Invio messaggio (testo o foto) |
| `forwardMessage` | Inoltro messaggio |
| `deleteMessage` | Eliminazione messaggio |
| `commentMessage` | Aggiunta reazione (emoticon) |
| `uncommentMessage` | Rimozione reazione |
| `addToGroup` | Aggiunta utente a un gruppo |
| `leaveGroup` | Abbandono di un gruppo |
| `setGroupName` | Rinomina gruppo |
| `setGroupPhoto` | Aggiornamento foto gruppo |

---

## Sviluppi futuri / possibili miglioramenti

Le seguenti funzionalità sono richieste dalla specifica ma potrebbero essere incomplete o assenti:

### Funzionalità frontend
- **Reply a un messaggio**: la consegna richiede di poter rispondere a un messaggio specifico con riferimento visivo al messaggio originale.
- **Preview corretta nella lista conversazioni**: ogni elemento deve mostrare testo del messaggio (snippet) oppure un'icona per messaggi foto — distinguere i due casi visivamente.
- **Ordine cronologico inverso**: sia la lista conversazioni che i messaggi all'interno di una chat devono essere ordinati dal più recente al più vecchio.
- **Checkmark stato messaggio**: un checkmark = messaggio consegnato a tutti; due checkmark = messaggio letto da tutti. La logica di tracking "letto da tutti" nei gruppi è complessa.

### Funzionalità backend / logica
- **Autorizzazione gruppi**: verificare che un utente non possa accedere a gruppi di cui non fa parte (né visualizzarli né cercarne i messaggi).
- **CORS**: la specifica richiede `Allow-All-Origins` e `Max-Age: 1` secondo (non il default). Verificare che il middleware CORS sia configurato esattamente così.

### Qualità e documentazione
- **README**: attualmente vuoto — aggiungere istruzioni di avvio, requisiti e descrizione del progetto.
- **Gestione errori**: verificare che tutti gli endpoint restituiscano i codici HTTP corretti (400, 401, 404, 500) come da specifica OpenAPI.
- **Validazione input**: il campo `name` per il login richiede lunghezza 3–16 caratteri — verificare che la validazione sia applicata sia lato frontend che backend.
