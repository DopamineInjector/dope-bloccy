# dope-bloccy
A http rest server used for interfacing between main application backend and dopecoin blockchain ecosystem

## Endpoint docs
### User management
#### GET /api/wallet/<user_id>
Fetches info about user's wallet, i.e. public key and creation date.  
The route param `<user_id` is the id of the user in our backend system
Response (200):
```ts
class WalletResponseDTO {
  id: string; // User's id as in our casino system
  publicKey: string; // User's blockchain id
  created: string; // Wallet's creation date
}
```
Throws 404 if user's wallet does not exist in the system
#### POST /api/wallet/<user_id> (__AUTH__)
Creates a new blockchain wallet for a user with `<user_id>` in our backend system.
Response (201):
empty
Throws 409 if user already has a wallet in our system

### NFT
#### GET /api/wallet/<user_id>/nfts
Fetches info about all the nfts that a user with `<user_id>` owns.  
Response (200):
```ts
class NftResponseDTO {
  nfts: {
    id: string // Nft metadata id
    description: string // NFT description string
    image: string | null // A base64 encoded byte array representing nft image
  }[]
}
```
Throws 404 if user's wallet does not exist in the system

#### POST /api/nft/mint (__AUTH__)
Mints a new nft and transfers it to a user's wallet
Request:
```ts
class MintRequestDTO {
  user: string // User's id in our backend system
  image: string // A base64 encoded byte array representing nft image
  description: // NFT description
}
```
Response (201):
empty

Throws many things, mostly 500 i think
## Auth
We do have an optional auth 'middleware' for use with the creation endpoints. It uses RSA signatures in order to establish the identity of the server issuing the request. 
To pass an auth check, the request body has to be signed with the admin's private key utilizing SHA256+RSA. Then, the signature has to be base64 encoded and included in the `x-auth-signature` header of the request.

