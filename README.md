<p align="center">
<img src="assets/fluxo_logo.png" height="50">
<br>
<br>
Fluxo é uma linguagem de script concatenativa.
<br>
Dê uma olhada na <a href="https://github.com/fabricioh/fluxo/wiki/Introdu%C3%A7%C3%A3o-%C3%A0-linguagem">introdução à linguagem</a>.
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

  if: {less: 10} -> (idem: {
    self: (arg | add: 1)
  })
}
```

