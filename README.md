# Baymax

## API setup

Go API is now self-contained under `api/`.

```sh
cp api/.env.example api/.env
make -C api run
```

The API listens on `http://localhost:8080` by default.

For hot reload:

```sh
make -C api air
```
