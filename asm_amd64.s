
#define get_tls(r)      MOVL TLS, r
#define g(r)    0(r)(TLS*1)
TEXT ·GetG(SB),NOSPLIT,$0-8
  get_tls(CX)
  MOVQ  g(CX), AX
  MOVQ  AX, ret+0(FP)
  RET
