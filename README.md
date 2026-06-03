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

La struttura del progetto è basata sul template ufficiale del corso [fantastic-coffee-decaffeinated](https://github.com/sapienzaapps/fantastic-coffee-decaffeinated), fornito dal professore come punto di partenza consigliato.

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
 
Le seguenti funzionalità presentano problemi noti o sono parzialmente implementate:
 
### Bug e funzionalità incomplete
- **Lista conversazioni**: la proprietà computed `orderedChats` contiene un side effect inatteso — la logica di ordinamento andrebbe spostata in un metodo separato o gestita tramite un getter puro.
- **Invio messaggi**: non è possibile inviare testo e foto contemporaneamente nello stesso messaggio; i due tipi di contenuto sono mutuamente esclusivi - sistemare visualizzazione grafica del messaggio e riattivare l'opzione già implementata nel backend.
- **Forward**: l'inoltro di un messaggio funziona solo verso conversazioni già esistenti — andrebbe esteso alla possibilità di inoltrare a qualsiasi utente, mostrando la lista degli utenti totali invece che delle chat già esistenti.
- **Reazioni**: la funzionalità di aggiunta e rimozione emoticon sui messaggi non funziona correttamente e va corretta - implemetare la funzionalità lato frontend e completarne il funzionamento nel backend per la parte del database.

### Funzionalità backend / logica
- **Autorizzazione gruppi**: verificare che un utente non possa accedere a gruppi di cui non fa parte e quindi anche che non possa visualizzarne i messaggi (al momento della creazione di questo README.md: non ricordo se l'avevo implementata del tutto, mi pare di si)
- **CORS**: la specifica richiede `Allow-All-Origins` e `Max-Age: 1` secondo (non il default). Verificare che il middleware CORS sia configurato esattamente così.

### Qualità e documentazione
- **Gestione errori**: verificare che tutti gli endpoint restituiscano i codici HTTP corretti (400, 401, 404, 500) come da specifica OpenAPI (al momento della creazione di questo README.md: ormai non ricordo se funziona tutto).
- **Popup operazioni**: creare un semplice componente di popup nella sezione vue da poter usare per mostrare correttamente notifiche quali errori(già implementato), conferme, info.

### Roadmap: progetto → prodotto

I seguenti punti raccolgono le modifiche che ho pensato di poter affrontare per trasformare questo progetto per esame universitario in un'applicazione reale, funzionante, e scalabile.

- **CORS**: ristretto `AllowedOrigins` ai soli domini autorizzati e aumentato `MaxAge` a 86400 (24 ore) per ridurre le richieste OPTIONS superflue.
- **Sistema di autenticazione**: il login semplificato tramite username senza password è una scelta intenzionale per il progetto, ma in un contesto reale lo sostituirei con un sistema di autenticazione più robusto, o delegando la gestione delle identità a un provider esterno (es. OAuth2/OpenID Connect) o implementando un flusso completo con password, token JWT e refresh token.
- **Database**: SQLite è adeguato per lo sviluppo e per un singolo nodo, ma soffre di limitazioni nella gestione della concorrenza (locking a livello di file, già causavano problemi per questo progetto, limitati con dei timer nel frontend). Per un'app in produzione con accessi simultanei andrei a sostituito con un database come **MySQL** o **PostgreSQL**, che offrono una gestione nativa delle transazioni concorrenti, migliori performance in lettura/scrittura parallela e supporto a deployment multi-istanza.
- **Architettura WebSocket**: l'attuale architettura REST richiede che il frontend faccia polling periodico per ricevere nuovi messaggi, il che introduce latenza e spreco di risorse. Per una messaggistica più immediata, si potrebbe passare ad un sistema di long polling, o per un sistema realmente istantaneo bisogna adottare un sistema basato su **WebSocket**, che mantiene una connessione persistente tra client e server e consente al server di inviare messaggi in push non appena disponibili.

