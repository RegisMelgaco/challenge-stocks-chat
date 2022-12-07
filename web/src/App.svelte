<script>
  let authToken = "";
  let username = "";
  let pass = "";
  let msgContent = "";
  const chat = [];

  let socket;

  async function handleCreateUser(e) {
    e.preventDefault();

    const res = await fetch("http://localhost:3000/auth/signup", {
      method: "POST",
      body: JSON.stringify({ username, pass }),
      mode: "cors",
      headers: {
        Origin: "localhost:5137",
      },
    });

    const json = await res.json();

    if (res.ok) {
      alert("success: " + JSON.stringify(json));
    } else {
      alert("failed: " + JSON.stringify(json));
    }
  }

  async function handleLogin(e) {
    e.preventDefault();

    const res = await fetch("http://localhost:3000/auth/login", {
      method: "POST",
      body: JSON.stringify({ username, pass }),
      mode: "cors",
      headers: {
        Origin: "localhost:5137",
      },
    });

    const json = await res.json();

    if (res.ok) {
      authToken = json.token;

      const ws = new WebSocket("ws://localhost:3000/chat/listen");

      ws.onmessage = (event) => {
        chat.push(event.data)

        console.log(chat)
      };

      ws.onopen = () => {
        ws.send(JSON.stringify({token: authToken}))
      
        socket = ws
      }

      ws.onerror = (err) => {
        console.log(err)
      }

      alert("success: " + JSON.stringify(json));
    } else {
      alert("failed: " + JSON.stringify(json));
    }
  }

  async function handleCreateMessage(e) {
    e.preventDefault();

    socket.send(JSON.stringify({content: msgContent}))
  }
</script>

<main>
  <nav class="navbar bg-light">
    <div class="container container-fluid">
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

  <section>
    {#each chat as msg }
      <p>{msg.author}: {msg.content}</p>
    {/each}
  </section>
</main>

<style>
</style>
