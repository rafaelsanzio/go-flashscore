#### Creating a Tournament

```http
  POST /tournaments
```

| Header  | Type     | Description                |
| :------ | :------- | :------------------------- |
| `Token` | `Bearer` | **Required**. Your API key |

| Parameter | Type       | Description                   |
| :-------- | :--------- | :---------------------------- |
| `name`    | `string`   | **Required**. Tournament name |
| `teams`   | `[]string` | **Required**. Teams id        |

#### Updating a Tournament

```http
  PUT /tournaments/{id}
```

| Header  | Type     | Description                |
| :------ | :------- | :------------------------- |
| `Token` | `Bearer` | **Required**. Your API key |

| Parameter | Type     | Description                   |
| :-------- | :------- | :---------------------------- |
| `name`    | `string` | **Required**. Tournament name |

#### Deleting a Tournament

```http
  DELETE /tournaments/{id}
```

| Header  | Type     | Description                |
| :------ | :------- | :------------------------- |
| `Token` | `Bearer` | **Required**. Your API key |

#### Getting a Tournament

```http
  GET /tournaments/{id}
```

| Header  | Type     | Description                |
| :------ | :------- | :------------------------- |
| `Token` | `Bearer` | **Required**. Your API key |

#### Listing all Tournaments

```http
  GET /tournaments
```

| Header  | Type     | Description                |
| :------ | :------- | :------------------------- |
| `Token` | `Bearer` | **Required**. Your API key |

#### Addind teams to a Tournament

```http
  POST /tournaments/{id}/add-teams
```

| Header  | Type     | Description                |
| :------ | :------- | :------------------------- |
| `Token` | `Bearer` | **Required**. Your API key |

| Parameter | Type       | Description            |
| :-------- | :--------- | :--------------------- |
| `teams`   | `[]string` | **Required**. Teams id |
