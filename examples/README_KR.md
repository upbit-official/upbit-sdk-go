[English](./README.md) | 한국어

# 예제 코드

이 디렉토리에는 Upbit Go SDK의 주요 기능을 시나리오 형태로 보여주는 예제가 포함되어 있습니다.

이 디렉토리에는 시나리오 예제와 API별 생성 예제가 있습니다.
`DRY_RUN`을 구현한 시나리오 예제는 기본 실행에서 해당 분기의 쓰기 요청을 생략합니다. 조회와 테스트 주문 API는 실제로 호출합니다.
[`orders/`](./orders/) 같은 기능별 하위 디렉토리의 생성 예제에는 `DRY_RUN`이 없습니다. 실행하면 해당 API를 호출하며, 실제 주문·출금 요청도 보낼 수 있습니다.

> **참고:** 각 파일은 독립된 `package main` 프로그램입니다. `go run examples/<file>.go`로 실행하세요. API 엔드포인트 단위의 예제는 기능별 하위 디렉토리에 정리되어 있습니다. 예를 들어 주문 관련 예제는 [`examples/orders/`](./orders/)를 참고하세요.

---

## 빠르게 시작하기

처음이라면 아래 순서로 살펴보는 것을 추천합니다.

1. [`quotation_kr.go`](./quotation_kr.go) — 인증 없이 시세/캔들/체결/호가 예제를 빠르게 확인
2. [`indicators_kr.go`](./indicators_kr.go) — Quotation 데이터를 활용해 투자 지표 계산 흐름 이해
3. [`orders_kr.go`](./orders_kr.go) — 주문 생성, 조회, 취소 흐름을 Dry run으로 먼저 확인

### 인증이 필요 없는 예제

아래 예제는 API 키 없이 바로 실행할 수 있습니다.
- [`quotation_kr.go`](./quotation_kr.go)
- [`indicators_kr.go`](./indicators_kr.go)

### 인증이 필요한 예제

아래 예제는 Upbit API 키가 필요합니다.
- [`orders_kr.go`](./orders_kr.go)
- [`orders_validate_kr.go`](./orders_validate_kr.go)
- [`deposits_kr.go`](./deposits_kr.go)
- [`withdrawals_kr.go`](./withdrawals_kr.go)
- [`dca_kr.go`](./dca_kr.go)
- [`tp_sl_kr.go`](./tp_sl_kr.go)

### 사전 준비

```bash
# Upbit API 키 발급 후 환경변수 설정
export UPBIT_ACCESS_KEY=<your-access-key>
export UPBIT_SECRET_KEY=<your-secret-key>
```

### 실행 전 알아두기

- **조회 중심으로 먼저 확인하고 싶다면** `quotation_kr.go`, `indicators_kr.go`부터 시작하세요.
- **실제 주문/출금이 가능한 예제를 볼 때는** 먼저 Dry run 모드로 동작 방식을 확인하세요.
- **`DRY_RUN=false`로 실행하는 경우** 주문, 출금, 자동 매매 관련 쓰기 작업이 실제로 발생할 수 있습니다.

---

## 예제 한눈에 보기

| 예제 | 용도 | 인증 | 기본 실행 특성 |
| --- | --- | --- | --- |
| [`quotation_kr.go`](./quotation_kr.go) | 시세/캔들/체결/호가 조회 | 불필요 | 안전한 조회 전용 |
| [`indicators_kr.go`](./indicators_kr.go) | RSI 등 투자 지표 계산 | 불필요 | 안전한 조회 전용 |
| [`orders_kr.go`](./orders_kr.go) | 주문 생성·조회·취소 흐름 | 필요 | 기본은 Dry run, 실제 주문 가능 |
| [`orders_validate_kr.go`](./orders_validate_kr.go) | 주문 생성 테스트 API로 주문 유형 전체 검증 (KRW) | 필요 | 항상 안전 — 실제 주문 없음 |
| [`deposits_kr.go`](./deposits_kr.go) | 입금 주소/내역 관리 | 필요 | 기본은 Dry run, 일부 작업은 실반영 가능 |
| [`withdrawals_kr.go`](./withdrawals_kr.go) | 출금 가능 정보/출금 흐름 | 필요 | 기본은 Dry run, 실제 출금은 주의 필요 |
| [`dca_kr.go`](./dca_kr.go) | 정기적 시장가 매수 자동화 | 필요 | 기본은 Dry run, 실제 주문 가능 |
| [`tp_sl_kr.go`](./tp_sl_kr.go) | 익절/손절 자동 매도 | 필요 | 기본은 Dry run, 실제 주문 가능 |

---

## 시나리오 예제 상세

### `quotation_kr.go` — 시세 조회 입문

인증 없이 사용할 수 있는 Quotation API 전반을 빠르게 익히기 좋은 예제입니다.
처음 SDK를 살펴보는 경우 가장 먼저 실행해 보기 좋습니다.

**포함 기능:**
- 페어(마켓) 목록 조회 및 유의 종목 필터링
- 현재가(Ticker) 조회 — 단일 / 복수 / 마켓 단위
- 캔들(OHLCV) 조회 — 분봉(5분), 일봉, 주봉, 월봉 / `to` 파라미터 활용
- 최근 체결 내역 조회 — `days_ago` 파라미터 활용
- 호가(Orderbook) 조회

```bash
go run examples/quotation_kr.go
```

---

### `orders_kr.go` — 주문 생성 및 관리

지정가 매수 주문의 생성부터 조회, 취소까지의 전체 흐름을 다룹니다.

**포함 기능:**
- 주문 가능 정보 확인 (수수료율, 최소 주문금액)
- 호가 단위(tick size) 조회 및 주문 가격 계산
- 지정가 매수 주문 생성 → 상태 조회 → 취소
- 보유 자산 확인 및 시장가 매도 가능 여부 점검
- 완료 주문 목록 조회

```bash
# Dry run (기본) — 조회만, 주문 없음
UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> go run examples/orders_kr.go

# 실제 실행 — 지정가 매수 → 조회 → 취소
DRY_RUN=false UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> go run examples/orders_kr.go
```

---

### `orders_validate_kr.go` — 주문 유형 전체 검증 (KRW, TestNew)

주문 생성 테스트 API를 사용하여 KRW 마켓에서 모든 주문 유형을 검증합니다.
주문 생성 테스트 API는 실제 주문과 동일한 검증 과정을 거치지만, 주문이 생성되지 않아 수수료 없이 안전하게 테스트할 수 있습니다. 응답으로 반환된 UUID는 주문 조회/취소에 사용할 수 없습니다.

**포함 기능:**
- 주문 가능 정보 확인 (수수료율, 최소 주문금액, 지원 주문 유형)
- 가격을 호가 단위(tick size)로 내림한 뒤 지정가 매수/매도
- 시장가 매수(price) / 시장가 매도(volume)
- 최유리 매수/매도 (IOC)
- 에러 검증 (존재하지 않는 마켓, 가격 누락)

```bash
UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> go run examples/orders_validate_kr.go
```

---

### `deposits_kr.go` — 입금 주소와 입금 내역 관리

입금 주소 관리부터 내역 조회까지 입금 관련 기능 전반을 다룹니다.

**포함 기능:**
- 입금 주소 목록 조회
- 특정 코인 입금 주소 조회 및 신규 생성
- 입금 가능 정보 확인 (입금 한도, 최소 입금액)
- 입금 내역 조회 — 목록 / UUID, TxID 기반 단건 조회
- 트래블룰(Travel Rule) 지원 거래소(VASP) 목록 조회

```bash
# Dry run (기본) — 조회만
UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> go run examples/deposits_kr.go

# 실제 실행 — 입금 주소 생성 포함
DRY_RUN=false UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> go run examples/deposits_kr.go
```

---

### `withdrawals_kr.go` — 출금 관리

출금 허용 주소 확인부터 실제 출금까지 출금 관련 기능 전반을 다룹니다.
실행 전에는 반드시 주소, 네트워크, 금액을 다시 확인하세요.

> **주의:** 출금은 취소할 수 없습니다. 실제 출금 실행 전 주소와 금액을 반드시 확인하세요.

**포함 기능:**
- 출금 허용 주소 목록 조회
- 입출금 서비스 상태 확인 (BTC, ETH, USDT 기준)
- 출금 가능 정보 조회 (출금 수수료, 최소 금액, 일일 한도)
- 출금 내역 조회 — 목록 / 단건
- 디지털 자산 출금 및 원화(KRW) 출금

```bash
# Dry run (기본) — 조회만
UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> go run examples/withdrawals_kr.go

# 실제 실행 — 원화 출금 요청 포함 (기본 10,000원, kakao)
DRY_RUN=false UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> go run examples/withdrawals_kr.go
```

---

`DRY_RUN=false`이면 [`withdrawals_kr.go`](./withdrawals_kr.go)는 코드 수정 없이 기본 금액 10,000원·`kakao` 인증수단으로 원화 출금 요청을 보낼 수 있습니다. 디지털 자산 출금은 실제 출금 주소를 입력해야 합니다.

### `dca_kr.go` — DCA 정기 매수 자동화

일정 금액을 반복 시장가 매수하는 DCA(Dollar Cost Averaging) 전략을 구현합니다.
주기적 매수 자동화를 어떻게 구성하는지 흐름 중심으로 살펴볼 수 있습니다.

> 이 예제는 서버(또는 스케줄러) 환경이 갖추어져 있음을 전제로 작성되었습니다.
> 지속적인 자동 실행 환경은 직접 구성해야 합니다.

**포함 기능:**
- 현재가 조회 및 매수 설정 확인
- 지정 횟수만큼 시장가 매수 반복
- 주문 체결 대기 및 상태 폴링
- 매수 결과 요약 (총 매수금액, 총 매수수량, 평균 단가)

```bash
# Dry run (기본) — 현재가 조회만, 주문 없음
UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> go run examples/dca_kr.go

# 실제 실행
DRY_RUN=false UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> go run examples/dca_kr.go
```

---

### `tp_sl_kr.go` — TP/SL 익절/손절 자동 매도

현재가를 폴링하여 목표가(익절) 또는 손절가에 도달하면 시장가 매도를 실행합니다.

> 이 예제는 서버(또는 스케줄러) 환경이 갖추어져 있음을 전제로 작성되었습니다.

**포함 기능:**
- 현재가 조회 및 TP/SL 가격 계산
  - Dry run: 현재가 기준
  - 실제 실행: 보유 자산의 평균 매수가 기준
- REST 폴링 방식 가격 모니터링 (1초 간격)
- 목표가 도달 시 시장가 매도 실행
- 체결 대기 및 결과 확인

```bash
# Dry run (기본) — 모니터링만, 주문 없음
UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> go run examples/tp_sl_kr.go

# 실제 실행 — 평균 매수가 기준 자동 매도
DRY_RUN=false UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> go run examples/tp_sl_kr.go
```

---

### `indicators_kr.go` — 투자 지표 산출

Quotation API 데이터를 활용하여 기술적 투자 지표를 계산합니다.
외부 라이브러리 없이 순수 Go로 구현되어 있습니다.

**포함 기능:**
- 24시간 누적 거래대금 상위 5개 페어 조회
- RSI(상대강도지수) 계산
  - 일봉 데이터 200개 기반
  - Wilder 평활법
  - RSI 값 해석 (과매수 / 과매도 / 중립)

```bash
# 인증 불필요
go run examples/indicators_kr.go
```

필요하다면 이 README에서 흐름을 먼저 파악한 뒤, 더 세분화된 엔드포인트 예제는 기능별 하위 디렉토리에서 확인하세요.
