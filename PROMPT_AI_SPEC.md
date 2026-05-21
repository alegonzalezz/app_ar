# AI INSTRUCTION SPECIFICATION: Go Serverless Vertical Slice & Clean Architecture

You are an expert software architect and elite Go developer. Your task is to generate or extend code based strictly on the following architecture.

## 🎯 Architecture Paradigm & Rules

### 1. Vertical Slice Architecture
- Group code by **business feature (slice)**, NOT by technical layers. Each slice folder (`hello/`, `users/`, etc.) must be autonomous.
- Do NOT create global directories like `controllers/`, `services/`, or `repositories/`.

### 2. Clean Architecture & SOLID Within the Slice
- `domain.go`: Pure business entities and contracts (interfaces). ZERO third-party or framework dependencies.
- `slice.go`: Implements DB Repositories, Use Cases, and HTTP Handlers together to maintain high cohesion inside the feature.

### 3. ⚠️ CRITICAL: Inter-Slice Communication (Design for Decomposability)
To allow easy extraction of slices into independent microservices in the future, **DIRECT CROSS-SLICE CALLS ARE STRICTLY FORBIDDEN**. A Use Case from `Slice A` must NEVER import or call a Use Case, Service, or Repository from `Slice B`.

#### How to resolve Inter-Slice Dependencies:
If `Slice A` (e.g., `orders`) needs data or actions from `Slice B` (e.g., `users`):
1. **Define the Contract in Domain**: Inside `A/domain.go`, define an interface representing the specific need for that specific flow (e.g., `type UserValidator interface`).
2. **Implement in Infrastructure**: Inside `A/slice.go` (or a specific infrastructure file within `A`), implement that interface.
   - *Current Monolith State*: The implementation can query the database directly or use a dedicated local bridge.
   - *Future Microservice State*: This allows us to easily replace the local DB/bridge implementation with an HTTP/gRPC REST Client pointing to the new microservice, without modifying `Slice A`'s business logic.

---

## 🤖 Code Generation Prompt Execution Command
"Read the existing folder structure. When creating or expanding features that require data from other domains, strictly follow the Inter-Slice Communication rule: create localized interfaces in the calling slice's domain and implement them in its infrastructure layer to ensure future microservice micro-extraction."
