# Secure Payment Service Challenge
Para levantar la aplicacion se puede usar Docker Desktop y ejecutar `docker compose up --build`

Esto disponibiliza algunos endpoints:
- `localhost:8080`: se pueden acceder a la base de datos. Hay que cambiar `Motor de base de datos` a `PostgreSQL`. El usuario es `postgres` y la contraseña `go`

- `localhost:8000`: una API que tiene los siguiente endpoints
    - `/ping`: esto se puede usar para chequear que la API haya levantado
    ```
    curl http://localhost:8000/ping
    ```
    - `/metrics`: devuelve algunas metricas generales como cantidad de usuarios registrados y cantidad de transferencias realizadas.
    ```
    curl http://localhost:8000/metrics
    ```
    - `/users`: crea un user. Para esto recibe una estructura como la siguiente:
        ```
        {
        "Name": "Rocio",
        "DNI": "14548789",
        "Email": "rocio@gmail.com"
        }
        ```
        Es un endpoint que necesita autenticación
        ```
        curl 'http://localhost:8000/users' \
            --header 'Authorization: 7rXBD1gtRblZtxyFhtEVqb6B1S1SG2dyCcbrKcOyrkSjTwIMgXoKGk0Ci7ek7FDy' \
            --header 'Role: admin' \
            --header 'Content-Type: application/json' \
            --data-raw '{
                "Name": "Rocio",
                "DNI": "14548789",
                "Email": "rocio@gmail.com"
        }'
        ```
    - `/users/{id}/balance`: devuelve el balance del usuario con el id dado. Requiere autenticación.
        ```
        curl --location 'localhost:8000/users/1/balance' \
        --header 'Authorization: 7rXBD1gtRblZtxyFhtEVqb6B1S1SG2dyCcbrKcOyrkSjTwIMgXoKGk0Ci7ek7FDy' \
        --header 'Role: admin'
        ```
    - `/transfers`: crea una transferencia. Para esto recibe una estructura como la siguiente: 
        ```
        {
            "from_user_id": 1,
            "to_user_id": 2,
            "amount": 100,
            "description": "alguna description"
        }
        ```
        Requiere autenticación.
        ```
        curl --location 'localhost:8000/transfers' \
            --header 'Authorization: 7rXBD1gtRblZtxyFhtEVqb6B1S1SG2dyCcbrKcOyrkSjTwIMgXoKGk0Ci7ek7FDy' \
            --header 'Role: admin' \
            --header 'Content-Type: application/json' \
            --data '{
                "from_user_id": 1,
                "to_user_id": 2,
                "amount": 100,
                "description": "alguna description"
            }'
        ```
    - `/transfers/finish`: Este endpoint es para notificar el estado de una transferencia. 
        ```
        {
            "id": 36,
            "status": "COMPLETE"
        }
        ```
        Requiere autenteticación como webhook
        ```
        curl --location 'localhost:8000/transfers/finish' \
            --header 'Authorization: buQs7IZKA5j05WuC21cgSt9DGQlqCqgrbO6EXfgWa48Lw9qkbFLRKWXKQMQMPLQb' \
            --header 'Role: webhook' \
            --header 'Content-Type: application/json' \
            --data '{
                "id": 36,
                "status": "COMPLETE"
            }'
        ```
    - `/transfers/{id}`: devuelve toda la información de una transferencia. Requiere autenticación como admin.
        ```
        curl --location 'localhost:8000/transfers/1' \
            --header 'Authorization: 7rXBD1gtRblZtxyFhtEVqb6B1S1SG2dyCcbrKcOyrkSjTwIMgXoKGk0Ci7ek7FDy' \
            --header 'Role: admin'
        ```

# Observaciones
Algunas decisiones que tome para simplicar algunas cosas:
- Se pueden crear transferencias para usuarios que no esten creados en el sistema pero si el usuario no esta en el sistema no se va a tener su balance. Me parecio que era lo mejor por un tema de ingresar dinero. Si verificaba que ambos esten, nunca se iba a poder generar dinero. Verificar que alguno este registrado me parecio que no tenia sentido. 
- Con respecto a autenticación decidí separarlo en admin y webhook. Por ahí lo correcto era generarle un token a los usuarios cuando se crean la cuenta pero esto no permitiria lo dicho en el punto anterior. 
