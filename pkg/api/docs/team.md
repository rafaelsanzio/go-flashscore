#### Creating a Team

```http
  POST /teams
```

| Header  | Type     | Description                |
| :------ | :------- | :------------------------- |
| `Token` | `Bearer` | **Required**. Your API key |

| Parameter    | Type     | Description                   |
| :----------- | :------- | :---------------------------- |
| `name`       | `string` | **Required**. Team name       |
| `short_code` | `string` | **Required**. Team short code |
| `country`    | `string` | **Required**. Team country    |
| `city`       | `string` | **Required**. Team city       |

#### Updating a Team

```http
  PUT /teams/{id}
```

| Header  | Type     | Description                |
| :------ | :------- | :------------------------- |
| `Token` | `Bearer` | **Required**. Your API key |

| Parameter    | Type     | Description                   |
| :----------- | :------- | :---------------------------- |
| `name`       | `string` | **Required**. Team name       |
| `short_code` | `string` | **Required**. Team short code |
| `country`    | `string` | **Required**. Team country    |
| `city`       | `string` | **Required**. Team city       |

#### Deleting a Team

```http
  DELETE /teams/{id}
```

| Header  | Type     | Description                |
| :------ | :------- | :------------------------- |
| `Token` | `Bearer` | **Required**. Your API key |

#### Getting a Team

```http
  GET /teams/{id}
```

| Header  | Type     | Description                |
| :------ | :------- | :------------------------- |
| `Token` | `Bearer` | **Required**. Your API key |

#### Listing all Teams

```http
  GET /teams
```

| Header  | Type     | Description                |
| :------ | :------- | :------------------------- |
| `Token` | `Bearer` | **Required**. Your API key |
