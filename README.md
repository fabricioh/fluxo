<p align="center">
<img src="assets/fluxo_logo.png" height="50">
<br>
<br>
Fluxo é uma linguagem de script concatenativa.
<br>
Dê uma olhada na <a href="https://github.com/fabricioh/fluxo/wiki/Introdu%C3%A7%C3%A3o-%C3%A0-linguagem">Introdução à linguagem</a>. Depois veja a <a href="https://github.com/fabricioh/fluxo/wiki/Documenta%C3%A7%C3%A3o">Documentação</a>.
<br>
Faça o download <a href="https://github.com/fabricioh/fluxo/releases">aqui</a>.
<br>
<br>
</p>

```
; Uma função anônima recursiva que
; imprime números de 0 a 10

do: 0 -> * {
  arg
  println

  if: {less: 10} -> {
    self: (add: 1)
  }
}
```

## Utilização

Basta baixar o zip [aqui](https://github.com/fabricioh/fluxo/releases), extrair e colocar o executável na variável de ambiente Path.

O seguinte programa é um Hello World em fluxo:

```
val: "Hello, World!" | println
```

Para incluir no seu script o arquivo `lib.fl` que vem com o release, basta usar a função `exec` passando o caminho para o arquivo como o primeiro valor de uma lista:

```
exec: ["lib.fl"]

val: ["hello" "fluxo!"]
each: {println}

; hello
; fluxo!
```