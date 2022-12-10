<script>
  let username = "";
  let pass = "";
  let msgContent = "";
  let chat = [];
  $: sortedChat = chat.sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())

  let socket = null;

  const mode = import.meta.env.DEV ? "cors" : "same-origin";
  const headers = import.meta.env.DEV ? {
      Origin: import.meta.env.BASE_URL,
  } : {};
  const origin = import.meta.env.DEV ? "localhost:3000" : "localhost:8080/api"

  async function handleCreateUser(e) {
    e.preventDefault();

    const res = await fetch(`http://${origin}/auth/signup`, {
      method: "POST",
      body: JSON.stringify({ username, pass }),
      mode, headers,
    });

    const json = await res.json();

    if (!res.ok) {
      alert("failed: " + JSON.stringify(json));

      return;
    } 

    alert("success: " + JSON.stringify(json));
  }

  async function handleLogin(e) {
    e.preventDefault();

    const res = await fetch(`http://${origin}/auth/login`, {
      method: "POST",
      body: JSON.stringify({ username, pass }),
      mode: "cors",
      headers: {
        Origin: "localhost:5137",
      },
    });

    const json = await res.json();

    if (!res.ok) {
      alert("failed: " + JSON.stringify(json));

      return;
    }



    const ws = new WebSocket(`ws://${origin}/chat/listen`);

    ws.onmessage = async (event) => {
      const json = await JSON.parse(event.data);

      if (json.type !== "message") {
        alert(json.payload);

        return;
      }

      chat = [...chat, json.payload];
    };

    ws.onopen = () => {
      ws.send(JSON.stringify({ token: json.token }));

      socket = ws;
    };

    ws.onerror = (err) => {
      alert(err);
    };

    alert("success: " + JSON.stringify(json));
  }

  async function handleCreateMessage(e) {
    e.preventDefault();

    if (socket === null) {
      alert("login is required to send messages")

      return;
    }

    socket.send(JSON.stringify({ content: msgContent }));
  }
</script>

<main>
  <nav class="navbar bg-light">
    <div class="container container-fluid navbar-dark bg-cyan-200">
      <a class="navbar-brand">Stocks Chat</a>
    </div>
  </nav>

  <section class="container mt-4">
    <h3>User</h3>
    <form class="d-flex">
      <input
        class="form-control me-2"
        placeholder="Username"
        aria-label="Username"
        bind:value={username}
      />
      <input
        class="form-control me-2"
        type="password"
        placeholder="Password"
        aria-label="Password"
        bind:value={pass}
      />

      <button class="btn btn-success me-2" on:click={handleCreateUser}
        >Create</button
      >
      <button class="btn btn-success" on:click={handleLogin}>Login</button>
    </form>
  </section>

  <section class="container mt-4">
    <h3>Chat</h3>
    <form class="d-flex">
      <input
        class="form-control me-2"
        placeholder="message"
        aria-label="message"
        bind:value={msgContent}
      />

      <button class="btn btn-success me-2" on:click={handleCreateMessage}
        >Create</button
      >
    </form>
  </section>

  <section class="container mt-4">
    {#each sortedChat as msg (msg.created_at)}
      <div class="card mb-2">
        <div class="card-body">
          <h5 class="card-title">{msg.author}</h5>
          {msg.content}
        </div>
      </div>
    {/each}
  </section>
</main>

<style>
</style>
