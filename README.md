<p align="center">
  <img src="assets/fluxo_logo.png" height="40">
</p>

Fluxo é uma linguagem de script concatenativa.

```
do: 0 -> * {
  arg
  println

  if: {less: 10} -> {
    self: (add: 1)
  }
}

; printa números de 0 a 10
```
