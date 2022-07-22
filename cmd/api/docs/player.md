#### Creating a Player

```http
  POST /players
```

| Header  | Type     | Description                |
| :------ | :------- | :------------------------- |
| `Token` | `Bearer` | **Required**. Your API key |

| Parameter       | Type     | Description                   |
| :-------------- | :------- | :---------------------------- |
| `name`          | `string` | **Required**. Player name     |
| `team`          | `string` | **Required**. Team player     |
| `country`       | `string` | **Required**. Player country  |
| `birthday_date` | `date`   | **Required**. Player birthday |

#### Updating a player

```http
  PUT /players/{id}
```

| Header  | Type     | Description                |
| :------ | :------- | :------------------------- |
| `Token` | `Bearer` | **Required**. Your API key |

| Parameter       | Type     | Description                   |
| :-------------- | :------- | :---------------------------- |
| `name`          | `string` | **Required**. Player name     |
| `short_code`    | `string` | **Required**. Team player     |
| `country`       | `string` | **Required**. Player country  |
| `birthday_date` | `date`   | **Required**. Player birthday |

#### Deleting a player

```http
  DELETE /players/{id}
```

| Header  | Type     | Description                |
| :------ | :------- | :------------------------- |
| `Token` | `Bearer` | **Required**. Your API key |

#### Getting a Player

```http
  GET /players/{id}
```

| Header  | Type     | Description                |
| :------ | :------- | :------------------------- |
| `Token` | `Bearer` | **Required**. Your API key |

#### Listing all Players

```http
  GET /players
```

| Header  | Type     | Description                |
| :------ | :------- | :------------------------- |
| `Token` | `Bearer` | **Required**. Your API key |
