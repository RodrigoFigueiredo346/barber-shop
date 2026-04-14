<template>
  <div class="container" style="display:flex;flex-direction:column;justify-content:center;min-height:100vh">
    <h1 class="title">✂️ Barbearia</h1>
    <p style="text-align:center;margin-bottom:2rem;color:#aaa">Agende seu horário de forma rápida e prática</p>
    <button class="btn btn-primary" @click="$router.push('/login')" style="margin-bottom:1rem">Entrar</button>
    <button class="btn btn-secondary" @click="$router.push('/register')">Cadastrar</button>
  </div>
</template>
