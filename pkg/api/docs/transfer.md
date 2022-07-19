#### Creating a Transfer

```http
  POST /transfers
```

| Header  | Type     | Description                |
| :------ | :------- | :------------------------- |
| `Token` | `Bearer` | **Required**. Your API key |

| Parameter          | Type     | Description                      |
| :----------------- | :------- | :------------------------------- |
| `player`           | `string` | **Required**. Player id          |
| `team_destiny`     | `string` | **Required**. Team id            |
| `amount`           | `money`  | **Required**. Amount of transfer |
| `date_of_transfer` | `date`   | **Required**. Date of Transfer   |

#### Getting a Transfer

```http
  GET /transfers/{id}
```

| Header  | Type     | Description                |
| :------ | :------- | :------------------------- |
| `Token` | `Bearer` | **Required**. Your API key |

#### Listing all transfers

```http
  GET /transfers
```

| Header  | Type     | Description                |
| :------ | :------- | :------------------------- |
| `Token` | `Bearer` | **Required**. Your API key |
